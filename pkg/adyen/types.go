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
