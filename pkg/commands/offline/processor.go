package offline

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

// Run runs parsing & offline payments processing.
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
		return fmt.Errorf("failed to process offline payments: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process terminals")
	return nil
}

func (p *Processor) process(ctx context.Context, record *Record) error {
	terminalID := record.TerminalID
	if terminalID == "" && record.Serial != "" {
		// Get terminal ID by serial number
		terminals, err := p.adyenAPI.SearchTerminals(ctx, "", record.Serial)
		if err != nil {
			return fmt.Errorf("failed to process terminals: %w", err)
		}
		if terminals.ItemsTotal != 1 {
			return fmt.Errorf("expected 1 terminal, got %d", terminals.ItemsTotal)
		}
		terminalID = terminals.Data[0].ID
	}
	if terminalID == "" {
		return fmt.Errorf("no terminal id and serial number defined")
	}

	// Get existing terminal settings
	settings, err := p.adyenAPI.TerminalSettings(ctx, terminalID)
	if err != nil {
		return fmt.Errorf("failed to process settings: %w", err)
	}

	// Update offline payments structure and push it to Adyen
	// We should set all limits to 0
	update := adyen.SetOfflinePaymentsRequest{}
	update.OfflineProcessing = settings.OfflineProcessing
	update.StoreAndForward = settings.StoreAndForward
	p.setZero(&update)

	if p.dryRun {
		return nil
	}
	if err := p.adyenAPI.DisableOfflinePayments(ctx, terminalID, update); err != nil {
		return fmt.Errorf("failed to process offline payments: %w", err)
	}
	return nil
}

func (p *Processor) setZero(update *adyen.SetOfflinePaymentsRequest) {
	update.OfflineProcessing.ChipFloorLimit = 0
	for i := range update.OfflineProcessing.OfflineSwipeLimits {
		update.OfflineProcessing.OfflineSwipeLimits[i].Amount = 0
	}

	update.StoreAndForward.MaxPayments = 0
	for i := range update.StoreAndForward.MaxAmount {
		update.StoreAndForward.MaxAmount[i].Amount = 0
	}
}
