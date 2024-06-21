package sweep

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

// Run runs fix sweep configuration.
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
		return fmt.Errorf("failed to process fix sweep configuration record: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process restaurants")
	return nil
}

func (p *Processor) process(ctx context.Context, record *Record) error {
	balanceID, legalEntityID, err := p.ids(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to get ids: %w", err)
	}

	legalEntity, err := p.adyenAPI.LegalEntity(ctx, legalEntityID)
	if err != nil {
		return fmt.Errorf("failed to get legal entity: %w", err)
	}
	if len(legalEntity.TransferInstruments) != 1 {
		return fmt.Errorf("expected 1 legal instrument, got %d (%s)", len(legalEntity.TransferInstruments), legalEntityID)
	}

	sweeps, err := p.adyenAPI.Sweeps(ctx, balanceID)
	if err != nil {
		return fmt.Errorf("failed to get sweeps: %w", err)
	}
	if len(sweeps.Sweeps) != 1 {
		return fmt.Errorf("expected 1 sweep configuration, got %d (%s)", len(sweeps.Sweeps), balanceID)
	}

	if sweeps.Sweeps[0].Counterparty.TransferInstrumentID == legalEntity.TransferInstruments[0].ID {
		p.logger.
			With(zap.String("BalanceID", balanceID)).
			With(zap.String("LegalEntityID", legalEntityID)).
			Info("Transfer instrument already valid")
		return nil
	}

	if !p.dryRun {
		_, err := p.adyenAPI.UpdateSweep(ctx, balanceID, sweeps.Sweeps[0].ID, legalEntity.TransferInstruments[0].ID)
		if err != nil {
			return fmt.Errorf("failed to update sweep (%s): %w", balanceID, err)
		}
	}
	return nil
}

func (p *Processor) ids(ctx context.Context, record *Record) (string, string, error) {
	accountHolderID := record.AccountHolderID
	if accountHolderID == "" && record.BalanceID != "" {
		acc, err := p.adyenAPI.BalanceAccount(ctx, record.BalanceID)
		if err != nil {
			return "", "", fmt.Errorf("failed to get balance account (%s): %w", record.BalanceID, err)
		}
		accountHolderID = acc.AccountHolderID
	}
	if accountHolderID == "" {
		return "", "", fmt.Errorf("no balance account holder identified: %s", record.BalanceID)
	}

	acc, err := p.adyenAPI.BalanceAccountHolder(ctx, accountHolderID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get balance account by holder (%s): %w", record.AccountHolderID, err)
	}

	balanceID := acc.PrimaryBalanceAccount
	if balanceID == "" {
		return "", "", fmt.Errorf("no balance account identified: %s", accountHolderID)
	}

	legalEntityID := acc.LegalEntityID
	if legalEntityID == "" {
		return "", "", fmt.Errorf("no legal entity identified: %s", accountHolderID)
	}
	return balanceID, legalEntityID, nil
}
