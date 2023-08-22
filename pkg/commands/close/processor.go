package close

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	logger           *zap.Logger
	client           *http.Client
	adyenAPI         *adyen.API
	csvFilePath      string
	shouldCloseStore bool
	dryRun           bool
}

// New creates new instance of Processor.
func New(
	logger *zap.Logger, client *http.Client, config *commands.Config,
	csvFilePath string, shouldCloseStore, production, dryRun bool) *Processor {
	var calURL, calKey, mgmtURL, mgmtKey string
	switch {
	case production:
		calURL = config.AdyenCalURL
		calKey = config.AdyenCalKey
		mgmtURL = config.AdyenMgmtURL
		mgmtKey = config.AdyenMgmtKey
	case !production:
		calURL = config.AdyenCalTestURL
		calKey = config.AdyenCalTestKey
		mgmtURL = config.AdyenMgmtTestURL
		mgmtKey = config.AdyenMgmtTestKey
	}

	gocsv.SetHeaderNormalizer(strings.ToUpper)

	return &Processor{
		logger:           logger,
		client:           client,
		adyenAPI:         adyen.New(logger, client, calURL, calKey, mgmtURL, mgmtKey),
		csvFilePath:      csvFilePath,
		shouldCloseStore: shouldCloseStore,
		dryRun:           dryRun,
	}
}

// Run runs closer of merchant accounts and stores.
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
	if err := p.checkAccountHolder(accountHolder, record.StoreID); err != nil {
		return fmt.Errorf("failed to close account holder: %w", err)
	}
	if !p.dryRun {
		if p.shouldCloseStore {
			if err := p.closeStore(ctx, record.StoreID); err != nil {
				return fmt.Errorf("failed to close the store (%s): %w", record.StoreID, err)
			}
		}

		if err := p.adyenAPI.CloseAccountHolder(ctx, record.AccountHolderCode); err != nil {
			return fmt.Errorf("failed to close account holder: %w", err)
		}
	}
	return nil
}

func (p *Processor) checkAccountHolder(accountHolder *adyen.GetAccountHolderResponse, storeID string) error {
	if len(accountHolder.AccountHolderDetails.StoreDetails) != len(accountHolder.Accounts) {
		return ErrInvalidResponse
	}
	if len(accountHolder.AccountHolderDetails.StoreDetails) != 1 {
		return ErrInvalidResponse
	}
	if accountHolder.AccountHolderDetails.StoreDetails[0].StoreID != storeID {
		return fmt.Errorf("store ID not found: %s %s", accountHolder.AccountHolderDetails.StoreDetails[0].StoreID, storeID)
	}
	return nil
}

func (p *Processor) closeStore(ctx context.Context, storeID string) error {
	stores, err := p.adyenAPI.GetAllStores(ctx, storeID)
	if err != nil {
		return fmt.Errorf("failed to get all stores: %w", err)
	}
	if stores.ItemsTotal != 1 || len(stores.Data) != 1 {
		return ErrInvalidResponse
	}
	if stores.Data[0].Reference != storeID {
		return fmt.Errorf("store ID not found: %s %s", stores.Data[0].Reference, storeID)
	}

	if err := p.adyenAPI.SetStoreStatus(ctx, stores.Data[0].ID, "inactive"); err != nil {
		return fmt.Errorf("failed to set inactive status: %w", err)
	}
	if err := p.adyenAPI.SetStoreStatus(ctx, stores.Data[0].ID, "closed"); err != nil {
		return fmt.Errorf("failed to set closed status: %w", err)
	}

	return nil
}
