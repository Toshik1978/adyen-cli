package reassign

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

// Run runs parsing & terminal re-assignment.
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
			Error("Failed to process terminals")
		return fmt.Errorf("failed to process re-assignment: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process terminals")
	return nil
}

func (p *Processor) process(ctx context.Context, record *Record) error {
	terminalID, err := p.searchTerminal(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to search terminal: %w", err)
	}

	storeID, err := p.searchStore(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to search store: %w", err)
	}

	if p.dryRun {
		return nil
	}
	if err := p.adyenAPI.ReassignTerminal(ctx, terminalID, record.MerchantID, storeID); err != nil {
		return fmt.Errorf("failed to process terminal re-assignment: %w", err)
	}
	return nil
}

func (p *Processor) searchTerminal(ctx context.Context, record *Record) (string, error) {
	terminalID := record.TerminalID
	if terminalID == "" && record.Serial != "" {
		// Get terminal ID by serial number
		terminals, err := p.adyenAPI.SearchTerminals(ctx, "", record.Serial)
		if err != nil {
			return "", fmt.Errorf("failed to process terminals: %w", err)
		}
		if terminals.ItemsTotal != 1 {
			return "", fmt.Errorf("expected 1 terminal, got %d", terminals.ItemsTotal)
		}
		terminalID = terminals.Data[0].ID
	}
	if terminalID == "" {
		return "", fmt.Errorf("no terminal id and serial number defined")
	}
	return terminalID, nil
}

func (p *Processor) searchStore(ctx context.Context, record *Record) (string, error) {
	storeID := record.StoreID
	if storeID != "" {
		// Need to convert Adyen Store GUID to the management ID.
		stores, err := p.adyenAPI.SearchStores(ctx, record.StoreID)
		if err != nil {
			return "", fmt.Errorf("failed to get all stores: %w", err)
		}
		if stores.ItemsTotal != 1 || len(stores.Data) != 1 {
			return "", ErrInvalidResponse
		}
		if stores.Data[0].Reference != record.StoreID {
			return "", fmt.Errorf("store ID not found: %s %s", stores.Data[0].Reference, record.StoreID)
		}
		storeID = stores.Data[0].ID
	}
	if storeID == "" {
		return "", fmt.Errorf("no store id defined")
	}
	return storeID, nil
}
