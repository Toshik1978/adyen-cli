package adyen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
}

// New instantiate the new API object instance.
func New(logger *zap.Logger, client *http.Client, calURL, calKey, mgmtURL, mgmtKey string) *API {
	return &API{
		calURL:  calURL,
		calKey:  calKey,
		mgmtURL: mgmtURL,
		mgmtKey: mgmtKey,
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
		return nil, fmt.Errorf("failed to get account holder: %w", err)
	}

	var accountHolder GetAccountHolderResponse
	if err := json.Unmarshal(response, &accountHolder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.Any("AccountHolder", accountHolder)).
		Info("<< Get Account Holder")
	return &accountHolder, nil
}

// UpdateAccountHolder updates account holder.
func (a *API) UpdateAccountHolder(ctx context.Context, accountHolder *UpdateAccountHolderRequest) error { //nolint:dupl
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
		With(zap.Any("AccountHolder", updated)).
		Info("<< Update Account Holder")
	return nil
}

// CloseAccountHolder closes account holder.
func (a *API) CloseAccountHolder(ctx context.Context, accountHolderCode string) error { //nolint:dupl
	a.logger.
		With(zap.Any("AccountHolder", accountHolderCode)).
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
		With(zap.Any("AccountHolder", closed)).
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

// GetAllStores gets store's management IDs by store ID.
func (a *API) GetAllStores(ctx context.Context, storeID string) (*GetAllStoresResponse, error) {
	a.logger.
		With(zap.Any("StoreID", storeID)).
		Info(">> Get All Store")

	response, err := a.call(
		ctx,
		http.MethodGet,
		fmt.Sprintf("https://%s/v1/stores?reference=%s", a.mgmtURL, storeID),
		a.mgmtKey,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get all stores: %w", err)
	}

	var stores GetAllStoresResponse
	if err := json.Unmarshal(response, &stores); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Adyen response: %w", err)
	}

	a.logger.
		With(zap.Any("StoreID", storeID)).
		With(zap.Any("Response", stores)).
		Info("<< Get All Store")
	return &stores, nil
}

// SetStoreStatus set store status by management ID.
func (a *API) SetStoreStatus(ctx context.Context, storeMgmtID, status string) error {
	a.logger.
		With(zap.Any("StoreID", storeMgmtID)).
		With(zap.Any("Status", status)).
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
		With(zap.Any("StoreID", storeMgmtID)).
		With(zap.Any("Status", status)).
		With(zap.Any("Response", store)).
		Info("<< Set Store Status")
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
		body, _ := io.ReadAll(response.Body)
		a.logger.
			With(zap.ByteString("Response", body)).
			Error("Failed call")
		return nil, fmt.Errorf("failed to call Adyen, HTTP status: %d", response.StatusCode)
	}
	if response == nil {
		return nil, fmt.Errorf("failed to call Adyen, surprising nil response")
	}

	return io.ReadAll(response.Body)
}

func closeResponse(response *http.Response) {
	if response != nil {
		if response.Body != nil {
			_ = response.Body.Close()
		}
	}
}
