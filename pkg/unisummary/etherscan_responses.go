package unisummary

type EtherscanTokenTransactionsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		From              string `json:"from"`
		ContractAddress   string `json:"contractAddress"`
		To                string `json:"to"`
		Value             string `json:"value"`
		TokenName         string `json:"tokenName"`
		TokenSymbol       string `json:"tokenSymbol"`
		TokenDecimal      string `json:"tokenDecimal"`
		TransactionIndex  string `json:"transactionIndex"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		GasUsed           string `json:"gasUsed"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		Input             string `json:"input"`
		Confirmations     string `json:"confirmations"`
	}
}

type EtherscanInternalTransactionsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber     string `json:"blockNumber"`
		TimeStamp       string `json:"timeStamp"`
		Hash            string `json:"hash"`
		From            string `json:"from"`
		To              string `json:"to"`
		Value           string `json:"value"`
		ContractAddress string `json:"contractAddress"`
		Input           string `json:"input"`
		Type            string `json:"type"`
		Gas             string `json:"gas"`
		GasUsed         string `json:"gasUsed"`
		TraceId         string `json:"traceId"`
		IsError         string `json:"isError"`
		ErrCode         string `json:"errCode"`
	}
}

type EtherscanNormalTransactionsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		TransactionIndex  string `json:"transactionIndex"`
		From              string `json:"from"`
		To                string `json:"to"`
		Value             string `json:"value"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		IsError           string `json:"isError"`
		Txreceipt_status  string `json:"txreceipt_status"`
		Input             string `json:"input"`
		ContractAddress   string `json:"contractAddress"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		GasUsed           string `json:"gasUsed"`
		Confirmations     string `json:"confirmations"`
	}
}
