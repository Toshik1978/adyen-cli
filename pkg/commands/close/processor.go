package close

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/gocarina/gocsv"
	"go.uber.org/zap"

	"github.com/Toshik1978/csv2adyen/pkg/adyen"
	"github.com/Toshik1978/csv2adyen/pkg/commands"
)

var (
	// ErrInvalidResponse means we have a wrong Adyen response.
	ErrInvalidResponse = errors.New("store details count not equal to accounts count")
)

// Processor declare implementation of the main module.
type Processor struct {
	logger      *zap.Logger
	client      *http.Client
	adyenAPI    *adyen.API
	csvFilePath string
	dryRun      bool
}

// New creates new instance of Processor.
func New(
	logger *zap.Logger, client *http.Client, config *commands.Config,
	csvFilePath string, production, dryRun bool) *Processor {
	var apiURL, apiKey string
	switch {
	case production:
		apiURL = config.AdyenCalURL
		apiKey = config.AdyenCalKey
	case !production:
		apiURL = config.AdyenCalTestURL
		apiKey = config.AdyenCalTestKey
	}

	gocsv.SetHeaderNormalizer(strings.ToUpper)

	return &Processor{
		logger:      logger,
		client:      client,
		adyenAPI:    adyen.New(logger, client, apiURL, apiKey),
		csvFilePath: csvFilePath,
		dryRun:      dryRun,
	}
}

// Run runs parsing & split config updating.
func (p *Processor) Run(ctx context.Context) error {
	file, err := os.OpenFile(p.csvFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open CSV: %w", err)
	}
	defer file.Close()

	var records []*Record
	if err := gocsv.UnmarshalFile(file, &records); err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	var successCnt int
	var failureCnt int
	errs := make([]error, 0, 2)
	for _, record := range records {
		if err := p.process(ctx, record); err != nil {
			failureCnt++
			errs = append(errs, err)
		} else {
			successCnt++
		}
	}

	if failureCnt > 0 {
		p.logger.
			With(zap.Errors("Errors", errs)).
			With(zap.Int("Success Count", successCnt)).
			With(zap.Int("Failure Count", failureCnt)).
			Error("Failed to process restaurants")
		return fmt.Errorf("failed to process close record: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process restaurants")
	return nil
}

func (p *Processor) process(ctx context.Context, record *Record) error {
	accountHolder, err := p.adyenAPI.AccountHolder(ctx, record.AccountHolderCode)
	if err != nil {
		return fmt.Errorf("failed to get account holder: %w", err)
	}
	if err := p.updateAccountHolder(accountHolder, record.StoreID); err != nil {
		return fmt.Errorf("failed to close account: %w", err)
	}
	if !p.dryRun {
		accountHolder.AccountHolderDetails.StoreDetails[0].Status = pointer.ToString("InactiveWithModifications")
		if err := p.adyenAPI.UpdateAccountHolder(ctx, &accountHolder.UpdateAccountHolderRequest); err != nil {
			return fmt.Errorf("failed to update account holder: %w", err)
		}

		accountHolder.AccountHolderDetails.StoreDetails[0].Status = pointer.ToString("Closed")
		if err := p.adyenAPI.UpdateAccountHolder(ctx, &accountHolder.UpdateAccountHolderRequest); err != nil {
			return fmt.Errorf("failed to update account holder: %w", err)
		}

		if err := p.adyenAPI.CloseAccountHolder(ctx, record.AccountHolderCode); err != nil {
			return fmt.Errorf("failed to close account holder: %w", err)
		}
	}
	return nil
}

func (p *Processor) updateAccountHolder(accountHolder *adyen.GetAccountHolderResponse, storeID string) error {
	if len(accountHolder.AccountHolderDetails.StoreDetails) != len(accountHolder.Accounts) {
		return ErrInvalidResponse
	}
	if len(accountHolder.AccountHolderDetails.StoreDetails) != 1 {
		return ErrInvalidResponse
	}
	if accountHolder.AccountHolderDetails.StoreDetails[0].StoreID != storeID {
		return fmt.Errorf("store ID not found: %s %s", accountHolder.AccountHolderDetails.StoreDetails[0].StoreID, storeID)
	}

	accountHolder.AccountHolderDetails.StoreDetails[0].VirtualAccount = nil
	accountHolder.AccountHolderDetails.StoreDetails[0].SplitConfigurationUUID = nil
	return nil
}
