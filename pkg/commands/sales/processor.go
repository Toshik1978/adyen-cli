package sales

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
	csvFilePath string, production, dryRun bool) *Processor {
	var calURL, calKey, mgmtURL, mgmtKey, kycURL, kycKey, balURL, balKey string
	switch {
	case production:
		calURL = config.AdyenCalURL
		calKey = config.AdyenCalKey
		mgmtURL = config.AdyenMgmtURL
		mgmtKey = config.AdyenMgmtKey
		kycURL = config.AdyenKycURL
		kycKey = config.AdyenKycKey
		balURL = config.AdyenBalURL
		balKey = config.AdyenBalKey
	case !production:
		calURL = config.AdyenCalTestURL
		calKey = config.AdyenCalTestKey
		mgmtURL = config.AdyenMgmtTestURL
		mgmtKey = config.AdyenMgmtTestKey
		kycURL = config.AdyenKycTestURL
		kycKey = config.AdyenKycTestKey
		balURL = config.AdyenBalTestURL
		balKey = config.AdyenBalTestKey
	}

	gocsv.SetHeaderNormalizer(strings.ToUpper)

	return &Processor{
		logger:      logger,
		client:      client,
		adyenAPI:    adyen.New(logger, client, calURL, calKey, mgmtURL, mgmtKey, kycURL, kycKey, balURL, balKey),
		csvFilePath: csvFilePath,
		dryRun:      dryRun,
	}
}

// Run runs sales close time updater.
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
	errs := make([]error, 0, len(records))
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
		return fmt.Errorf("failed to process sales close time record: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process restaurants")
	return nil
}

func (p *Processor) process(ctx context.Context, record *Record) error {
	balanceID := record.BalanceID
	if record.AccountHolderID != "" && record.BalanceID == "" {
		acc, err := p.adyenAPI.BalanceAccountHolder(ctx, record.AccountHolderID)
		if err != nil {
			return fmt.Errorf("failed to get balance account by holder (%s): %w", record.AccountHolderID, err)
		}
		balanceID = acc.PrimaryBalanceAccount
	}
	if balanceID == "" {
		return fmt.Errorf("no balance account identified: %s", record.AccountHolderID)
	}

	if !p.dryRun {
		if _, err := p.adyenAPI.SetSalesCloseTime(ctx, balanceID, record.CloseTime, record.Delays); err != nil {
			return fmt.Errorf("failed to change sales close time: %w", err)
		}
	}
	return nil
}
