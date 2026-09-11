package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

// User represents a registered user
type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	IsPro     bool
	UsageMin  int
}

// AnonymousSession tracks usage for anonymous users via cookie
type AnonymousSession struct {
	ID         string
	CookieHash string
	IP         string
	UsageMin   int
	CreatedAt  time.Time
}

// NewDB creates a new database connection with schema
func NewDB(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Create tables if not exist
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		password TEXT,
		created_at INTEGER NOT NULL,
		is_pro INTEGER DEFAULT 0,
		usage_minutes INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		cookie_hash TEXT,
		ip TEXT,
		usage_minutes INTEGER DEFAULT 0,
		created_at INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		cookie_hash TEXT,
		file_name TEXT,
		file_size INTEGER,
		status TEXT DEFAULT 'pending',
		result TEXT,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// CreateUser creates a new user
func (s *Store) CreateUser(id, email, password string) error {
	now := time.Now().Unix()
	_, err := s.db.Exec(
		"INSERT INTO users (id, email, password, created_at) VALUES (?, ?, ?, ?)",
		id, email, password, now,
	)
	return err
}

// GetUserByEmail finds a user by email
func (s *Store) GetUserByEmail(email string) (*User, error) {
	var u User
	var createdAt int64
	err := s.db.QueryRow(
		"SELECT id, email, password, created_at, is_pro, usage_minutes FROM users WHERE email = ?",
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &createdAt, &u.IsPro, &u.UsageMin)
	if err != nil {
		return nil, err
	}
	u.CreatedAt = time.Unix(createdAt, 0)
	return &u, nil
}

// GetUserByID finds a user by ID
func (s *Store) GetUserByID(id string) (*User, error) {
	var u User
	var createdAt int64
	err := s.db.QueryRow(
		"SELECT id, email, password, created_at, is_pro, usage_minutes FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Email, &u.Password, &createdAt, &u.IsPro, &u.UsageMin)
	if err != nil {
		return nil, err
	}
	u.CreatedAt = time.Unix(createdAt, 0)
	return &u, nil
}

// IncrementUserUsage increments user's used minutes
func (s *Store) IncrementUserUsage(id string, minutes int) error {
	_, err := s.db.Exec(
		"UPDATE users SET usage_minutes = usage_minutes + ? WHERE id = ?",
		minutes, id,
	)
	return err
}

// GetOrCreateSession gets or creates an anonymous session
func (s *Store) GetOrCreateSession(cookieHash, ip string) (*AnonymousSession, error) {
	var sess AnonymousSession
	var createdAt int64
	err := s.db.QueryRow(
		"SELECT id, cookie_hash, ip, usage_minutes, created_at FROM sessions WHERE cookie_hash = ?",
		cookieHash,
	).Scan(&sess.ID, &sess.CookieHash, &sess.IP, &sess.UsageMin, &createdAt)

	if err == sql.ErrNoRows {
		// Create new session
		sess = AnonymousSession{
			ID:         uuid.New().String(),
			CookieHash: cookieHash,
			IP:         ip,
			UsageMin:   0,
			CreatedAt:  time.Now(),
		}
		_, err = s.db.Exec(
			"INSERT INTO sessions (id, cookie_hash, ip, usage_minutes, created_at) VALUES (?, ?, ?, ?, ?)",
			sess.ID, sess.CookieHash, sess.IP, sess.UsageMin, sess.CreatedAt.Unix(),
		)
		if err != nil {
			return nil, err
		}
		return &sess, nil
	}

	sess.CreatedAt = time.Unix(createdAt, 0)
	return &sess, err
}

// IncrementSessionUsage increments anonymous session's used minutes
func (s *Store) IncrementSessionUsage(cookieHash string, minutes int) error {
	_, err := s.db.Exec(
		"UPDATE sessions SET usage_minutes = usage_minutes + ? WHERE cookie_hash = ?",
		minutes, cookieHash,
	)
	return err
}

// FreeQuota returns the free quota in minutes
const FreeQuota = 30

// CheckQuota checks if user has remaining quota
func (s *Store) CheckQuota(userID string, cookieHash string) (bool, int, error) {
	if userID != "" {
		user, err := s.GetUserByID(userID)
		if err != nil {
			return false, 0, err
		}
		if user.IsPro {
			return true, -1, nil // unlimited for pro
		}
		return user.UsageMin < FreeQuota, FreeQuota - user.UsageMin, nil
	}

	if cookieHash != "" {
		sess, err := s.GetOrCreateSession(cookieHash, "")
		if err != nil {
			return false, 0, err
		}
		return sess.UsageMin < FreeQuota, FreeQuota - sess.UsageMin, nil
	}

	return false, 0, nil
}

// CreateJob creates a new audio processing job
func (s *Store) CreateJob(id, userID, cookieHash, fileName string, fileSize int64) error {
	now := time.Now().Unix()
	_, err := s.db.Exec(
		"INSERT INTO jobs (id, user_id, cookie_hash, file_name, file_size, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 'pending', ?, ?)",
		id, userID, cookieHash, fileName, fileSize, now, now,
	)
	return err
}

// UpdateJobStatus updates job status and result
func (s *Store) UpdateJobStatus(id, status, result string) error {
	_, err := s.db.Exec(
		"UPDATE jobs SET status = ?, result = ?, updated_at = ? WHERE id = ?",
		status, result, time.Now().Unix(), id,
	)
	return err
}

// GetJob retrieves a job by ID
func (s *Store) GetJob(id string) (map[string]interface{}, error) {
	var job struct {
		ID        string `json:"id"`
		UserID    string `json:"user_id"`
		FileName  string `json:"file_name"`
		FileSize  int64  `json:"file_size"`
		Status    string `json:"status"`
		Result    string `json:"result"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
	}
	// Initialize result to empty string to handle NULL
	var result sql.NullString
	err := s.db.QueryRow(
		"SELECT id, user_id, file_name, file_size, status, result, created_at, updated_at FROM jobs WHERE id = ?",
		id,
	).Scan(&job.ID, &job.UserID, &job.FileName, &job.FileSize, &job.Status, &result, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}
	job.Result = result.String
	return map[string]interface{}{
		"id":         job.ID,
		"user_id":    job.UserID,
		"file_name":  job.FileName,
		"file_size":  job.FileSize,
		"status":     job.Status,
		"result":     job.Result,
		"created_at": job.CreatedAt,
		"updated_at": job.UpdatedAt,
	}, nil
}
