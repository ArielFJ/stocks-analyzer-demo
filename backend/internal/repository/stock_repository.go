package repository

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"stock-api/internal/models"
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) DB() *sql.DB {
	return r.db
}

func (r *StockRepository) CreateStock(stock *models.Stock) error {
	// First try to get existing stock
	existing, err := r.GetStockBySymbol(stock.Symbol)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing stock %s: %w", stock.Symbol, err)
	}

	if existing != nil {
		// Stock exists, update it and return the existing data
		updateQuery := `
			UPDATE stocks 
			SET name = $1, updated_at = NOW() 
			WHERE symbol = $2
			RETURNING id, created_at, updated_at`

		err := r.db.QueryRow(updateQuery, stock.Name, stock.Symbol).
			Scan(&stock.ID, &stock.CreatedAt, &stock.UpdatedAt)
		if err != nil {
			return fmt.Errorf("error updating stock %s: %w", stock.Symbol, err)
		}
		return nil
	}

	// Stock doesn't exist, create new one
	insertQuery := `
		INSERT INTO stocks (symbol, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(insertQuery, stock.Symbol, stock.Name).
		Scan(&stock.ID, &stock.CreatedAt, &stock.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error inserting stock %s: %w", stock.Symbol, err)
	}

	return nil
}

func (r *StockRepository) GetStockBySymbol(symbol string) (*models.Stock, error) {
	query := `
		SELECT id, symbol, name, created_at, updated_at
		FROM stocks WHERE symbol = $1`

	stock := &models.Stock{}
	err := r.db.QueryRow(query, symbol).Scan(
		&stock.ID, &stock.Symbol, &stock.Name, &stock.CreatedAt, &stock.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return stock, err
}


func (r *StockRepository) CreateStockAnalysis(analysis *models.StockAnalysis) error {
	// First check if analysis already exists
	checkQuery := `
		SELECT id, created_at FROM stock_analysis 
		WHERE stock_id = $1 AND analysis_date = $2 AND brokerage = $3`

	var existingID int
	var existingCreatedAt time.Time
	err := r.db.QueryRow(checkQuery, analysis.StockID, analysis.AnalysisDate, analysis.Brokerage).
		Scan(&existingID, &existingCreatedAt)

	if err == nil {
		// Analysis exists, check if we need to update
		updateQuery := `
			UPDATE stock_analysis 
			SET target_from = $1, target_to = $2, action = $3, rating_from = $4, rating_to = $5
			WHERE stock_id = $6 AND analysis_date = $7 AND brokerage = $8
			AND (target_from IS DISTINCT FROM $1 OR target_to IS DISTINCT FROM $2 
				 OR action IS DISTINCT FROM $3 OR rating_from IS DISTINCT FROM $4 
				 OR rating_to IS DISTINCT FROM $5)
			RETURNING id, created_at`

		err = r.db.QueryRow(updateQuery, analysis.TargetFrom, analysis.TargetTo, analysis.Action,
			analysis.RatingFrom, analysis.RatingTo, analysis.StockID, analysis.AnalysisDate, analysis.Brokerage).
			Scan(&analysis.ID, &analysis.CreatedAt)

		if err == sql.ErrNoRows {
			// No update needed, use existing values
			analysis.ID = existingID
			analysis.CreatedAt = existingCreatedAt
			return nil
		}
		return err
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing analysis: %w", err)
	}

	// Analysis doesn't exist, create new one
	insertQuery := `
		INSERT INTO stock_analysis (stock_id, target_from, target_to, action, brokerage, rating_from, rating_to, analysis_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	err = r.db.QueryRow(insertQuery, analysis.StockID, analysis.TargetFrom, analysis.TargetTo,
		analysis.Action, analysis.Brokerage, analysis.RatingFrom, analysis.RatingTo, analysis.AnalysisDate).
		Scan(&analysis.ID, &analysis.CreatedAt)

	if err != nil {
		return fmt.Errorf("error inserting analysis for stock_id %d: %w", analysis.StockID, err)
	}

	return nil
}

