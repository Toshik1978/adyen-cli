package cellular

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
	disable     bool
	dryRun      bool
}

// New creates new instance of Processor.
func New(
	logger *zap.Logger, client *http.Client, config *commands.Config,
	csvFilePath string, disable, production, dryRun bool) *Processor {
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
		logger:      logger,
		client:      client,
		adyenAPI:    adyen.New(logger, client, calURL, calKey, mgmtURL, mgmtKey),
		csvFilePath: csvFilePath,
		disable:     disable,
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
			Error("Failed to process terminals")
		return fmt.Errorf("failed to process cellular: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process terminals")
	return nil
}

func (p *Processor) process(ctx context.Context, record *Record) error {
	if p.dryRun {
		return nil
	}

	if err := p.adyenAPI.SetSimCardStatus(ctx, record.TerminalID, p.disable); err != nil {
		return fmt.Errorf("failed to process cellular: %w", err)
	}
	return nil
}
