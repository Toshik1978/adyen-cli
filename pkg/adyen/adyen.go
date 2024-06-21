package adyen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AlekSi/pointer"
	"go.uber.org/zap"
)

// API define Adyen API.
type API struct {
	logger  *zap.Logger
	client  *http.Client
	calURL  string
	calKey  string
	mgmtURL string
	mgmtKey string
	kycURL  string
	kycKey  string
	balURL  string
	balKey  string
}

// New instantiate the new API object instance.
func New(
	logger *zap.Logger, client *http.Client, calURL, calKey, mgmtURL, mgmtKey, kycURL, kycKey, balURL, balKey string,
) *API {
	return &API{
		calURL:  calURL,
		calKey:  calKey,
		mgmtURL: mgmtURL,
		mgmtKey: mgmtKey,
		kycURL:  kycURL,
		kycKey:  kycKey,
		balURL:  balURL,
		balKey:  balKey,
		logger:  logger,
		client:  client,
	}
}

// AccountHolder retrieve account holder by the code.
func (a *API) AccountHolder(ctx context.Context, accountHolderCode string) (*GetAccountHolderResponse, error) {
	a.logger.
		With(zap.String("AccountHolderCode", accountHolderCode)).
		Info(">> Get Account Holder")

	response, err := a.call(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://%s/cal/services/Account/v6/getAccountHolder", a.calURL),
		a.calKey,
		&AccountHolderRequest{AccountHolderCode: accountHolderCode})
	if err != nil {
		return nil, fmt.Errorf("failed to get account holder (%s): %w", accountHolderCode, err)
	}

	var accountHolder GetAccountHolderResponse
	if err := json.Unmarshal(response, &accountHolder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("AccountHolderCode", accountHolderCode)).
		With(zap.Any("Response", accountHolder)).
		Info("<< Get Account Holder")
	return &accountHolder, nil
}

// UpdateAccountHolder updates account holder.
func (a *API) UpdateAccountHolder(ctx context.Context, accountHolder *UpdateAccountHolderRequest) error {
	a.logger.
		With(zap.Any("AccountHolder", accountHolder)).
		Info(">> Update Account Holder")

	response, err := a.call(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://%s/cal/services/Account/v6/updateAccountHolder", a.calURL),
		a.calKey,
		accountHolder)
	if err != nil {
		return fmt.Errorf("failed to update account holder: %w", err)
	}

	var updated GetAccountHolderResponse
	if err := json.Unmarshal(response, &updated); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.Any("AccountHolder", accountHolder)).
		With(zap.Any("Response", updated)).
		Info("<< Update Account Holder")
	return nil
}

// CloseAccountHolder closes account holder.
func (a *API) CloseAccountHolder(ctx context.Context, accountHolderCode string) error {
	a.logger.
		With(zap.String("AccountHolderCode", accountHolderCode)).
		Info(">> Close Account Holder")

	response, err := a.call(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://%s/cal/services/Account/v6/closeAccountHolder", a.calURL),
		a.calKey,
		&AccountHolderRequest{AccountHolderCode: accountHolderCode})
	if err != nil {
		return fmt.Errorf("failed to close account holder: %w", err)
	}

	var closed CloseAccountHolderResponse
	if err := json.Unmarshal(response, &closed); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("AccountHolderCode", accountHolderCode)).
		With(zap.Any("Response", closed)).
		Info("<< Close Account Holder")
	return nil
}

// UpdateSplitConfiguration updates split configuration on Balance.
func (a *API) UpdateSplitConfiguration(
	ctx context.Context, merchantID, storeID string, config *UpdateSplitConfigurationRequest) error {
	a.logger.
		With(zap.String("MerchantID", merchantID)).
		With(zap.Any("SplitConfiguration", config)).
		Info(">> Update Split Configuration")

	response, err := a.call(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("https://%s/v1/merchants/%s/stores/%s", a.mgmtURL, merchantID, storeID),
		a.mgmtKey,
		config)
	if err != nil {
		return fmt.Errorf("failed to update split configuration: %w", err)
	}

	var updated UpdateSplitConfigurationResponse
	if err := json.Unmarshal(response, &updated); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("MerchantID", merchantID)).
		With(zap.Any("SplitConfiguration", config)).
		With(zap.Any("Response", updated)).
		Info("<< Update Split Configuration")
	return nil
}

