package adyen

import "time"

// Adyen API types

// Error declare Adyen API error response.
type Error struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	ErrorCode string `json:"errorCode"`
	Status    int    `json:"status"`
}

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
	ID              string   `json:"id"`
	MerchantID      string   `json:"merchantId"`
	BusinessLineIDs []string `json:"businessLineIds"`
	Reference       string   `json:"reference"`
	Status          string   `json:"status"`
}

// SearchStoresResponse declare get all stores response.
type SearchStoresResponse struct {
	Data       []GetStoreResponse `json:"data"`
	ItemsTotal int64              `json:"itemsTotal"`
}

// AddPaymentMethodRequest declare add payment method request.
type AddPaymentMethodRequest struct {
	Type           string   `json:"type"`
	BusinessLineID string   `json:"businessLineId"`
	StoreIDs       []string `json:"storeIds"`
	Currencies     []string `json:"currencies"`
}

// AddPaymentMethodResponse declare add payment method response.
type AddPaymentMethodResponse struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`
	BusinessLineID     string   `json:"businessLineId"`
	StoreIDs           []string `json:"storeIds"`
	Currencies         []string `json:"currencies"`
	Enabled            bool     `json:"enabled"`
	Allowed            bool     `json:"allowed"`
	VerificationStatus string   `json:"verificationStatus"`
}

// GetBalanceAccountResponse declare balance account information.
type GetBalanceAccountResponse struct {
	ID                           string `json:"id"`
	AccountHolderID              string `json:"accountHolderId"`
	DefaultCurrencyCode          string `json:"defaultCurrencyCode"`
	Description                  string `json:"description"`
	PlatformPaymentConfiguration struct {
		SalesDayClosingTime string `json:"salesDayClosingTime"`
		SettlementDelayDays int    `json:"settlementDelayDays"`
	} `json:"platformPaymentConfiguration"`
	TimeZone string `json:"timeZone"`
	Balances []struct {
		Available int    `json:"available"`
		Balance   int    `json:"balance"`
		Currency  string `json:"currency"`
		Pending   int    `json:"pending"`
		Reserved  int    `json:"reserved"`
	} `json:"balances"`
	Status string `json:"status"`
}

// GetBalanceAccountHolderResponse declare the response for balance get account holder request.
type GetBalanceAccountHolderResponse struct {
	ID                    string `json:"id"`
	PrimaryBalanceAccount string `json:"primaryBalanceAccount"`
	BalancePlatform       string `json:"balancePlatform"`
	Description           string `json:"description"`
	LegalEntityID         string `json:"legalEntityId"`
	Capabilities          struct {
		ReceiveFromPlatformPayments struct {
			Enabled            bool   `json:"enabled"`
			Requested          bool   `json:"requested"`
			Allowed            bool   `json:"allowed"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"receiveFromPlatformPayments"`
		ReceiveFromBalanceAccount struct {
			Enabled            bool   `json:"enabled"`
			Requested          bool   `json:"requested"`
			Allowed            bool   `json:"allowed"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"receiveFromBalanceAccount"`
		SendToBalanceAccount struct {
			Enabled            bool   `json:"enabled"`
			Requested          bool   `json:"requested"`
			Allowed            bool   `json:"allowed"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"sendToBalanceAccount"`
		SendToTransferInstrument struct {
			Enabled            bool   `json:"enabled"`
			Requested          bool   `json:"requested"`
			Allowed            bool   `json:"allowed"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"sendToTransferInstrument"`
		ReceivePayments struct {
			Enabled            bool   `json:"enabled"`
			Requested          bool   `json:"requested"`
			Allowed            bool   `json:"allowed"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"receivePayments"`
	} `json:"capabilities"`
	Status string `json:"status"`
}

// GetLegalEntityResponse declare legal entity structure.
type GetLegalEntityResponse struct {
	ID           string `json:"id"`
	Capabilities struct {
		ReceiveFromBalanceAccount struct {
			Allowed            bool   `json:"allowed"`
			Requested          bool   `json:"requested"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"receiveFromBalanceAccount"`
		ReceivePayments struct {
			Allowed            bool   `json:"allowed"`
			Requested          bool   `json:"requested"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"receivePayments"`
		SendToTransferInstrument struct {
			Allowed            bool   `json:"allowed"`
			Requested          bool   `json:"requested"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"sendToTransferInstrument"`
		ReceiveFromPlatformPayments struct {
			Allowed            bool   `json:"allowed"`
			Requested          bool   `json:"requested"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"receiveFromPlatformPayments"`
		SendToBalanceAccount struct {
			Allowed            bool   `json:"allowed"`
			Requested          bool   `json:"requested"`
			VerificationStatus string `json:"verificationStatus"`
		} `json:"sendToBalanceAccount"`
	} `json:"capabilities"`
	EntityAssociations []struct {
		AssociatorID  string `json:"associatorId"`
		EntityType    string `json:"entityType"`
		JobTitle      string `json:"jobTitle,omitempty"`
		LegalEntityID string `json:"legalEntityId"`
		Name          string `json:"name"`
		Type          string `json:"type"`
	} `json:"entityAssociations"`
	Organization struct {
		DoingBusinessAs   string `json:"doingBusinessAs"`
		LegalName         string `json:"legalName"`
		RegisteredAddress struct {
			City            string `json:"city"`
			Country         string `json:"country"`
			PostalCode      string `json:"postalCode"`
			StateOrProvince string `json:"stateOrProvince"`
			Street          string `json:"street"`
		} `json:"registeredAddress"`
		RegistrationNumber string `json:"registrationNumber"`
		TaxInformation     []struct {
			Country string `json:"country"`
			Number  string `json:"number"`
			Type    string `json:"type"`
		} `json:"taxInformation"`
		Type string `json:"type"`
	} `json:"organization"`
	Type            string `json:"type"`
	DocumentDetails []struct {
		ID               string    `json:"id"`
		Active           bool      `json:"active"`
		Description      string    `json:"description"`
		FileName         string    `json:"fileName"`
		ModificationDate time.Time `json:"modificationDate"`
		Pages            []struct {
			PageName   string `json:"pageName"`
			PageNumber int    `json:"pageNumber"`
			Type       string `json:"type"`
		} `json:"pages"`
		Type string `json:"type"`
	} `json:"documentDetails"`
	Documents []struct {
		ID string `json:"id"`
	} `json:"documents"`
	TransferInstruments []struct {
		ID                string `json:"id"`
		AccountIdentifier string `json:"accountIdentifier"`
		TrustedSource     bool   `json:"trustedSource"`
	} `json:"transferInstruments"`
}

// Sweep declare one sweep configuration.
type Sweep struct {
	ID       string `json:"id"`
	Schedule struct {
		CronExpression string `json:"cronExpression"`
		Type           string `json:"type"`
	} `json:"schedule"`
	Status       string `json:"status"`
	TargetAmount struct {
		Currency string `json:"currency"`
		Value    int    `json:"value"`
	} `json:"targetAmount"`
	TriggerAmount struct {
		Currency string `json:"currency"`
		Value    int    `json:"value"`
	} `json:"triggerAmount"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Counterparty struct {
		TransferInstrumentID string `json:"transferInstrumentId"`
	} `json:"counterparty"`
	Currency   string   `json:"currency"`
	Priorities []string `json:"priorities"`
}

