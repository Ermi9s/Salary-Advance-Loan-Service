package dto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadCSVParsesValidRows(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "sample.csv")
	content := "name,account\nAlice,1234567890\nBob,9876543210\n"
	if err := os.WriteFile(csvPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	d := &DTOMethodes{SampleCustomerFilePath: csvPath}
	customers, err := d.ReadCSV()
	if err != nil {
		t.Fatalf("ReadCSV returned error: %v", err)
	}

	if len(customers) != 2 {
		t.Fatalf("expected 2 valid records, got %d", len(customers))
	}
	if customers[0].CustomerName != "Alice" || customers[0].AccountNo != "1234567890" {
		t.Fatalf("unexpected first customer: %#v", customers[0])
	}
	if customers[1].CustomerName != "Bob" || customers[1].AccountNo != "9876543210" {
		t.Fatalf("unexpected second customer: %#v", customers[1])
	}
}

func TestReadCSVReturnsErrorForMalformedRow(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "sample.csv")
	content := "name,account\nAlice,1234567890\ninvalid-only-one-column\n"
	if err := os.WriteFile(csvPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	d := &DTOMethodes{SampleCustomerFilePath: csvPath}
	if _, err := d.ReadCSV(); err == nil {
		t.Fatalf("expected malformed CSV row error")
	}
}

func TestReadSampleDataUsesCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "sample.csv")
	content := "name,account\nAlice,1234567890\nBob,9876543210\n"
	if err := os.WriteFile(csvPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	d := &DTOMethodes{SampleCustomerFilePath: csvPath}
	customers, err := d.ReadSampleData()
	if err != nil {
		t.Fatalf("ReadSampleData returned error: %v", err)
	}

	if len(customers) != 2 {
		t.Fatalf("expected 2 customers, got %d", len(customers))
	}
	if customers[0].CustomerName != "Alice" || customers[1].AccountNo != "9876543210" {
		t.Fatalf("unexpected customers: %#v", customers)
	}
}

func TestReadCustomerDataMissingFile(t *testing.T) {
	d := &DTOMethodes{CustomersFilePath: filepath.Join(t.TempDir(), "missing.json")}
	_, err := d.ReadCustomerData()
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}
