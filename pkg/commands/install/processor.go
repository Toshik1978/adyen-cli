package install

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"go.uber.org/zap"

	"github.com/Toshik1978/csv2adyen/pkg/adyen"
	"github.com/Toshik1978/csv2adyen/pkg/commands"
)

var (
	// ErrInvalidResponse means we have a wrong Adyen response.
	ErrInvalidResponse = errors.New("store details count not equal to accounts count")
	// ErrTooManyTerminals means we have more than 100 terminals per store.
	ErrTooManyTerminals = errors.New("too many terminals assigned to the store")
	// ErrNoAppFound means we could not find the relevant app version.
	ErrNoAppFound = errors.New("no app found")
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
		dryRun:      dryRun,
	}
}

// Run runs parsing & app installation.
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
			Error("Failed to process installations")
		return fmt.Errorf("failed to process installations: %w", errors.Join(errs...))
	}

	p.logger.
		With(zap.Int("Success Count", successCnt)).
		Info("Finished to process installations")
	return nil
}

func (p *Processor) process(ctx context.Context, record *Record) error {
	var terminalIDs []string
	if record.StoreID != "" {
		// Need to convert Adyen Store GUID to the management ID.
		storeID, err := p.searchStore(ctx, record.StoreID)
		if err != nil {
			return fmt.Errorf("failed to get store ID by UUID: %w", err)
		}

		ids, err := p.searchTerminals(ctx, storeID, record.TerminalFilter)
		if err != nil {
			return fmt.Errorf("failed to get terminals: %w", err)
		}
		terminalIDs = ids
	} else {
		terminalIDs = []string{record.TerminalID}
	}
	appID, err := p.appID(ctx, record.CompanyID, record.PackageName, record.VersionName)
	if err != nil {
		return fmt.Errorf("failed to get app id: %w", err)
	}

	if p.dryRun {
		return nil
	}

	scheduledAt := record.Date
	if scheduledAt == "" {
		scheduledAt = time.Now().Add(2 * time.Minute).Format(time.RFC3339)
		if scheduledAt[len(scheduledAt)-1] == 'Z' {
			scheduledAt = strings.Replace(scheduledAt, "Z", "+0000", 1)
		}
		if scheduledAt[len(scheduledAt)-3] == ':' {
			scheduledAt = scheduledAt[:len(scheduledAt)-3] + scheduledAt[len(scheduledAt)-2:]
		}
		scheduledAt = strings.Replace(scheduledAt, "Z", "", 1)
	}

	if err := p.adyenAPI.InstallAndroidApp(ctx, appID, "", terminalIDs, scheduledAt); err != nil {
		return fmt.Errorf("failed to process app installations: %w", err)
	}
	return nil
}

func (p *Processor) searchStore(ctx context.Context, storeID string) (string, error) {
	stores, err := p.adyenAPI.SearchStores(ctx, storeID)
	if err != nil {
		return "", fmt.Errorf("failed to get all stores: %w", err)
	}
	if stores.ItemsTotal != 1 || len(stores.Data) != 1 {
		return "", ErrInvalidResponse
	}
	if stores.Data[0].Reference != storeID {
		return "", fmt.Errorf("store ID not found: %s %s", stores.Data[0].Reference, storeID)
	}
	return stores.Data[0].ID, nil
}

func (p *Processor) searchTerminals(ctx context.Context, storeID, searchQuery string) ([]string, error) {
	terminals, err := p.adyenAPI.SearchTerminals(ctx, storeID, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get terminals: %w", err)
	}
	if terminals.PagesTotal > 1 {
		return nil, ErrTooManyTerminals
	}

	terminalIDs := make([]string, 0, len(terminals.Data))
	for i := range terminals.Data {
		terminalIDs = append(terminalIDs, terminals.Data[i].ID)
	}
	return terminalIDs, nil
}

func (p *Processor) appID(ctx context.Context, companyID, packageName, versionName string) (string, error) {
	apps, err := p.adyenAPI.SearchAndroidApps(ctx, companyID, packageName)
	if err != nil {
		return "", fmt.Errorf("failed to get all apps: %w", err)
	}

	for i := range apps.Data {
		if apps.Data[i].VersionName == versionName {
			return apps.Data[i].ID, nil
		}
	}
	return "", ErrNoAppFound
}
