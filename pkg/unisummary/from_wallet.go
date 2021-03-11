package unisummary

import (
	"encoding/json"
	"fmt"
	"time"
)

func fromWalletAddress(us UniswapSummaryRequest) []UniswapSummaryResponse {

	normalTransactions := fetchAllNormalTransactions(us)
	fmt.Println(normalTransactions)

	/*
		normalTransactions = filterNormalTransactions(normalTransactions)
		transactions := fromNormalTransactions(normalTransactions)
	*/

	/*
		tokenTransactions := fetchAllTokenTransactions(walletAddress)
		tokenTransactions = filterTokenTransactions(tokenTransactions, normalHashes, walletAddress)
		tokenTransactions = groupByHash(tokenTransactions)

		internalTransactions := fetchAllInternalTransactions(walletAddress)
		internalTransactions = filterInternalTransactions(internalTransactions, normalHashes)
	*/

	return []UniswapSummaryResponse{}
}

func fetchAllNormalTransactions(us UniswapSummaryRequest) EtherscanNormalTransactionsResponse {
	endpoint := fmt.Sprintf(us.EtherscanNormalTransactionsEndpoint, us.EtherscanApiKey, us.UserAddress)
	responseBody := callEndpoint(endpoint)
	var response EtherscanNormalTransactionsResponse
	err := json.Unmarshal([]byte(responseBody), &response)
	handleError(err)
	return response
}

type Transactions []Transaction

func (ts Transactions) Hashes() []string {
	hashes := []string{}
	for _, t := range ts {
		hashes = append(hashes, t.Hash)
	}
	return hashes
}

type Transaction struct {
	Hash              string
	GasUsed           float64
	GasPrice          float64
	Date              time.Time
	TokenTransactions []TokenTransaction
}

type SendOrReceive string

var send = SendOrReceive("send")
var receive = SendOrReceive("receive")

type TokenTransaction struct {
	TokenSymbol              string
	TokenDecimal             int
	ContractAddress          string
	Value                    float64
	SendOrReceive            SendOrReceive
	IsLiquidityProviderToken bool
}
