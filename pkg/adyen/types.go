package adyen

import "time"

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

// GetStoreTerminalsResponse declare response for get store terminals request.
type GetStoreTerminalsResponse struct { //nolint:govet
	ItemsTotal int `json:"itemsTotal"`
	PagesTotal int `json:"pagesTotal"`
	Data       []struct {
		ID                string    `json:"id"`
		Model             string    `json:"model"`
		SerialNumber      string    `json:"serialNumber"`
		LastActivityAt    time.Time `json:"lastActivityAt"`
		LastTransactionAt time.Time `json:"lastTransactionAt"`
		FirmwareVersion   string    `json:"firmwareVersion"`
		Assignment        struct {
			CompanyID  string `json:"companyId"`
			MerchantID string `json:"merchantId"`
			StoreID    string `json:"storeId"`
			Status     string `json:"status"`
		} `json:"assignment"`
		Connectivity struct {
			Cellular struct {
				Status string `json:"status"`
				Iccid  string `json:"iccid"`
			} `json:"cellular"`
			Wifi struct {
				IPAddress  string `json:"ipAddress"`
				MACAddress string `json:"macAddress"`
			} `json:"wifi"`
		} `json:"connectivity"`
	} `json:"data"`
}

// GetAndroidAppsResponse declare response for get android apps request.
type GetAndroidAppsResponse struct {
	Data []struct { //nolint:govet
		ID          string `json:"id"`
		PackageName string `json:"packageName"`
		VersionCode int    `json:"versionCode"`
		Description string `json:"description"`
		Label       string `json:"label"`
		VersionName string `json:"versionName"`
		Status      string `json:"status"`
	} `json:"data"`
}

// ScheduleActionRequest declare structure for schedule action request.
type ScheduleActionRequest struct { //nolint:govet
	TerminalIDs   []string `json:"terminalIds"`
	StoreID       string   `json:"storeId"`
	ScheduledAt   string   `json:"scheduledAt"`
	ActionDetails struct {
		Type  string `json:"type"`
		AppID string `json:"appId"`
	} `json:"actionDetails"`
}

// ScheduleActionResponse declare response for schedule action request.
type ScheduleActionResponse struct {
	ActionDetails struct {
		AppID string `json:"appId"`
		Type  string `json:"type"`
	} `json:"actionDetails"`
	ScheduledAt string `json:"scheduledAt"`
	StoreID     string `json:"storeId"`
	Items       []struct {
		ID         string `json:"id"`
		TerminalID string `json:"terminalId"`
	} `json:"items"`
	TerminalsWithErrors struct {
	} `json:"terminalsWithErrors"`
	TotalScheduled int `json:"totalScheduled"`
	TotalErrors    int `json:"totalErrors"`
}
