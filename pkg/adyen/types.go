package adyen

// Adyen API types

// AccountHolderRequest declare general account holder request.
type AccountHolderRequest struct {
	AccountHolderCode string `json:"accountHolderCode"`
}

// GetAccountHolderResponse declare get account holder response.
type GetAccountHolderResponse struct {
	UpdateAccountHolderRequest

	Accounts []struct {
		AccountCode    string `json:"accountCode"`
		Description    string `json:"description"`
		PayoutSchedule struct {
			Schedule string `json:"schedule"`
		} `json:"payoutSchedule"`
		PayoutSpeed string `json:"payoutSpeed"`
		Status      string `json:"status"`
	} `json:"accounts"`
}

// UpdateAccountHolderRequest declare account holder update request.
type UpdateAccountHolderRequest struct {
	AccountHolderCode    string `json:"accountHolderCode"`
	AccountHolderDetails struct {
		StoreDetails []struct {
			SplitConfigurationUUID *string `json:"splitConfigurationUUID,omitempty"`
			VirtualAccount         *string `json:"virtualAccount,omitempty"`
			Status                 *string `json:"status,omitempty"`
			StoreID                string  `json:"store"`
		} `json:"storeDetails"`
	} `json:"accountHolderDetails"`
}

// CloseAccountHolderResponse declare close account holder response.
type CloseAccountHolderResponse struct {
	AccountHolderStatus struct {
		Status       string `json:"status"`
		PspReference string `json:"pspReference"`
		ResultCode   string `json:"resultCode"`
	}
}

// UpdateSplitConfigurationRequest declare split configuration update request on Balance.
type UpdateSplitConfigurationRequest struct {
	SplitConfiguration struct {
		BalanceAccountID     string `json:"balanceAccountId"`
		SplitConfigurationID string `json:"splitConfigurationId"`
	} `json:"splitConfiguration"`
}

// UpdateSplitConfigurationResponse declare split configuration update response on Balance.
type UpdateSplitConfigurationResponse struct {
	StoreID            string `json:"id"`
	MerchantID         string `json:"merchantId"`
	Status             string `json:"status"`
	SplitConfiguration struct {
		BalanceAccountID     string `json:"balanceAccountId"`
		SplitConfigurationID string `json:"splitConfigurationId"`
	} `json:"splitConfiguration"`
}

// SetStoreStatusRequest declare set store status request.
type SetStoreStatusRequest struct {
	Status string `json:"status"`
}

// GetStoreResponse declare get store information response.
type GetStoreResponse struct {
	ID         string `json:"id"`
	MerchantID string `json:"merchantId"`
	Reference  string `json:"reference"`
	Status     string `json:"status"`
}

// GetAllStoresResponse declare get all stores response.
type GetAllStoresResponse struct {
	Data       []GetStoreResponse `json:"data"`
	ItemsTotal int64              `json:"itemsTotal"`
}

// ReassignTerminalRequest declare reassign terminal request.
type ReassignTerminalRequest struct {
	StoreID    *string `json:"storeId,omitempty"`
	MerchantID *string `json:"merchantId,omitempty"`
	Inventory  *bool   `json:"inventory,omitempty"`
}

// Sim Card statuses.

const (
	SimCardActivated = "ACTIVATED"
	SimCardInventory = "INVENTORY"
)

// SetSimCardStatusRequest declare request to update sim-card status.
type SetSimCardStatusRequest struct {
	Connectivity struct {
		Status string `json:"simcardStatus"`
	} `json:"connectivity"`
}

// TerminalSettingsResponse declare response with terminal settings.
type TerminalSettingsResponse struct { //nolint:govet
	Timeouts struct {
		FromActiveToSleep int `json:"fromActiveToSleep"`
	} `json:"timeouts"`
	Hardware struct {
		DisplayMaximumBackLight int `json:"displayMaximumBackLight"`
		RestartHour             int `json:"restartHour"`
	} `json:"hardware"`
	Connectivity struct {
		Status string `json:"simcardStatus"`
	} `json:"connectivity"`
}
