package template_test

import (
	"testing"

	t_ "design_patterns/template"
)

// --- CSVMigration ---

func TestCSVMigration_Success(t *testing.T) {
	m := t_.NewCSVMigration("data.csv", []string{
		"ID, Name, City",
		"1, Alice, NYC",
		"2, Bob, LA",
	})
	if err := t_.RunMigration(m); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCSVMigration_TransformLowercasesValues(t *testing.T) {
	m := t_.NewCSVMigration("data.csv", []string{
		"Name",
		"ALICE",
	})
	records, _ := m.ExtractData()
	transformed, err := m.TransformData(records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transformed[0]["Name"] != "alice" {
		t.Errorf("expected lowercase 'alice', got %q", transformed[0]["Name"])
	}
}

func TestCSVMigration_EmptySource(t *testing.T) {
	m := t_.NewCSVMigration("", []string{"ID", "1"})
	if err := t_.RunMigration(m); err == nil {
		t.Fatal("expected connect error for empty source, got nil")
	}
}

func TestCSVMigration_NoRows(t *testing.T) {
	m := t_.NewCSVMigration("data.csv", []string{})
	if err := t_.RunMigration(m); err == nil {
		t.Fatal("expected extract error for empty rows, got nil")
	}
}

// --- JSONMigration ---

func TestJSONMigration_Success(t *testing.T) {
	m := t_.NewJSONMigration("https://api.example.com/data", []t_.Record{
		{"id": "1", "name": " alice "},
		{"id": "2", "name": " bob "},
	})
	if err := t_.RunMigration(m); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestJSONMigration_TransformTrimsWhitespace(t *testing.T) {
	m := t_.NewJSONMigration("https://api.example.com/data", []t_.Record{
		{"name": "  alice  "},
	})
	records, _ := m.ExtractData()
	transformed, err := m.TransformData(records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transformed[0]["name"] != "alice" {
		t.Errorf("expected trimmed 'alice', got %q", transformed[0]["name"])
	}
}

func TestJSONMigration_EmptyEndpoint(t *testing.T) {
	m := t_.NewJSONMigration("", []t_.Record{{"id": "1"}})
	if err := t_.RunMigration(m); err == nil {
		t.Fatal("expected connect error for empty endpoint, got nil")
	}
}

func TestJSONMigration_EmptyPayload(t *testing.T) {
	m := t_.NewJSONMigration("https://api.example.com/data", []t_.Record{})
	if err := t_.RunMigration(m); err == nil {
		t.Fatal("expected extract error for empty payload, got nil")
	}
}
