package repository

import (
	"database/sql"
	"fmt"
	"time"

	"stock-api/internal/models"
)

const (
	ProcessStockSync = "stock_sync"
)

type ProcessControlRepository struct {
	db *sql.DB
}

func NewProcessControlRepository(db *sql.DB) *ProcessControlRepository {
	return &ProcessControlRepository{db: db}
}

func (r *ProcessControlRepository) GetProcessControl(processName string) (*models.ProcessControl, error) {
	query := `
		SELECT id, process_name, is_running, last_execution, interval_minutes, created_at, updated_at
		FROM process_control WHERE process_name = $1`
	
	process := &models.ProcessControl{}
	err := r.db.QueryRow(query, processName).Scan(
		&process.ID, &process.ProcessName, &process.IsRunning, &process.LastExecution,
		&process.IntervalMinutes, &process.CreatedAt, &process.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return process, err
}

func (r *ProcessControlRepository) CanStartProcess(processName string) (bool, error) {
	process, err := r.GetProcessControl(processName)
	if err != nil {
		return false, err
	}
	
	if process == nil {
		return false, fmt.Errorf("process %s not found in control table", processName)
	}
	
	// If process is already running, cannot start
	if process.IsRunning {
		return false, nil
	}
	
	// If no last execution, can start
	if process.LastExecution == nil {
		return true, nil
	}
	
	// Check if enough time has passed since last execution
	now := time.Now()
	nextAllowedExecution := process.LastExecution.Add(time.Duration(process.IntervalMinutes) * time.Minute)
	
	return now.After(nextAllowedExecution), nil
}

func (r *ProcessControlRepository) StartProcess(processName string) error {
	query := `
		UPDATE process_control 
		SET is_running = true, updated_at = NOW()
		WHERE process_name = $1 AND is_running = false`
	
	result, err := r.db.Exec(query, processName)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("process %s is already running or not found", processName)
	}
	
	return nil
}

func (r *ProcessControlRepository) FinishProcess(processName string) error {
	query := `
		UPDATE process_control 
		SET is_running = false, last_execution = NOW(), updated_at = NOW()
		WHERE process_name = $1`
	
	_, err := r.db.Exec(query, processName)
	return err
}

func (r *ProcessControlRepository) ForceStopProcess(processName string) error {
	query := `
		UPDATE process_control 
		SET is_running = false, updated_at = NOW()
		WHERE process_name = $1`
	
	_, err := r.db.Exec(query, processName)
	return err
}

// Convenience methods for specific processes
func (r *ProcessControlRepository) CanStartStockSync() (bool, error) {
	return r.CanStartProcess(ProcessStockSync)
}

func (r *ProcessControlRepository) StartStockSync() error {
	return r.StartProcess(ProcessStockSync)
}

func (r *ProcessControlRepository) FinishStockSync() error {
	return r.FinishProcess(ProcessStockSync)
}

func (r *ProcessControlRepository) ForceStopStockSync() error {
	return r.ForceStopProcess(ProcessStockSync)
}