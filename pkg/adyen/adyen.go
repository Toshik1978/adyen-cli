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
	logger *zap.Logger
	client *http.Client
	apiURL string
	apiKey string
	dryRun bool
}

// New instantiate the new API object instance.
func New(logger *zap.Logger, client *http.Client, apiURL, apiKey string, dryRun bool) *API {
	return &API{
		apiURL: apiURL,
		apiKey: apiKey,
		dryRun: dryRun,
		logger: logger,
		client: client,
	}
}

// AccountHolder retrieve account holder by the code.
func (a *API) AccountHolder(ctx context.Context, accountHolderCode string) (*GetAccountHolderResponse, error) {
	a.logger.
		With(zap.String("AccountHolderCode", accountHolderCode)).
		Info(">> Get Account Holder")

	response, err := a.call(
		ctx,
		fmt.Sprintf("https://%s/cal/services/Account/v6/getAccountHolder", a.apiURL),
		GetAccountHolderRequest{AccountHolderCode: accountHolderCode})
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
func (a *API) UpdateAccountHolder(ctx context.Context, accountHolder *UpdateAccountHolderRequest) error {
	a.logger.
		With(zap.Any("AccountHolder", accountHolder)).
		With(zap.Bool("Dry Run", a.dryRun)).
		Info(">> Update Account Holder")

	var updated GetAccountHolderResponse
	if !a.dryRun {
		response, err := a.call(
			ctx,
			fmt.Sprintf("https://%s/cal/services/Account/v6/updateAccountHolder", a.apiURL),
			accountHolder)
		if err != nil {
			return fmt.Errorf("failed to get update holder: %w", err)
		}

		if err := json.Unmarshal(response, &updated); err != nil {
			return fmt.Errorf("failed to unmarshal Adyen response: %w", err)
		}
	}

	a.logger.
		With(zap.Any("AccountHolder", updated)).
		With(zap.Bool("Dry Run", a.dryRun)).
		Info("<< Update Account Holder")
	return nil
}

func (a *API) call(ctx context.Context, url string, data interface{}) ([]byte, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-API-key", a.apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := a.client.Do(request) // nolint:bodyclose
	defer closeResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to call Adyen: %w", err)
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
