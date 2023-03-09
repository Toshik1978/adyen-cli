package adyen

// GetAccountHolderRequest declare account holder request Adyen type.
type GetAccountHolderRequest struct {
	AccountHolderCode string `json:"accountHolderCode"`
}

// GetAccountHolderResponse declare account holder response Adyen type.
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

// UpdateAccountHolderRequest declare account holder update request Adyen type.
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