// SearchStores gets store's management IDs by store ID.
func (a *API) SearchStores(ctx context.Context, storeID string) (*SearchStoresResponse, error) {
	a.logger.
		With(zap.String("StoreID", storeID)).
		Info(">> Get All Store")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/v3/stores?reference=%s", a.mgmtURL, storeID),
		a.mgmtKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stores: %w", err)
	}

	var stores SearchStoresResponse
	if err := json.Unmarshal(response, &stores); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("StoreID", storeID)).
		With(zap.Any("Response", stores)).
		Info("<< Get All Store")
	return &stores, nil
}

// SetStoreStatus set store status by management ID.
func (a *API) SetStoreStatus(ctx context.Context, storeMgmtID, status string) error {
	a.logger.
		With(zap.String("StoreID", storeMgmtID)).
		With(zap.String("Status", status)).
		Info(">> Set Store Status")

	response, err := a.call(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("https://%s/v1/stores/%s", a.mgmtURL, storeMgmtID),
		a.mgmtKey,
		&SetStoreStatusRequest{Status: status})
	if err != nil {
		return fmt.Errorf("failed to set store status: %w", err)
	}

	var store GetStoreResponse
	if err := json.Unmarshal(response, &store); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("StoreID", storeMgmtID)).
		With(zap.String("Status", status)).
		With(zap.Any("Response", store)).
		Info("<< Set Store Status")
	return nil
}

// AddPaymentMethod adds payment method to the store.
func (a *API) AddPaymentMethod(
	ctx context.Context, merchantID, storeID, businessLineID, method, currency string,
) (*AddPaymentMethodResponse, error) {
	a.logger.
		With(zap.String("MerchantID", merchantID)).
		With(zap.String("StoreID", storeID)).
		With(zap.String("BusinessLineID", businessLineID)).
		With(zap.String("Method", method)).
		With(zap.String("Currency", currency)).
		Info(">> Add Payment Method")

	req := AddPaymentMethodRequest{
		Type:           method,
		BusinessLineID: businessLineID,
		StoreIDs:       []string{storeID},
		Currencies:     []string{currency},
	}
	response, err := a.call(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://%s/v3/merchants/%s/paymentMethodSettings", a.mgmtURL, merchantID),
		a.mgmtKey,
		&req)
	if err != nil {
		return nil, fmt.Errorf("failed to add payment method: %w", err)
	}

	var resp AddPaymentMethodResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("MerchantID", merchantID)).
		With(zap.String("StoreID", storeID)).
		With(zap.String("BusinessLineID", businessLineID)).
		With(zap.String("Method", method)).
		With(zap.String("Currency", currency)).
		With(zap.Any("Response", resp)).
		Info("<< Add Payment Method")
	return &resp, nil
}

// BalanceAccount retrieve information about balance account.
func (a *API) BalanceAccount(
	ctx context.Context, balanceID string,
) (*GetBalanceAccountResponse, error) {
	a.logger.
		With(zap.String("BalanceID", balanceID)).
		Info(">> Balance Account")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/bcl/v2/balanceAccounts/%s", a.balURL, balanceID),
		a.balKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance account: %w", err)
	}

	var resp GetBalanceAccountResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("BalanceID", balanceID)).
		With(zap.Any("Response", resp)).
		Info("<< Balance Account")
	return &resp, nil
}

