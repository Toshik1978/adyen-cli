package adyen

// Adyen API types

// GetAccountHolderRequest declare account holder request.
type GetAccountHolderRequest struct {
	AccountHolderCode string `json:"accountHolderCode"`
}

// GetAccountHolderResponse declare account holder response.
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
			SplitConfigurationUUID string `json:"splitConfigurationUUID"`
			VirtualAccount         string `json:"virtualAccount"`
			StoreID                string `json:"store"`
		} `json:"storeDetails"`
	} `json:"accountHolderDetails"`
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
