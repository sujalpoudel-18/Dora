package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFileUnsafeSQL(t *testing.T) {
	tmpDir := t.TempDir()
	sampleFile := filepath.Join(tmpDir, "main.go")

	content := `package main
import "database/sql"
import "fmt"

func getUser(db *sql.DB, input string) {
	query := "SELECT * FROM users WHERE username = '" + input + "'"
	db.Query(query)
	db.Exec(fmt.Sprintf("DELETE FROM users WHERE id = %s", input))
}
`
	if err := os.WriteFile(sampleFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	findings, err := ScanFile(sampleFile, "Go")
	if err != nil {
		t.Fatalf("ScanFile failed: %v", err)
	}

	if len(findings) == 0 {
		t.Errorf("Expected findings for unsafe dynamic SQL, got 0")
	}
}
