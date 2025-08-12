-- Initialize the stock database
-- This script runs automatically when the database starts

-- Note: Database should already be selected via connection string
-- CREATE DATABASE IF NOT EXISTS stockdb;
-- \c stockdb;

-- Create stocks table
CREATE TABLE IF NOT EXISTS stocks (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create stock_analysis table
CREATE TABLE IF NOT EXISTS stock_analysis (
    id SERIAL PRIMARY KEY,
    stock_id INT REFERENCES stocks(id),
    target_from VARCHAR(20),
    target_to VARCHAR(20),
    action VARCHAR(100),
    brokerage VARCHAR(100),
    rating_from VARCHAR(50),
    rating_to VARCHAR(50),
    analysis_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create process_control table for managing background processes
CREATE TABLE IF NOT EXISTS process_control (
    id SERIAL PRIMARY KEY,
    process_name VARCHAR(100) NOT NULL UNIQUE,
    is_running BOOLEAN NOT NULL DEFAULT FALSE,
    last_execution TIMESTAMP,
    interval_minutes INT NOT NULL DEFAULT 60,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create recommendation_scores table for pre-calculated scores
CREATE TABLE IF NOT EXISTS recommendation_scores (
    id SERIAL PRIMARY KEY,
    stock_id INT REFERENCES stocks(id) ON DELETE CASCADE,
    total_score DECIMAL(10,2) NOT NULL DEFAULT 0,
    rating_score DECIMAL(10,2) NOT NULL DEFAULT 0,
    rating_change_score DECIMAL(10,2) NOT NULL DEFAULT 0,
    target_change_score DECIMAL(10,2) NOT NULL DEFAULT 0,
    action_score DECIMAL(10,2) NOT NULL DEFAULT 0,
    coverage_score DECIMAL(10,2) NOT NULL DEFAULT 0,
    confidence VARCHAR(10) NOT NULL DEFAULT 'Low',
    reason TEXT,
    latest_analysis_id INT REFERENCES stock_analysis(id),
    calculated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(stock_id)
);

-- Create unique constraint to prevent duplicate analysis for same stock on same date
CREATE UNIQUE INDEX IF NOT EXISTS idx_stock_analysis_unique ON stock_analysis(stock_id, analysis_date, brokerage);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_stocks_symbol ON stocks(symbol);
CREATE INDEX IF NOT EXISTS idx_stock_analysis_stock_id ON stock_analysis(stock_id);
CREATE INDEX IF NOT EXISTS idx_stock_analysis_date ON stock_analysis(analysis_date);
CREATE INDEX IF NOT EXISTS idx_process_control_name ON process_control(process_name);
CREATE INDEX IF NOT EXISTS idx_recommendation_scores_total_score ON recommendation_scores(total_score DESC);
CREATE INDEX IF NOT EXISTS idx_recommendation_scores_stock_id ON recommendation_scores(stock_id);
CREATE INDEX IF NOT EXISTS idx_recommendation_scores_confidence ON recommendation_scores(confidence);

-- Insert process control entries
INSERT INTO process_control (process_name, interval_minutes) VALUES 
('stock_sync', 30)
ON CONFLICT (process_name) DO NOTHING;

-- Insert some sample data for testing (optional)
-- INSERT INTO stocks (symbol, name) VALUES 
-- ('AAPL', 'Apple Inc.'),
-- ('GOOGL', 'Alphabet Inc.'),
-- ('MSFT', 'Microsoft Corporation')
-- ON CONFLICT (symbol) DO NOTHING;

SHOW TABLES;