// GetSweepsResponse declare all available sweeps per balance account.
type GetSweepsResponse struct {
	HasNext     bool    `json:"hasNext"`
	HasPrevious bool    `json:"hasPrevious"`
	Sweeps      []Sweep `json:"sweeps"`
}

// UpdateSweepRequest declare update sweep configuration request.
type UpdateSweepRequest struct {
	Counterparty struct {
		TransferInstrumentID string `json:"transferInstrumentId"`
	} `json:"counterparty"`
	Status string `json:"status"`
}

// Sweep configuration status

const (
	SweepActive = "active"
)

// SetSalesCloseTimeRequest declare change sales closing time request.
type SetSalesCloseTimeRequest struct {
	PlatformPaymentConfiguration struct {
		SalesDayClosingTime string `json:"salesDayClosingTime"`
		SettlementDelayDays int    `json:"settlementDelayDays"`
	} `json:"platformPaymentConfiguration"`
	TimeZone string `json:"timeZone"`
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

// SetOfflinePaymentsRequest declare request to update offline payments settings.
type SetOfflinePaymentsRequest struct {
	OfflineProcessing struct {
		ChipFloorLimit     int `json:"chipFloorLimit"`
		OfflineSwipeLimits []struct {
			Amount       int    `json:"amount"`
			CurrencyCode string `json:"currencyCode"`
		} `json:"offlineSwipeLimits"`
	} `json:"offlineProcessing"`
	StoreAndForward struct {
		MaxPayments int `json:"maxPayments"`
		MaxAmount   []struct {
			Amount       int    `json:"amount"`
			CurrencyCode string `json:"currencyCode"`
		} `json:"maxAmount"`
		SupportedCardTypes struct {
			Credit        bool `json:"credit"`
			Debit         bool `json:"debit"`
			DeferredDebit bool `json:"deferredDebit"`
			Prepaid       bool `json:"prepaid"`
			Unknown       bool `json:"unknown"`
		} `json:"supportedCardTypes"`
	} `json:"storeAndForward"`
}

// TerminalSettingsResponse declare response with terminal settings.
type TerminalSettingsResponse struct {
	CardholderReceipt struct {
		HeaderForAuthorizedReceipt string `json:"headerForAuthorizedReceipt"`
	} `json:"cardholderReceipt"`
	Gratuities []struct {
		Currency                string   `json:"currency"`
		UsePredefinedTipEntries bool     `json:"usePredefinedTipEntries"`
		PredefinedTipEntries    []string `json:"predefinedTipEntries"`
		AllowCustomAmount       bool     `json:"allowCustomAmount"`
	} `json:"gratuities"`
	Nexo struct {
		DisplayUrls struct {
			LocalUrls []struct {
				Password string `json:"password"`
				URL      string `json:"url"`
				Username string `json:"username"`
			} `json:"localUrls"`
			PublicUrls []struct {
				Password string `json:"password"`
				URL      string `json:"url"`
				Username string `json:"username"`
			} `json:"publicUrls"`
		} `json:"displayUrls"`
		EncryptionKey struct {
			Identifier string `json:"identifier"`
			Passphrase string `json:"passphrase"`
			Version    int    `json:"version"`
		} `json:"encryptionKey"`
		EventUrls struct {
			EventLocalUrls []struct {
				Password string `json:"password"`
				URL      string `json:"url"`
				Username string `json:"username"`
			} `json:"eventLocalUrls"`
			EventPublicUrls []struct {
				Password string `json:"password"`
				URL      string `json:"url"`
				Username string `json:"username"`
			} `json:"eventPublicUrls"`
		} `json:"eventUrls"`
		Notification struct {
			Category   string `json:"category"`
			Details    string `json:"details"`
			Enabled    bool   `json:"enabled"`
			ShowButton bool   `json:"showButton"`
			Title      string `json:"title"`
		} `json:"notification"`
	} `json:"nexo"`
	Opi struct {
		EnablePayAtTable bool `json:"enablePayAtTable"`
	} `json:"opi"`
	ReceiptPrinting struct {
		MerchantApproved        bool `json:"merchantApproved"`
		MerchantRefused         bool `json:"merchantRefused"`
		MerchantCancelled       bool `json:"merchantCancelled"`
		MerchantRefundApproved  bool `json:"merchantRefundApproved"`
		MerchantRefundRefused   bool `json:"merchantRefundRefused"`
		MerchantCaptureApproved bool `json:"merchantCaptureApproved"`
		MerchantCaptureRefused  bool `json:"merchantCaptureRefused"`
		MerchantVoid            bool `json:"merchantVoid"`
		ShopperApproved         bool `json:"shopperApproved"`
		ShopperRefused          bool `json:"shopperRefused"`
		ShopperCancelled        bool `json:"shopperCancelled"`
		ShopperRefundApproved   bool `json:"shopperRefundApproved"`
		ShopperRefundRefused    bool `json:"shopperRefundRefused"`
		ShopperCaptureApproved  bool `json:"shopperCaptureApproved"`
		ShopperCaptureRefused   bool `json:"shopperCaptureRefused"`
		ShopperVoid             bool `json:"shopperVoid"`
	} `json:"receiptPrinting"`
	Signature struct {
		AskSignatureOnScreen bool `json:"askSignatureOnScreen"`
		SkipSignature        bool `json:"skipSignature"`
	} `json:"signature"`
	Timeouts struct {
		FromActiveToSleep int `json:"fromActiveToSleep"`
	} `json:"timeouts"`
	Hardware struct {
		RestartHour int `json:"restartHour"`
	} `json:"hardware"`
	Connectivity struct {
		SimcardStatus string `json:"simcardStatus"`
	} `json:"connectivity"`
	OfflineProcessing struct {
		ChipFloorLimit     int `json:"chipFloorLimit"`
		OfflineSwipeLimits []struct {
			Amount       int    `json:"amount"`
			CurrencyCode string `json:"currencyCode"`
		} `json:"offlineSwipeLimits"`
	} `json:"offlineProcessing"`
	Passcodes struct {
		AdminMenuPin string `json:"adminMenuPin"`
		TxMenuPin    string `json:"txMenuPin"`
	} `json:"passcodes"`
	Standalone struct {
		EnableStandalone bool   `json:"enableStandalone"`
		CurrencyCode     string `json:"currencyCode"`
	} `json:"standalone"`
	StoreAndForward struct {
		MaxPayments int `json:"maxPayments"`
		MaxAmount   []struct {
			Amount       int    `json:"amount"`
			CurrencyCode string `json:"currencyCode"`
		} `json:"maxAmount"`
		SupportedCardTypes struct {
			Credit        bool `json:"credit"`
			Debit         bool `json:"debit"`
			DeferredDebit bool `json:"deferredDebit"`
			Prepaid       bool `json:"prepaid"`
			Unknown       bool `json:"unknown"`
		} `json:"supportedCardTypes"`
	} `json:"storeAndForward"`
	Payment struct {
		ContactlessCurrency string `json:"contactlessCurrency"`
	} `json:"payment"`
	Localization struct {
		Language string `json:"language"`
		Timezone string `json:"timezone"`
	} `json:"localization"`
	TerminalInstructions struct {
		AdyenAppRestart bool `json:"adyenAppRestart"`
	} `json:"terminalInstructions"`
}

// SearchTerminalsResponse declare response for search terminals request.
type SearchTerminalsResponse struct {
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

// SearchAndroidAppsResponse declare response for get android apps request.
type SearchAndroidAppsResponse struct {
	Data []struct {
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
type ScheduleActionRequest struct {
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