// BalanceAccountHolder retrieve information about balance account holder.
func (a *API) BalanceAccountHolder(
	ctx context.Context, accountHolderID string,
) (*GetBalanceAccountHolderResponse, error) {
	a.logger.
		With(zap.String("AccountHolderID", accountHolderID)).
		Info(">> Balance Account Holder")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/bcl/v2/accountHolders/%s", a.balURL, accountHolderID),
		a.balKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance account holder: %w", err)
	}

	var resp GetBalanceAccountHolderResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("AccountHolderID", accountHolderID)).
		With(zap.Any("Response", resp)).
		Info("<< Balance Account Holder")
	return &resp, nil
}

// LegalEntity retrieve information about legal entity.
func (a *API) LegalEntity(ctx context.Context, legalEntityID string) (*GetLegalEntityResponse, error) {
	a.logger.
		With(zap.String("LegalEntityID", legalEntityID)).
		Info(">> Legal Entity")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/lem/v3/legalEntities/%s", a.kycURL, legalEntityID),
		a.kycKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get legal entity: %w", err)
	}

	var resp GetLegalEntityResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("LegalEntityID", legalEntityID)).
		With(zap.Any("Response", resp)).
		Info("<< Legal Entity")
	return &resp, nil
}

// Sweeps retrieve all sweep configurations.
func (a *API) Sweeps(ctx context.Context, balanceID string) (*GetSweepsResponse, error) {
	a.logger.
		With(zap.String("BalanceID", balanceID)).
		Info(">> Sweeps")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/bcl/v2/balanceAccounts/%s/sweeps", a.balURL, balanceID),
		a.balKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get sweeps: %w", err)
	}

	var resp GetSweepsResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("BalanceID", balanceID)).
		With(zap.Any("Response", resp)).
		Info("<< Sweeps")
	return &resp, nil
}

// UpdateSweep updates sweep configuration.
func (a *API) UpdateSweep(ctx context.Context, balanceID, sweepID, transferInstrumentID string) (*Sweep, error) {
	a.logger.
		With(zap.String("BalanceID", balanceID)).
		With(zap.String("SweepID", sweepID)).
		With(zap.String("TransferInstrumentID", transferInstrumentID)).
		Info(">> Update Sweep")

	req := UpdateSweepRequest{}
	req.Counterparty.TransferInstrumentID = transferInstrumentID
	req.Status = SweepActive

	response, err := a.call(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("https://%s/bcl/v2/balanceAccounts/%s/sweeps/%s", a.balURL, balanceID, sweepID),
		a.balKey,
		&req)
	if err != nil {
		return nil, fmt.Errorf("failed to update sweep configuration: %w", err)
	}

	var resp Sweep
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("BalanceID", balanceID)).
		With(zap.String("SweepID", sweepID)).
		With(zap.String("TransferInstrumentID", transferInstrumentID)).
		With(zap.Any("Response", resp)).
		Info("<< Update Sweep")
	return &resp, nil
}

// SetSalesCloseTime changes sales close time for the store.
func (a *API) SetSalesCloseTime(
	ctx context.Context, balanceID, closingTime string, delays int,
) (*GetBalanceAccountResponse, error) {
	a.logger.
		With(zap.String("BalanceID", balanceID)).
		With(zap.String("ClosingTime", closingTime)).
		With(zap.Int("Delays", delays)).
		Info(">> Change Sales Close Time")

	req := SetSalesCloseTimeRequest{}
	req.PlatformPaymentConfiguration.SalesDayClosingTime = closingTime
	req.PlatformPaymentConfiguration.SettlementDelayDays = delays

	response, err := a.call(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("https://%s/bcl/v2/balanceAccounts/%s", a.balURL, balanceID),
		a.balKey,
		&req)
	if err != nil {
		return nil, fmt.Errorf("failed to change sales close time: %w", err)
	}

	var resp GetBalanceAccountResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("BalanceID", balanceID)).
		With(zap.String("ClosingTime", closingTime)).
		With(zap.Int("Delays", delays)).
		With(zap.Any("Response", resp)).
		Info("<< Change Sales Close Time")
	return &resp, nil
}

