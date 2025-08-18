package services

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"stock-api/internal/clients"
	"stock-api/internal/models"
	"stock-api/internal/repository"
)

type StockService struct {
	repo           *repository.StockRepository
	processRepo    *repository.ProcessControlRepository
	karenAIClient  *clients.KarenAIClient
	recommendation *RecommendationEngine
	recScoreRepo   *repository.RecommendationScoreRepository
}

func NewStockService(db *sql.DB, apiKey string) *StockService {
	repo := repository.NewStockRepository(db)
	processRepo := repository.NewProcessControlRepository(db)
	karenAIClient := clients.NewKarenAIClient(apiKey)
	recScoreRepo := repository.NewRecommendationScoreRepository(db)

	return &StockService{
		repo:           repo,
		processRepo:    processRepo,
		karenAIClient:  karenAIClient,
		recommendation: NewRecommendationEngine(),
		recScoreRepo:   recScoreRepo,
	}
}

func (s *StockService) CanStartStockSync() (bool, error) {
	return s.processRepo.CanStartStockSync()
}


func (s *StockService) GetStocksWithMetricsPaginated(page, pageSize int, filters models.StockFilterParams) (*models.PaginatedResponse[models.StockWithAnalysis], error) {
	return s.repo.GetStocksWithAnalysisPaginated(page, pageSize, filters)
}

func (s *StockService) GetStockWithMetrics(symbol string) (*models.StockWithAnalysis, error) {
	stock, err := s.repo.GetStockBySymbol(symbol)
	if err != nil || stock == nil {
		return nil, err
	}

	stockWithAnalysis := models.StockWithAnalysis{Stock: *stock}

	analyses, err := s.repo.GetLatestAnalysisForStock(stock.ID, 5)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stockWithAnalysis.LatestAnalysis = analyses

	return &stockWithAnalysis, nil
}

func (s *StockService) SearchAndAddStock(symbol string) (*models.StockWithAnalysis, error) {
	existingStock, err := s.repo.GetStockBySymbol(symbol)
	if err != nil {
		return nil, err
	}

	if existingStock != nil {
		return s.GetStockWithMetrics(symbol)
	}

	return nil, fmt.Errorf("stock not found and cannot search individual stocks in KarenAI API")
}

func (s *StockService) RefreshStockData(symbol string) error {
	stock, err := s.repo.GetStockBySymbol(symbol)
	if err != nil {
		return err
	}

	if stock == nil {
		return fmt.Errorf("stock not found: %s", symbol)
	}

	return s.SyncAllStocks()
}

func (s *StockService) SyncAllStocks() error {
	// Start the process
	if err := s.processRepo.StartStockSync(); err != nil {
		return fmt.Errorf("failed to start stock sync process: %w", err)
	}

	// Ensure process is marked as finished even if there's an error
	defer func() {
		if finishErr := s.processRepo.FinishStockSync(); finishErr != nil {
			fmt.Printf("Warning: failed to finish stock sync process: %v\n", finishErr)
		}
	}()

	nextPage := ""
	totalProcessed := 0

	for {
		response, err := s.karenAIClient.GetStocksList(nextPage)
		if err != nil {
			fmt.Println("🔴🔴 ~ Error fetching stocks from KarenAI API:", err)
			return fmt.Errorf("failed to fetch stocks from API: %w", err)
		}

		for _, apiAnalysis := range response.Items {
			stock := &models.Stock{
				Symbol: apiAnalysis.Ticker,
				Name:   apiAnalysis.Company,
			}

			if err := s.repo.CreateStock(stock); err != nil {
				fmt.Printf("Error: failed to create/update stock %s: %v\n", stock.Symbol, err)
				continue
			}

			fmt.Printf("Successfully created/updated stock %s (ID: %d)\n", stock.Symbol, stock.ID)

			analysis := &models.StockAnalysis{
				StockID:      stock.ID,
				TargetFrom:   apiAnalysis.TargetFrom,
				TargetTo:     apiAnalysis.TargetTo,
				Action:       apiAnalysis.Action,
				Brokerage:    apiAnalysis.Brokerage,
				RatingFrom:   apiAnalysis.RatingFrom,
				RatingTo:     apiAnalysis.RatingTo,
				AnalysisDate: apiAnalysis.Time,
			}

			if err := s.repo.CreateStockAnalysis(analysis); err != nil {
				fmt.Printf("Error: failed to create analysis for stock %s: %v\n", stock.Symbol, err)
				continue
			}

			fmt.Printf("Successfully created/updated analysis for stock %s (Analysis ID: %d)\n", stock.Symbol, analysis.ID)

			if err := s.repo.DeleteOldAnalysis(stock.ID, 10); err != nil {
				fmt.Printf("Warning: failed to cleanup old analysis for stock %s: %v\n", stock.Symbol, err)
			}

			// Calculate and store recommendation score for this stock
			if err := s.calculateAndStoreRecommendationScore(stock.ID); err != nil {
				fmt.Printf("Warning: failed to calculate recommendation score for stock %s: %v\n", stock.Symbol, err)
			}

			totalProcessed++
		}

		fmt.Printf("Processed %d stocks in page %s (total: %d)\n", len(response.Items), nextPage, totalProcessed)

		if response.NextPage == "" {
			break
		}

		nextPage = response.NextPage
	}

	fmt.Printf("Successfully processed %d stocks total\n", totalProcessed)
	return nil
}


