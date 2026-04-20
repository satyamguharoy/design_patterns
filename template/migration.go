// Package template demonstrates the Template Method pattern in Go.
//
// Template Method defines the skeleton of an algorithm in one place and
// lets implementations fill in the variable steps — without changing the order
// or structure of those steps.
//
// Classical OOP achieves this via inheritance and method overriding.
// Go has no inheritance, so the pattern uses an interface for the variable
// steps and a standalone function as the fixed template:
//
//	func RunMigration(m Migrator) error { … fixed order … }
//
// This is strictly cleaner: the template is not buried in a base struct,
// and implementations cannot accidentally bypass it.
package template

import (
	"errors"
	"fmt"
	"strings"
)

// Record is a generic row of key-value data produced by extraction.
type Record map[string]string

// Migrator defines the variable steps of a data migration.
// Each implementation supplies its own logic for these steps;
// RunMigration controls their order.
type Migrator interface {
	Connect() error
	ExtractData() ([]Record, error)
	TransformData(records []Record) ([]Record, error)
	LoadData(records []Record) error
	Disconnect()
	Name() string
}

// RunMigration is the template. It enforces the algorithm skeleton —
// callers and implementations never reorder or skip steps.
func RunMigration(m Migrator) error {
	fmt.Printf("[%s] starting migration\n", m.Name())

	if err := m.Connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer m.Disconnect()

	records, err := m.ExtractData()
	if err != nil {
		return fmt.Errorf("extract: %w", err)
	}
	fmt.Printf("[%s] extracted %d records\n", m.Name(), len(records))

	transformed, err := m.TransformData(records)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}
	fmt.Printf("[%s] transformed to %d records\n", m.Name(), len(transformed))

	if err := m.LoadData(transformed); err != nil {
		return fmt.Errorf("load: %w", err)
	}

	fmt.Printf("[%s] migration complete\n", m.Name())
	return nil
}

// --- Concrete implementations ---

// CSVMigration reads comma-separated rows and normalises field names to lowercase.
type CSVMigration struct {
	source string
	rows   []string // simulated raw CSV lines
}

func NewCSVMigration(source string, rows []string) *CSVMigration {
	return &CSVMigration{source: source, rows: rows}
}

func (c *CSVMigration) Name() string { return "CSVMigration" }

func (c *CSVMigration) Connect() error {
	if c.source == "" {
		return errors.New("csv source path is empty")
	}
	fmt.Printf("[CSV] connected to %s\n", c.source)
	return nil
}

func (c *CSVMigration) ExtractData() ([]Record, error) {
	if len(c.rows) == 0 {
		return nil, errors.New("csv: no data rows")
	}
	headers := strings.Split(c.rows[0], ",")
	var records []Record
	for _, row := range c.rows[1:] {
		fields := strings.Split(row, ",")
		r := make(Record)
		for i, h := range headers {
			if i < len(fields) {
				r[strings.TrimSpace(h)] = strings.TrimSpace(fields[i])
			}
		}
		records = append(records, r)
	}
	return records, nil
}

// TransformData lowercases all values — the CSV-specific transformation.
func (c *CSVMigration) TransformData(records []Record) ([]Record, error) {
	for _, r := range records {
		for k, v := range r {
			r[k] = strings.ToLower(v)
		}
	}
	return records, nil
}

func (c *CSVMigration) LoadData(records []Record) error {
	for _, r := range records {
		fmt.Printf("[CSV] loading record: %v\n", r)
	}
	return nil
}

func (c *CSVMigration) Disconnect() {
	fmt.Println("[CSV] disconnected")
}

// JSONMigration reads pre-parsed JSON-like records and trims whitespace from values.
type JSONMigration struct {
	endpoint string
	payload  []Record // simulated parsed JSON
}

func NewJSONMigration(endpoint string, payload []Record) *JSONMigration {
	return &JSONMigration{endpoint: endpoint, payload: payload}
}

func (j *JSONMigration) Name() string { return "JSONMigration" }

func (j *JSONMigration) Connect() error {
	if j.endpoint == "" {
		return errors.New("json endpoint is empty")
	}
	fmt.Printf("[JSON] connected to %s\n", j.endpoint)
	return nil
}

func (j *JSONMigration) ExtractData() ([]Record, error) {
	if len(j.payload) == 0 {
		return nil, errors.New("json: empty payload")
	}
	return j.payload, nil
}

// TransformData trims whitespace from all values — the JSON-specific transformation.
func (j *JSONMigration) TransformData(records []Record) ([]Record, error) {
	for _, r := range records {
		for k, v := range r {
			r[k] = strings.TrimSpace(v)
		}
	}
	return records, nil
}

func (j *JSONMigration) LoadData(records []Record) error {
	for _, r := range records {
		fmt.Printf("[JSON] loading record: %v\n", r)
	}
	return nil
}

func (j *JSONMigration) Disconnect() {
	fmt.Println("[JSON] disconnected")
}
