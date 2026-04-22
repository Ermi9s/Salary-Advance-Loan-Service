package dto

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"salaryAdvance/internal/entity"
)

type DTOMethodes struct {
	CustomersFilePath      string
	TransactionFilePath    string
	SampleCustomerFilePath string
}

func readJSONFile(path string, target any) error {
	openedFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = openedFile.Close()
	}()

	decoder := json.NewDecoder(openedFile)
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}

func (d *DTOMethodes) ReadCustomerData() ([]entity.BankCustomer, error) {
	var records []entity.BankCustomer
	if err := readJSONFile(d.CustomersFilePath, &records); err != nil {
		log.Printf("Error Reading data from Json %v", err)
		return nil, err
	}
	return records, nil
}

func (d *DTOMethodes) ReadTransactionData() ([]entity.Transaction, error) {
	var records []entity.Transaction
	if err := readJSONFile(d.TransactionFilePath, &records); err != nil {
		log.Printf("Error Reading data from Json %v", err)
		return nil, err
	}
	return records, nil
}

func (d *DTOMethodes) ReadSampleData() ([]entity.InputCustomer, error) {
	return d.ReadCSV()
}

func (d *DTOMethodes) ReadCSV() ([]entity.InputCustomer, error) {
	openedFile, err := os.Open(d.SampleCustomerFilePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := openedFile.Close(); closeErr != nil {
			log.Printf("Error closing file: %v", closeErr)
		}
	}()

	reader := csv.NewReader(openedFile)
	var customers []entity.InputCustomer

	// skip header
	if _, err := reader.Read(); err != nil {
		log.Printf("Error reading header: %v", err)
		return nil, err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record: %v", err)
			return nil, err
		}

		// safety check
		if len(record) < 2 {
			log.Printf("Invalid record: %v", record)
			continue
		}

		customer := entity.InputCustomer{
			CustomerName: record[0],
			AccountNo:    record[1],
		}
		customers = append(customers, customer)
	}

	return customers, nil
}
