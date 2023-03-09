package processor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gocarina/gocsv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Toshik1978/csv2adyen/pkg/adyen"
)

var (
	// ErrInvalidResponse means we have a wrong Adyen response.
	ErrInvalidResponse = errors.New("store details count not equal to accounts count")
	// ErrStoreNotFound means we could not find store in the store details.
	ErrStoreNotFound = errors.New("store ID not found")
)

// Processor declare implementation of the main module.
type Processor struct {
	logger   *zap.Logger
	client   *http.Client
	adyenAPI *adyen.API
	csvPath  string
}

// LinkRecord declare one split configuration record.
type LinkRecord struct {
	AccountHolderCode string `csv:"Account Holder Code"`
	ToastGUID         string `csv:"Toast GUID"`
	StoreID           string `csv:"Store ID"`
	SplitID           string `csv:"Split ID"`
}

// newLogger initializes logger for console.
func newLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	config.Encoding = "console"
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return config.Build()
}

// New creates new instance of Processor.
func New(csvPath, apiURL, apiKey string, dryRun bool) *Processor {
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}
	client := http.DefaultClient

	return &Processor{
		csvPath:  csvPath,
		logger:   logger,
		client:   client,
		adyenAPI: adyen.New(logger, client, apiURL, apiKey, dryRun),
	}
}

// Run runs parsing & split config updating.
func (p *Processor) Run(ctx context.Context) error {
	file, err := os.OpenFile(p.csvPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open CSV: %w", err)
	}
	defer file.Close()

	records := []*LinkRecord{}
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
		return fmt.Errorf("failed to process link record: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process restaurants")
	return nil
}

func (p *Processor) process(ctx context.Context, record *LinkRecord) error {
	accountHolder, err := p.adyenAPI.AccountHolder(ctx, record.AccountHolderCode)
	if err != nil {
		return fmt.Errorf("failed to get account holder: %w", err)
	}
	if err := p.updateSplitConfiguration(accountHolder, record.StoreID, record.SplitID); err != nil {
		return fmt.Errorf("failed to replace split configuration: %w", err)
	}
	if err := p.adyenAPI.UpdateAccountHolder(ctx, &accountHolder.UpdateAccountHolderRequest); err != nil {
		return fmt.Errorf("failed to update account holder: %w", err)
	}
	return nil
}

func (p *Processor) updateSplitConfiguration(accountHolder *adyen.GetAccountHolderResponse, storeID, splitID string) error {
	if len(accountHolder.AccountHolderDetails.StoreDetails) != len(accountHolder.Accounts) {
		return ErrInvalidResponse
	}
	if len(accountHolder.AccountHolderDetails.StoreDetails) != 1 {
		return ErrInvalidResponse
	}
	if accountHolder.AccountHolderDetails.StoreDetails[0].StoreID != storeID {
		return ErrStoreNotFound
	}

	accountHolder.AccountHolderDetails.StoreDetails[0].VirtualAccount = accountHolder.Accounts[0].AccountCode
	accountHolder.AccountHolderDetails.StoreDetails[0].SplitConfigurationUUID = splitID
	return nil
}