func (r *StockRepository) GetLatestAnalysisForStock(stockID int, limit int) ([]models.StockAnalysis, error) {
	query := `
		SELECT id, stock_id, target_from, target_to, action, brokerage, rating_from, rating_to, analysis_date, created_at
		FROM stock_analysis
		WHERE stock_id = $1
		ORDER BY analysis_date DESC
		LIMIT $2`

	rows, err := r.db.Query(query, stockID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []models.StockAnalysis
	for rows.Next() {
		var analysis models.StockAnalysis
		err := rows.Scan(
			&analysis.ID, &analysis.StockID, &analysis.TargetFrom, &analysis.TargetTo,
			&analysis.Action, &analysis.Brokerage, &analysis.RatingFrom, &analysis.RatingTo,
			&analysis.AnalysisDate, &analysis.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		analyses = append(analyses, analysis)
	}

	return analyses, nil
}


func (r *StockRepository) GetStocksWithAnalysisPaginated(page, pageSize int, filters models.StockFilterParams) (*models.PaginatedResponse[models.StockWithAnalysis], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Build where conditions and parameters for filtering
	whereConditions := []string{}
	queryArgs := []any{}
	argIndex := 1

	// Filter by brokerage
	if filters.Brokerage != "" && filters.Brokerage != "all" {
		whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.brokerage) LIKE LOWER($%d))", argIndex))
		queryArgs = append(queryArgs, "%"+filters.Brokerage+"%")
		argIndex++
	}

	// Filter by action type
	if filters.ActionType != "" && filters.ActionType != "all" {
		switch filters.ActionType {
		case "initiated":
			whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.action) LIKE '%%initiated%%')"))
		case "raised":
			whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.action) LIKE '%%raised%%')"))
		case "lowered":
			whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.action) LIKE '%%lowered%%')"))
		case "upgraded":
			whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.action) LIKE '%%upgraded%%')"))
		case "downgraded":
			whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.action) LIKE '%%downgraded%%')"))
		case "reiterated":
			whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.action) LIKE '%%reiterated%%')"))
		case "target-set":
			whereConditions = append(whereConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM stock_analysis sa WHERE sa.stock_id = s.id AND LOWER(sa.action) LIKE '%%target set%%')"))
		}
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + fmt.Sprintf("(%s)", whereConditions[0])
		for i := 1; i < len(whereConditions); i++ {
			whereClause += " AND " + fmt.Sprintf("(%s)", whereConditions[i])
		}
	}

	// Get total count with filters
	var totalCount int
	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT s.id) FROM stocks s %s", whereClause)
	err := r.db.QueryRow(countQuery, queryArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	offset := (page - 1) * pageSize
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Build main query with filters
	var queryTemplate string
	var finalOrderBy string

	if filters.SortBy == "analysis-newest" || filters.SortBy == "analysis-oldest" {
		// For analysis date sorting, use a subquery approach
		var sortDirection string

		if filters.SortBy == "analysis-newest" {
			sortDirection = "DESC"
		} else {
			sortDirection = "ASC"
		}

		queryTemplate = `
			WITH stocks_with_max_date AS (
				SELECT s.id, s.symbol, s.name, s.created_at, s.updated_at,
					   (SELECT MAX(analysis_date) FROM stock_analysis WHERE stock_id = s.id) as max_analysis_date
				FROM stocks s
				%s
			),
			paginated_stocks AS (
				SELECT id, symbol, name, created_at, updated_at, max_analysis_date
				FROM stocks_with_max_date
				ORDER BY max_analysis_date ` + sortDirection + ` NULLS LAST, symbol ASC
				LIMIT $%d OFFSET $%d
			)`

		// For analysis date sorting, preserve the order from the CTE
		finalOrderBy = "ORDER BY s.max_analysis_date " + sortDirection + " NULLS LAST, s.symbol ASC, sa.analysis_date DESC"
	} else {
		// Standard sorting approach
		orderBy := "s.symbol"
		switch filters.SortBy {
		case "newest":
			orderBy = "s.updated_at DESC"
		case "oldest":
			orderBy = "s.updated_at ASC"
		case "ticker-a-z":
			orderBy = "s.symbol ASC"
		case "company-a-z":
			orderBy = "s.name ASC"
		default:
			orderBy = "s.symbol ASC"
		}

		queryTemplate = `
			WITH paginated_stocks AS (
				SELECT s.id, s.symbol, s.name, s.created_at, s.updated_at, NULL as max_analysis_date
				FROM stocks s
				%s
				ORDER BY ` + orderBy + `
				LIMIT $%d OFFSET $%d
			)`

		finalOrderBy = "ORDER BY s.symbol ASC, sa.analysis_date DESC"
	}

	// Add the common SELECT part
	queryTemplate += `
		SELECT 
			s.id, s.symbol, s.name, s.created_at, s.updated_at, s.max_analysis_date,
			COALESCE(sa.id, 0) as analysis_id, COALESCE(sa.target_from, '') as target_from, 
			COALESCE(sa.target_to, '') as target_to, COALESCE(sa.action, '') as action,
			COALESCE(sa.brokerage, '') as brokerage, COALESCE(sa.rating_from, '') as rating_from,
			COALESCE(sa.rating_to, '') as rating_to, sa.analysis_date, sa.created_at as analysis_created_at
		FROM paginated_stocks s
		LEFT JOIN LATERAL (
			SELECT id, target_from, target_to, action, brokerage, rating_from, rating_to, analysis_date, created_at
			FROM stock_analysis
			WHERE stock_id = s.id
			ORDER BY analysis_date DESC
			LIMIT 3
		) sa ON true
		` + finalOrderBy

	query := fmt.Sprintf(queryTemplate, whereClause, argIndex, argIndex+1)
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		fmt.Printf("Query error for sort_by=%s: %v\n", filters.SortBy, err)
		fmt.Printf("Query: %s\n", query)
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	// Use a slice to preserve order instead of map
	var result []models.StockWithAnalysis
	stocksMap := make(map[int]int) // stock_id -> index in result slice

	for rows.Next() {
		var stock models.Stock
		var analysis models.StockAnalysis
		var analysisID sql.NullInt64
		var analysisDate sql.NullTime
		var analysisCreatedAt sql.NullTime
		var maxAnalysisDate sql.NullTime

		err := rows.Scan(
			&stock.ID, &stock.Symbol, &stock.Name, &stock.CreatedAt, &stock.UpdatedAt, &maxAnalysisDate,
			&analysisID, &analysis.TargetFrom, &analysis.TargetTo, &analysis.Action,
			&analysis.Brokerage, &analysis.RatingFrom, &analysis.RatingTo, &analysisDate, &analysisCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Check if we've seen this stock before
		stockIndex, exists := stocksMap[stock.ID]
		if !exists {
			// New stock - add to result slice
			stockWithAnalysis := models.StockWithAnalysis{
				Stock: stock,
			}
			result = append(result, stockWithAnalysis)
			stockIndex = len(result) - 1
			stocksMap[stock.ID] = stockIndex
		}

		// Add analysis if it exists
		if analysisID.Valid {
			analysis.ID = int(analysisID.Int64)
			analysis.StockID = stock.ID
			analysis.AnalysisDate = analysisDate.Time
			analysis.CreatedAt = analysisCreatedAt.Time
			result[stockIndex].LatestAnalysis = append(result[stockIndex].LatestAnalysis, analysis)
		}
	}

	meta := models.PaginationMeta{
		Page:        page,
		PageSize:    pageSize,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	return &models.PaginatedResponse[models.StockWithAnalysis]{
		Data: result,
		Meta: meta,
	}, nil
}

func (r *StockRepository) DeleteOldAnalysis(stockID int, keepCount int) error {
	query := `
		DELETE FROM stock_analysis 
		WHERE stock_id = $1 
		AND id NOT IN (
			SELECT id FROM stock_analysis 
			WHERE stock_id = $1 
			ORDER BY analysis_date DESC 
			LIMIT $2
		)`

	_, err := r.db.Exec(query, stockID, keepCount)
	return err
}

func (r *StockRepository) GetFilterOptions() (*models.FilterOptions, error) {
	// Get distinct action types
	actionQuery := `SELECT DISTINCT action FROM stock_analysis WHERE action IS NOT NULL AND action != '' ORDER BY action`
	actionRows, err := r.db.Query(actionQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct actions: %w", err)
	}
	defer actionRows.Close()

	var actionTypes []models.FilterOption
	actionTypes = append(actionTypes, models.FilterOption{Label: "All actions", Value: "all"})

	for actionRows.Next() {
		var action string
		if err := actionRows.Scan(&action); err != nil {
			continue
		}

		// Create label based on action type
		label := formatActionLabel(action)
		value := normalizeActionValue(action)

		// Check if we already have this normalized value
		exists := false
		for _, existing := range actionTypes {
			if existing.Value == value {
				exists = true
				break
			}
		}

		if !exists {
			actionTypes = append(actionTypes, models.FilterOption{
				Label: label,
				Value: value,
			})
		}
	}

	// Get distinct brokerages
	brokerageQuery := `SELECT DISTINCT brokerage FROM stock_analysis WHERE brokerage IS NOT NULL AND brokerage != '' ORDER BY brokerage`
	brokerageRows, err := r.db.Query(brokerageQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct brokerages: %w", err)
	}
	defer brokerageRows.Close()

	var brokerages []models.FilterOption
	brokerages = append(brokerages, models.FilterOption{Label: "All brokerages", Value: "all"})

	for brokerageRows.Next() {
		var brokerage string
		if err := brokerageRows.Scan(&brokerage); err != nil {
			continue
		}
		brokerages = append(brokerages, models.FilterOption{
			Label: brokerage,
			Value: brokerage,
		})
	}

	// Static sort options
	sortOptions := []models.FilterOption{
		{Label: "Newest", Value: "newest"},
		{Label: "Oldest", Value: "oldest"},
		{Label: "Ticker A-Z", Value: "ticker-a-z"},
		{Label: "Company A-Z", Value: "company-a-z"},
		{Label: "Analysis Date (Newest)", Value: "analysis-newest"},
		{Label: "Analysis Date (Oldest)", Value: "analysis-oldest"},
	}

	return &models.FilterOptions{
		ActionTypes: actionTypes,
		Brokerages:  brokerages,
		SortBy:      sortOptions,
	}, nil
}

// Helper function to format action labels
func formatActionLabel(action string) string {
	action = strings.ToLower(action)
	if strings.Contains(action, "initiated") {
		return "Initiated"
	} else if strings.Contains(action, "raised") {
		return "Target Raised"
	} else if strings.Contains(action, "lowered") {
		return "Target Lowered"
	} else if strings.Contains(action, "upgraded") {
		return "Upgraded"
	} else if strings.Contains(action, "downgraded") {
		return "Downgraded"
	} else if strings.Contains(action, "reiterated") {
		return "Reiterated"
	} else if strings.Contains(action, "target set") {
		return "Target Set"
	}
	// Capitalize first letter of each word
	words := strings.Fields(action)
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

// Helper function to normalize action values for filtering
func normalizeActionValue(action string) string {
	action = strings.ToLower(action)
	if strings.Contains(action, "initiated") {
		return "initiated"
	} else if strings.Contains(action, "raised") {
		return "raised"
	} else if strings.Contains(action, "lowered") {
		return "lowered"
	} else if strings.Contains(action, "upgraded") {
		return "upgraded"
	} else if strings.Contains(action, "downgraded") {
		return "downgraded"
	} else if strings.Contains(action, "reiterated") {
		return "reiterated"
	} else if strings.Contains(action, "target set") {
		return "target-set"
	}
	return strings.ReplaceAll(strings.ToLower(action), " ", "-")
}

func (r *StockRepository) GetMarketIntelligenceOverview() (*models.MarketIntelligenceOverview, error) {
	overview := &models.MarketIntelligenceOverview{}

	// Get total stocks
	var totalStocks int
	err := r.db.QueryRow("SELECT COUNT(*) FROM stocks").Scan(&totalStocks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stocks: %w", err)
	}
	overview.TotalStocks = totalStocks

	// Get total analysis count (last 30 days)
	thirtyDaysAgo := "NOW() - INTERVAL '30 days'"
	var recentAnalysis int
	recentAnalysisQuery := fmt.Sprintf("SELECT COUNT(*) FROM stock_analysis WHERE created_at >= %s", thirtyDaysAgo)
	err = r.db.QueryRow(recentAnalysisQuery).Scan(&recentAnalysis)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent analysis: %w", err)
	}
	overview.RecentAnalysis = recentAnalysis

	// Count upgrades and downgrades (last 30 days)
	var upgrades, downgrades int
	upgradeQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM stock_analysis 
		WHERE created_at >= %s 
		AND (LOWER(action) LIKE '%%raised%%' OR LOWER(action) LIKE '%%upgrade%%' OR LOWER(action) LIKE '%%initiated%%' 
		     OR LOWER(rating_to) LIKE '%%buy%%' OR LOWER(rating_to) LIKE '%%outperform%%')`, thirtyDaysAgo)
	err = r.db.QueryRow(upgradeQuery).Scan(&upgrades)
	if err != nil {
		return nil, fmt.Errorf("failed to get upgrades: %w", err)
	}
	overview.Upgrades = upgrades

	downgradeQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM stock_analysis 
		WHERE created_at >= %s 
		AND (LOWER(action) LIKE '%%lowered%%' OR LOWER(action) LIKE '%%downgrade%%' 
		     OR LOWER(rating_to) LIKE '%%sell%%' OR LOWER(rating_to) LIKE '%%underperform%%')`, thirtyDaysAgo)
	err = r.db.QueryRow(downgradeQuery).Scan(&downgrades)
	if err != nil {
		return nil, fmt.Errorf("failed to get downgrades: %w", err)
	}
	overview.Downgrades = downgrades

	// Get top brokerages (last 30 days)
	brokerageQuery := fmt.Sprintf(`
		SELECT brokerage, COUNT(*) as analysis_count
		FROM stock_analysis 
		WHERE created_at >= %s AND brokerage IS NOT NULL AND brokerage != ''
		GROUP BY brokerage 
		ORDER BY analysis_count DESC 
		LIMIT 5`, thirtyDaysAgo)

	brokerageRows, err := r.db.Query(brokerageQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get top brokerages: %w", err)
	}
	defer brokerageRows.Close()

	var topBrokerages []models.BrokerageAnalytics
	totalBrokerageAnalysis := 0
	for brokerageRows.Next() {
		var brokerage models.BrokerageAnalytics
		err := brokerageRows.Scan(&brokerage.Brokerage, &brokerage.AnalysisCount)
		if err != nil {
			continue
		}
		topBrokerages = append(topBrokerages, brokerage)
		totalBrokerageAnalysis += brokerage.AnalysisCount
	}

	// Calculate percentages for brokerages
	for i := range topBrokerages {
		if totalBrokerageAnalysis > 0 {
			topBrokerages[i].Percentage = float64(topBrokerages[i].AnalysisCount) / float64(totalBrokerageAnalysis) * 100
		}
	}
	overview.TopBrokerages = topBrokerages

	// Get top action types (last 30 days)
	actionQuery := fmt.Sprintf(`
		SELECT 
			CASE 
				WHEN LOWER(action) LIKE '%%initiated%%' THEN 'Initiated'
				WHEN LOWER(action) LIKE '%%raised%%' THEN 'Target Raised'
				WHEN LOWER(action) LIKE '%%lowered%%' THEN 'Target Lowered'
				WHEN LOWER(action) LIKE '%%upgraded%%' THEN 'Upgraded'
				WHEN LOWER(action) LIKE '%%downgraded%%' THEN 'Downgraded'
				WHEN LOWER(action) LIKE '%%reiterated%%' THEN 'Reiterated'
				WHEN LOWER(action) LIKE '%%target set%%' THEN 'Target Set'
				ELSE 'Other'
			END as action_type,
			COUNT(*) as count
		FROM stock_analysis 
		WHERE created_at >= %s AND action IS NOT NULL AND action != ''
		GROUP BY action_type 
		ORDER BY count DESC 
		LIMIT 5`, thirtyDaysAgo)

	actionRows, err := r.db.Query(actionQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get top action types: %w", err)
	}
	defer actionRows.Close()

	var topActionTypes []models.ActionTypeAnalytics
	totalActionAnalysis := 0
	for actionRows.Next() {
		var actionType models.ActionTypeAnalytics
		err := actionRows.Scan(&actionType.ActionType, &actionType.Count)
		if err != nil {
			continue
		}
		topActionTypes = append(topActionTypes, actionType)
		totalActionAnalysis += actionType.Count
	}

	// Calculate percentages for action types
	for i := range topActionTypes {
		if totalActionAnalysis > 0 {
			topActionTypes[i].Percentage = float64(topActionTypes[i].Count) / float64(totalActionAnalysis) * 100
		}
	}
	overview.TopActionTypes = topActionTypes

	// Get recent activity trend (last 7 days)
	trendQuery := `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM stock_analysis 
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
		LIMIT 7`

	trendRows, err := r.db.Query(trendQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity trend: %w", err)
	}
	defer trendRows.Close()

	var activityTrend []models.ActivityTrendPoint
	for trendRows.Next() {
		var trend models.ActivityTrendPoint
		err := trendRows.Scan(&trend.Date, &trend.Count)
		if err != nil {
			continue
		}
		activityTrend = append(activityTrend, trend)
	}
	overview.RecentActivityTrend = activityTrend

	return overview, nil
}