// ReassignTerminal reassign terminal to store or merchant (always inventory)
func (a *API) ReassignTerminal(ctx context.Context, terminalID, merchantID, storeID string) error {
	a.logger.
		With(zap.String("TerminalID", terminalID)).
		With(zap.String("MerchantID", merchantID)).
		With(zap.String("StoreID", storeID)).
		Info(">> Re-assign Terminal")

	req := ReassignTerminalRequest{}
	switch {
	case storeID != "":
		req.StoreID = pointer.ToString(storeID)
	case merchantID != "":
		req.MerchantID = pointer.ToString(merchantID)
		req.Inventory = pointer.ToBool(true)
	default:
		return fmt.Errorf("no merchant id and store id for terminal re-assignment: %s", terminalID)
	}

	_, err := a.call(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://%s/v3/terminals/%s/reassign", a.mgmtURL, terminalID),
		a.mgmtKey,
		&req)
	if err != nil {
		return fmt.Errorf("failed to re-assign terminal: %w", err)
	}

	a.logger.
		With(zap.String("TerminalID", terminalID)).
		With(zap.String("MerchantID", merchantID)).
		With(zap.String("StoreID", storeID)).
		Info("<< Re-assign Terminal")
	return nil
}

// TerminalSettings gets terminal settings.
func (a *API) TerminalSettings(ctx context.Context, terminalID string) (*TerminalSettingsResponse, error) {
	a.logger.
		With(zap.String("TerminalID", terminalID)).
		Info(">> Get Terminal Settings")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/v3/terminals/%s/terminalSettings", a.mgmtURL, terminalID),
		a.mgmtKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal settings: %w", err)
	}

	var settings TerminalSettingsResponse
	if err := json.Unmarshal(response, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("TerminalID", terminalID)).
		With(zap.Any("Response", settings)).
		Info("<< Get Terminal Settings")
	return &settings, nil
}

// SetSimCardStatus set sim card status.
func (a *API) SetSimCardStatus(ctx context.Context, terminalID string, disable bool) error {
	a.logger.
		With(zap.String("TerminalID", terminalID)).
		With(zap.Bool("Disable", disable)).
		Info(">> Set Sim Card Status")

	req := SetSimCardStatusRequest{}
	if disable {
		req.Connectivity.Status = SimCardInventory
	} else {
		req.Connectivity.Status = SimCardActivated
	}

	response, err := a.call(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("https://%s/v3/terminals/%s/terminalSettings", a.mgmtURL, terminalID),
		a.mgmtKey,
		&req)
	if err != nil {
		return fmt.Errorf("failed to set simcard status: %w", err)
	}

	var updated TerminalSettingsResponse
	if err := json.Unmarshal(response, &updated); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("TerminalID", terminalID)).
		With(zap.Bool("Disable", disable)).
		With(zap.Any("Response", updated)).
		Info("<< Set Sim Card Status")
	return nil
}

// DisableOfflinePayments disables offline payments per device.
func (a *API) DisableOfflinePayments(ctx context.Context, terminalID string, settings SetOfflinePaymentsRequest) error {
	a.logger.
		With(zap.String("TerminalID", terminalID)).
		With(zap.Any("Settings", settings)).
		Info(">> Disable Offline Payments")

	response, err := a.call(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("https://%s/v3/terminals/%s/terminalSettings", a.mgmtURL, terminalID),
		a.mgmtKey,
		&settings)
	if err != nil {
		return fmt.Errorf("failed to disable offline payments: %w", err)
	}

	var updated TerminalSettingsResponse
	if err := json.Unmarshal(response, &updated); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("TerminalID", terminalID)).
		With(zap.Any("Settings", settings)).
		With(zap.Any("Response", updated)).
		Info("<< Disable Offline Payments")
	return nil
}

