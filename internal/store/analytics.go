package store

import (
	"database/sql"
	"time"
)

// AnalyticsData holds statistical data for dashboard
type AnalyticsData struct {
	TotalJobs      int       `json:"total_jobs"`
	CompletedJobs  int       `json:"completed_jobs"`
	FailedJobs     int       `json:"failed_jobs"`
	AvgDurationSec float64   `json:"avg_duration_sec"`
	RecentJobs     []JobStat `json:"recent_jobs"`
}

// JobStat represents a recent job summary
type JobStat struct {
	ID          string  `json:"id"`
	FileName    string  `json:"file_name"`
	Status      string  `json:"status"`
	CreatedAt   int64   `json:"created_at"`
	DurationSec float64 `json:"duration_sec"`
}

// GetAnalytics returns usage statistics for the given user
func (s *Store) GetAnalytics(userID string, days int) (*AnalyticsData, error) {
	data := &AnalyticsData{}

	// Total jobs
	err := s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = ?", userID).Scan(&data.TotalJobs)
	if err != nil {
		return nil, err
	}

	// Completed jobs
	err = s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = ? AND status = 'completed'", userID).Scan(&data.CompletedJobs)
	if err != nil {
		return nil, err
	}

	// Failed jobs
	err = s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = ? AND status = 'failed'", userID).Scan(&data.FailedJobs)
	if err != nil {
		return nil, err
	}

	// Average duration (updated_at - created_at)
	var avgDur sql.NullFloat64
	err = s.db.QueryRow(`
		SELECT AVG(STRFTIME('%s', datetime(updated_at, 'unixepoch')) - 
		         STRFTIME('%s', datetime(created_at, 'unixepoch')))
		FROM jobs WHERE user_id = ? AND status = 'completed'
	`, userID).Scan(&avgDur)
	if err == nil && avgDur.Valid {
		data.AvgDurationSec = avgDur.Float64
	}

	// Recent jobs (last N days)
	startDate := time.Now().AddDate(0, 0, -days)
	rows, err := s.db.Query(`
		SELECT id, file_name, status, created_at, updated_at
		FROM jobs 
		WHERE user_id = ? AND created_at >= ?
		ORDER BY created_at DESC
		LIMIT 10
	`, userID, startDate.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var j JobStat
		var createdAt, updatedAt int64
		if err := rows.Scan(&j.ID, &j.FileName, &j.Status, &createdAt, &updatedAt); err != nil {
			continue
		}
		j.CreatedAt = createdAt
		j.DurationSec = float64(updatedAt - createdAt)
		data.RecentJobs = append(data.RecentJobs, j)
	}

	return data, nil
}