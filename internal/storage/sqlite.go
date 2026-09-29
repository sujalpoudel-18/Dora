package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dora/dora/pkg/models"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

// NewStore initializes a SQLite database file and creates required schema tables
func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	s := &Store{db: db}
	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

func (s *Store) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS scans (
		id TEXT PRIMARY KEY,
		created_at DATETIME,
		status TEXT
	);

	CREATE TABLE IF NOT EXISTS findings (
		id TEXT PRIMARY KEY,
		scan_id TEXT,
		template_id TEXT,
		name TEXT,
		severity TEXT,
		confidence TEXT,
		target TEXT,
		endpoint TEXT,
		source_file TEXT,
		line_number INTEGER,
		description TEXT,
		remediation TEXT,
		timestamp DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_findings_severity ON findings(severity);
	CREATE INDEX IF NOT EXISTS idx_findings_target ON findings(target);
	`
	_, err := s.db.Exec(query)
	return err
}

// SaveFinding inserts a vulnerability finding into SQLite
func (s *Store) SaveFinding(scanID string, f *models.Finding) error {
	query := `
	INSERT OR REPLACE INTO findings (id, scan_id, template_id, name, severity, confidence, target, endpoint, source_file, line_number, description, remediation, timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := s.db.Exec(query, f.ID, scanID, f.TemplateID, f.Name, string(f.Severity), string(f.Confidence), f.Target, f.Endpoint, f.SourceFile, f.LineNumber, f.Description, f.Remediation, time.Now())
	return err
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}