// SearchTerminals gets the list of terminals.
func (a *API) SearchTerminals(ctx context.Context, storeID, searchQuery string) (*SearchTerminalsResponse, error) {
	a.logger.
		With(zap.String("StoreID", storeID)).
		With(zap.String("SearchQuery", searchQuery)).
		Info(">> Get Store Terminals")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/v3/terminals?storeIds=%s&searchQuery=%s&pageSize=100", a.mgmtURL, storeID, searchQuery),
		a.mgmtKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all terminals: %w", err)
	}

	var terminals SearchTerminalsResponse
	if err := json.Unmarshal(response, &terminals); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("StoreID", storeID)).
		With(zap.String("SearchQuery", searchQuery)).
		With(zap.Any("Response", terminals)).
		Info("<< Get Store Terminals")
	return &terminals, nil
}

// SearchAndroidApps gets the list of android apps.
func (a *API) SearchAndroidApps(ctx context.Context, companyID, packageName string) (*SearchAndroidAppsResponse, error) {
	a.logger.
		With(zap.String("CompanyID", companyID)).
		With(zap.String("PackageName", packageName)).
		Info(">> Get Android Apps")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/v3/companies/%s/androidApps?packageName=%s", a.mgmtURL, companyID, packageName),
		a.mgmtKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all apps: %w", err)
	}

	var apps SearchAndroidAppsResponse
	if err := json.Unmarshal(response, &apps); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("CompanyID", companyID)).
		With(zap.String("PackageName", packageName)).
		With(zap.Any("Response", apps)).
		Info("<< Get Android Apps")
	return &apps, nil
}

// InstallAndroidApp schedule action to install android app.
func (a *API) InstallAndroidApp(ctx context.Context, appID, storeID string, terminalIDs []string, at string) error {
	a.logger.
		With(zap.String("AppID", appID)).
		With(zap.String("StoreID", storeID)).
		With(zap.Strings("TerminalIDs", terminalIDs)).
		With(zap.String("ScheduledAt", at)).
		Info(">> Install Android App")

	req := ScheduleActionRequest{
		TerminalIDs: terminalIDs,
		StoreID:     storeID,
		ScheduledAt: at,
	}
	req.ActionDetails.Type = "InstallAndroidApp"
	req.ActionDetails.AppID = appID

	response, err := a.call(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://%s/v3/terminals/scheduleActions", a.mgmtURL),
		a.mgmtKey,
		&req)
	if err != nil {
		return fmt.Errorf("failed to schedule app install: %w", err)
	}

	var scheduled ScheduleActionResponse
	if err := json.Unmarshal(response, &scheduled); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.String("AppID", appID)).
		With(zap.String("StoreID", storeID)).
		With(zap.Strings("TerminalIDs", terminalIDs)).
		With(zap.String("ScheduledAt", at)).
		With(zap.Any("Response", response)).
		Info("<< Install Android App")
	return nil
}

func (a *API) call(ctx context.Context, method, url, key string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data")
		}
		body = bytes.NewReader(buf)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-API-key", key)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := a.client.Do(request) // nolint:bodyclose
	defer closeResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to call Adyen: %w", err)
	}
	if response != nil && response.StatusCode != http.StatusOK {
		return nil, a.createError(response.Body)
	}
	if response == nil {
		return nil, fmt.Errorf("failed to call Adyen, surprising nil response")
	}

	return io.ReadAll(response.Body)
}

func (a *API) createError(response io.Reader) error {
	body, _ := io.ReadAll(response)
	a.logger.
		With(zap.ByteString("Response", body)).
		Error("Failed call")

	var adyenErr Error
	if err := json.Unmarshal(body, &adyenErr); err != nil {
		return fmt.Errorf("failed to unmarshal Adyen error: %w", err)
	}
	return fmt.Errorf("failed to call Adyen: %s (%s), HTTP status: %d", adyenErr.Title, adyenErr.Detail, adyenErr.Status)
}

func closeResponse(response *http.Response) {
	if response != nil {
		if response.Body != nil {
			_ = response.Body.Close()
		}
	}
}
