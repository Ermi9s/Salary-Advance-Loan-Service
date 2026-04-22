package services

import (
	"context"
	"fmt"
	"time"

	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/interfaces/dto"
)

type ValidationService struct {
	Repository    ValidationRepository
	DTOMethodes   *dto.DTOMethodes
	RatingService *RatingService
}

func (s *ValidationService) ValidateData(ctx context.Context) (entity.ValidationReport, error) {
	var emptyReport entity.ValidationReport

	customerList, err := s.DTOMethodes.ReadCustomerData()
	if err != nil {
		return emptyReport, err
	}

	sampleData, err := s.DTOMethodes.ReadSampleData()
	if err != nil {
		return emptyReport, err
	}

	batchValidationLog, verifiedCustomers, err := ProcessData(customerList, sampleData)
	if err != nil {
		return emptyReport, err
	}

	if err := s.Repository.SaveVerifiedCustomers(ctx, verifiedCustomers); err != nil {
		return emptyReport, err
	}

	report := BuildValidationReport(batchValidationLog)
	report.ProcessedAt = time.Now().UTC().Format(time.RFC3339)

	return report, nil
}

func (s *ValidationService) CalculateCustomerRatings(ctx context.Context) ([]entity.CustomerRating, error) {
	_ = ctx

	if s.RatingService == nil {
		return nil, fmt.Errorf("rating service is not configured")
	}

	customerList, err := s.DTOMethodes.ReadCustomerData()
	if err != nil {
		return nil, err
	}

	transactionData, err := s.DTOMethodes.ReadTransactionData()
	if err != nil {
		return nil, err
	}

	return s.RatingService.RateCustomers(customerList, transactionData)
}

func (s *ValidationService) GetVerifiedCustomers(ctx context.Context) ([]entity.BankCustomer, error) {
	return s.Repository.ListVerifiedCustomers(ctx)
}