func (s *StockService) GetRecommendationsPaginated(page, pageSize int) (*models.PaginatedResponse[models.StockRecommendation], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Get paginated recommendations from pre-calculated scores
	recommendations, err := s.recScoreRepo.GetTopRecommendationsPaginated(page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert to legacy format for backward compatibility
	var result []models.StockRecommendation
	for _, rec := range recommendations.Data {
		result = append(result, models.StockRecommendation{
			Stock:      rec.Stock,
			Score:      rec.TotalScore,
			Reason:     rec.Reason,
			Confidence: rec.Confidence,
		})
	}

	return &models.PaginatedResponse[models.StockRecommendation]{
		Data: result,
		Meta: recommendations.Meta,
	}, nil
}

type RecommendationEngine struct{}

func NewRecommendationEngine() *RecommendationEngine {
	return &RecommendationEngine{}
}

func (r *RecommendationEngine) AnalyzeStocks(stocks []models.StockWithAnalysis) []models.StockRecommendation {
	var recommendations []models.StockRecommendation

	for _, stock := range stocks {
		if len(stock.LatestAnalysis) == 0 {
			continue
		}

		score := r.calculateScore(stock)
		confidence := r.getConfidence(score)
		reason := r.generateReason(stock, score)

		recommendation := models.StockRecommendation{
			Stock:      stock,
			Score:      score,
			Reason:     reason,
			Confidence: confidence,
		}

		recommendations = append(recommendations, recommendation)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	return recommendations
}

func (r *RecommendationEngine) calculateScore(stock models.StockWithAnalysis) float64 {
	baseScore := 50.0
	ratingScore, ratingChangeScore := r.calculateRatingScores(stock)
	targetChangeScore := r.calculateTargetChangeScore(stock)
	actionScore := r.calculateActionScore(stock)
	coverageScore := r.calculateCoverageScore(stock)

	return baseScore + ratingScore + ratingChangeScore + targetChangeScore + actionScore + coverageScore
}

func (r *RecommendationEngine) calculateRatingScores(stock models.StockWithAnalysis) (float64, float64) {
	if len(stock.LatestAnalysis) == 0 {
		return 0, 0
	}

	latestAnalysis := stock.LatestAnalysis[0]
	ratingScore := 0.0
	ratingChangeScore := 0.0

	// Analyze rating changes
	if latestAnalysis.RatingTo != "" {
		ratingScore = r.getRatingScore(latestAnalysis.RatingTo) - 50 // Subtract base to get delta
	}

	if latestAnalysis.RatingTo != "" && latestAnalysis.RatingFrom != "" {
		toScore := r.getRatingScore(latestAnalysis.RatingTo)
		fromScore := r.getRatingScore(latestAnalysis.RatingFrom)

		// Bonus for rating upgrades
		if toScore > fromScore {
			ratingChangeScore = 15
		} else if toScore < fromScore {
			ratingChangeScore = -10
		}
	}

	return ratingScore, ratingChangeScore
}

func (r *RecommendationEngine) calculateTargetChangeScore(stock models.StockWithAnalysis) float64 {
	if len(stock.LatestAnalysis) == 0 {
		return 0
	}

	latestAnalysis := stock.LatestAnalysis[0]
	targetFromVal := r.extractPrice(latestAnalysis.TargetFrom)
	targetToVal := r.extractPrice(latestAnalysis.TargetTo)

	if targetFromVal > 0 && targetToVal > 0 {
		targetChange := (targetToVal - targetFromVal) / targetFromVal

		if targetChange > 0.1 { // Target raised by >10%
			return 20
		} else if targetChange > 0.05 { // Target raised by >5%
			return 10
		} else if targetChange < -0.1 { // Target lowered by >10%
			return -15
		} else if targetChange < -0.05 { // Target lowered by >5%
			return -8
		}
	}

	return 0
}

func (r *RecommendationEngine) calculateActionScore(stock models.StockWithAnalysis) float64 {
	if len(stock.LatestAnalysis) == 0 {
		return 0
	}

	latestAnalysis := stock.LatestAnalysis[0]
	action := strings.ToLower(latestAnalysis.Action)
	
	if strings.Contains(action, "initiated") {
		return 10
	} else if strings.Contains(action, "raised") {
		return 12
	} else if strings.Contains(action, "lowered") {
		return -8
	} else if strings.Contains(action, "maintained") {
		return 5
	}

	return 0
}

func (r *RecommendationEngine) calculateCoverageScore(stock models.StockWithAnalysis) float64 {
	score := 0.0

	// Multiple recent analyses bonus
	if len(stock.LatestAnalysis) >= 3 {
		score += 5
	}

	// Check for consistent positive sentiment
	positiveCount := 0
	for _, analysis := range stock.LatestAnalysis {
		if r.getRatingScore(analysis.RatingTo) > 60 {
			positiveCount++
		}
	}

	if positiveCount >= 2 {
		score += 8
	}

	return score
}

func (r *RecommendationEngine) getRatingScore(rating string) float64 {
	rating = strings.ToLower(rating)

	switch {
	case strings.Contains(rating, "strong buy") || strings.Contains(rating, "buy"):
		return 80
	case strings.Contains(rating, "outperform") || strings.Contains(rating, "overweight"):
		return 70
	case strings.Contains(rating, "hold") || strings.Contains(rating, "neutral"):
		return 50
	case strings.Contains(rating, "underperform") || strings.Contains(rating, "underweight"):
		return 30
	case strings.Contains(rating, "sell") || strings.Contains(rating, "strong sell"):
		return 10
	default:
		return 50
	}
}

func (r *RecommendationEngine) extractPrice(priceStr string) float64 {
	if priceStr == "" {
		return 0
	}

	// Remove $ and any other non-numeric characters except decimal point
	cleaned := ""
	for _, r := range priceStr {
		if r >= '0' && r <= '9' || r == '.' {
			cleaned += string(r)
		}
	}

	price, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0
	}

	return price
}

func (r *RecommendationEngine) getConfidence(score float64) string {
	if score >= 75 {
		return "High"
	} else if score >= 60 {
		return "Medium"
	} else {
		return "Low"
	}
}

func (r *RecommendationEngine) generateReason(stock models.StockWithAnalysis, score float64) string {
	if len(stock.LatestAnalysis) == 0 {
		return "No recent analyst coverage"
	}

	latestAnalysis := stock.LatestAnalysis[0]
	reasons := []string{}

	// Check rating
	if latestAnalysis.RatingTo != "" {
		rating := strings.ToLower(latestAnalysis.RatingTo)
		if strings.Contains(rating, "buy") {
			reasons = append(reasons, "Buy rating from "+latestAnalysis.Brokerage)
		} else if strings.Contains(rating, "outperform") {
			reasons = append(reasons, "Outperform rating from "+latestAnalysis.Brokerage)
		}
	}

	// Check price target
	targetFromVal := r.extractPrice(latestAnalysis.TargetFrom)
	targetToVal := r.extractPrice(latestAnalysis.TargetTo)

	if targetFromVal > 0 && targetToVal > 0 {
		targetChange := (targetToVal - targetFromVal) / targetFromVal
		if targetChange > 0.1 {
			reasons = append(reasons, fmt.Sprintf("Price target raised by %.1f%%", targetChange*100))
		}
	}

	// Check action
	if strings.Contains(strings.ToLower(latestAnalysis.Action), "initiated") {
		reasons = append(reasons, "New analyst coverage")
	}

	// Check multiple analyses
	if len(stock.LatestAnalysis) >= 3 {
		reasons = append(reasons, "Multiple recent analyst updates")
	}

	if len(reasons) == 0 {
		return "Analyst coverage available from " + latestAnalysis.Brokerage
	}

	return strings.Join(reasons, ", ")
}

func (s *StockService) GetFilterOptions() (*models.FilterOptions, error) {
	return s.repo.GetFilterOptions()
}

func (s *StockService) GetMarketIntelligenceOverview() (*models.MarketIntelligenceOverview, error) {
	// Get basic analytics from repository
	overview, err := s.repo.GetMarketIntelligenceOverview()
	if err != nil {
		return nil, err
	}

	// Get recommendation statistics from the scores table
	stats, err := s.recScoreRepo.GetRecommendationStats()
	if err != nil {
		// If recommendations fail, still return basic overview
		fmt.Printf("Warning: failed to get recommendation stats for analytics: %v\n", err)
		overview.TotalRecommendations = 0
		overview.HighConfidenceRecs = 0
		overview.SelectionRate = 0
		overview.AverageRecommendationScore = 0
		return overview, nil
	}

	// Extract recommendation metrics from stats
	overview.TotalRecommendations = stats["total_recommendations"].(int)
	overview.HighConfidenceRecs = stats["high_confidence"].(int)
	overview.AverageRecommendationScore = stats["average_score"].(float64)

	// Calculate selection rate (recommendations per stock)
	if overview.TotalStocks > 0 {
		overview.SelectionRate = float64(overview.TotalRecommendations) / float64(overview.TotalStocks) * 100
	}

	return overview, nil
}

func (s *StockService) calculateAndStoreRecommendationScore(stockID int) error {
	// Get stock with analysis data
	stockWithAnalysis, err := s.GetStockWithMetrics(s.getStockSymbolByID(stockID))
	if err != nil {
		return fmt.Errorf("failed to get stock with analysis: %w", err)
	}

	if stockWithAnalysis == nil {
		return fmt.Errorf("stock not found")
	}

	// Calculate individual scores
	ratingScore, ratingChangeScore := s.recommendation.calculateRatingScores(*stockWithAnalysis)
	targetChangeScore := s.recommendation.calculateTargetChangeScore(*stockWithAnalysis)
	actionScore := s.recommendation.calculateActionScore(*stockWithAnalysis)
	coverageScore := s.recommendation.calculateCoverageScore(*stockWithAnalysis)

	// Calculate total score
	baseScore := 50.0
	totalScore := baseScore + ratingScore + ratingChangeScore + targetChangeScore + actionScore + coverageScore

	// Get confidence and reason
	confidence := s.recommendation.getConfidence(totalScore)
	reason := s.recommendation.generateReason(*stockWithAnalysis, totalScore)

	// Get latest analysis ID if available
	var latestAnalysisID *int
	if len(stockWithAnalysis.LatestAnalysis) > 0 {
		latestAnalysisID = &stockWithAnalysis.LatestAnalysis[0].ID
	}

	// Create recommendation score record
	score := &models.RecommendationScore{
		StockID:           stockID,
		TotalScore:        totalScore,
		RatingScore:       ratingScore,
		RatingChangeScore: ratingChangeScore,
		TargetChangeScore: targetChangeScore,
		ActionScore:       actionScore,
		CoverageScore:     coverageScore,
		Confidence:        confidence,
		Reason:            reason,
		LatestAnalysisID:  latestAnalysisID,
	}

	// Store in database
	return s.recScoreRepo.UpsertRecommendationScore(score)
}

func (s *StockService) getStockSymbolByID(stockID int) string {
	// Helper method to get stock symbol by ID
	query := `SELECT symbol FROM stocks WHERE id = $1`
	var symbol string
	err := s.repo.DB().QueryRow(query, stockID).Scan(&symbol)
	if err != nil {
		return ""
	}
	return symbol
}
