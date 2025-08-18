package repository

import (
	"database/sql"
	"time"

	"stock-api/internal/models"
)

type RecommendationScoreRepository struct {
	db *sql.DB
}

func NewRecommendationScoreRepository(db *sql.DB) *RecommendationScoreRepository {
	return &RecommendationScoreRepository{db: db}
}

func (r *RecommendationScoreRepository) UpsertRecommendationScore(score *models.RecommendationScore) error {
	query := `
		INSERT INTO recommendation_scores (
			stock_id, total_score, rating_score, rating_change_score, 
			target_change_score, action_score, coverage_score, 
			confidence, reason, latest_analysis_id, calculated_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $11)
		ON CONFLICT (stock_id) DO UPDATE SET
			total_score = EXCLUDED.total_score,
			rating_score = EXCLUDED.rating_score,
			rating_change_score = EXCLUDED.rating_change_score,
			target_change_score = EXCLUDED.target_change_score,
			action_score = EXCLUDED.action_score,
			coverage_score = EXCLUDED.coverage_score,
			confidence = EXCLUDED.confidence,
			reason = EXCLUDED.reason,
			latest_analysis_id = EXCLUDED.latest_analysis_id,
			calculated_at = EXCLUDED.calculated_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at`

	now := time.Now()
	score.CalculatedAt = now
	score.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		score.StockID,
		score.TotalScore,
		score.RatingScore,
		score.RatingChangeScore,
		score.TargetChangeScore,
		score.ActionScore,
		score.CoverageScore,
		score.Confidence,
		score.Reason,
		score.LatestAnalysisID,
		now,
	).Scan(&score.ID, &score.CreatedAt)

	return err
}


func (r *RecommendationScoreRepository) GetTopRecommendationsPaginated(page, pageSize int) (*models.PaginatedResponse[models.RecommendationWithStock], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM recommendation_scores`
	var totalItems int
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	totalPages := int(float64(totalItems)/float64(pageSize) + 0.9) // Ceiling division
	if totalPages == 0 {
		totalPages = 1
	}

	offset := (page - 1) * pageSize

	// Get paginated data
	query := `
		SELECT 
			rs.id, rs.stock_id, rs.total_score, rs.rating_score, rs.rating_change_score,
			rs.target_change_score, rs.action_score, rs.coverage_score, rs.confidence,
			rs.reason, rs.latest_analysis_id, rs.calculated_at, rs.created_at, rs.updated_at,
			s.id, s.symbol, s.name, s.created_at, s.updated_at
		FROM recommendation_scores rs
		JOIN stocks s ON rs.stock_id = s.id
		ORDER BY rs.total_score DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []models.RecommendationWithStock
	var stockRepo = NewStockRepository(r.db)

	for rows.Next() {
		var rec models.RecommendationWithStock
		var stock models.Stock

		err := rows.Scan(
			&rec.ID, &rec.StockID, &rec.TotalScore, &rec.RatingScore, &rec.RatingChangeScore,
			&rec.TargetChangeScore, &rec.ActionScore, &rec.CoverageScore, &rec.Confidence,
			&rec.Reason, &rec.LatestAnalysisID, &rec.CalculatedAt, &rec.CreatedAt, &rec.UpdatedAt,
			&stock.ID, &stock.Symbol, &stock.Name, &stock.CreatedAt, &stock.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get latest analysis for this stock
		stockWithAnalysis := models.StockWithAnalysis{Stock: stock}
		analyses, err := stockRepo.GetLatestAnalysisForStock(stock.ID, 5)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		stockWithAnalysis.LatestAnalysis = analyses

		rec.Stock = stockWithAnalysis
		recommendations = append(recommendations, rec)
	}

	meta := models.PaginationMeta{
		Page:        page,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	return &models.PaginatedResponse[models.RecommendationWithStock]{
		Data: recommendations,
		Meta: meta,
	}, nil
}

func (r *RecommendationScoreRepository) DeleteRecommendationScore(stockID int) error {
	query := `DELETE FROM recommendation_scores WHERE stock_id = $1`
	_, err := r.db.Exec(query, stockID)
	return err
}

func (r *RecommendationScoreRepository) GetRecommendationScoreByStockID(stockID int) (*models.RecommendationScore, error) {
	query := `
		SELECT id, stock_id, total_score, rating_score, rating_change_score,
			   target_change_score, action_score, coverage_score, confidence,
			   reason, latest_analysis_id, calculated_at, created_at, updated_at
		FROM recommendation_scores 
		WHERE stock_id = $1`

	var score models.RecommendationScore
	err := r.db.QueryRow(query, stockID).Scan(
		&score.ID, &score.StockID, &score.TotalScore, &score.RatingScore, &score.RatingChangeScore,
		&score.TargetChangeScore, &score.ActionScore, &score.CoverageScore, &score.Confidence,
		&score.Reason, &score.LatestAnalysisID, &score.CalculatedAt, &score.CreatedAt, &score.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &score, nil
}

func (r *RecommendationScoreRepository) GetRecommendationStats() (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_recommendations,
			COUNT(CASE WHEN confidence = 'High' THEN 1 END) as high_confidence,
			COUNT(CASE WHEN confidence = 'Medium' THEN 1 END) as medium_confidence,
			COUNT(CASE WHEN confidence = 'Low' THEN 1 END) as low_confidence,
			AVG(total_score) as avg_score
		FROM recommendation_scores`

	var totalRecs, highConf, mediumConf, lowConf int
	var avgScore sql.NullFloat64

	err := r.db.QueryRow(query).Scan(&totalRecs, &highConf, &mediumConf, &lowConf, &avgScore)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_recommendations": totalRecs,
		"high_confidence":       highConf,
		"medium_confidence":     mediumConf,
		"low_confidence":        lowConf,
		"average_score":         0.0,
	}

	if avgScore.Valid {
		stats["average_score"] = avgScore.Float64
	}

	return stats, nil
}