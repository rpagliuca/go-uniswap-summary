package unisummary

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func FromWalletAddress(us *UniswapSummaryRequest) []LiquidityProviderPosition {

	normalTransactions := fetchAllNormalTransactions(us)
	transactions := processNormalTransactions(normalTransactions)

	tokenTransactions := fetchAllTokenTransactions(us)
	transactions = processTokenTransactions(us, transactions, tokenTransactions)

	internalTransactions := fetchAllInternalTransactions(us)
	transactions = processInternalTransactions(us, transactions, internalTransactions)

	transactions = normalizeAndRemoveSwaps(transactions)

	positions := makePositions(transactions)

	return positions
}

func makePositions(ts Transactions) []LiquidityProviderPosition {
	var positions []LiquidityProviderPosition
	for _, t := range ts {
		var token1, token2, pair int
		for i, tt := range t.TokenTransactions {
			if tt.IsLiquidityProviderToken == true {
				pair = i
			} else if token1 == 0 {
				token1 = i
			} else {
				token2 = i
			}
		}
		p := LiquidityProviderPosition{
			Pair: Token{
				Id: t.TokenTransactions[pair].TokenSymbol +
					" " + t.TokenTransactions[token1].TokenSymbol +
					" " + t.TokenTransactions[token2].TokenSymbol,
				Address:  t.TokenTransactions[pair].ContractAddress,
				Decimals: t.TokenTransactions[pair].TokenDecimal,
			},
			PairQuantity: parseTokenFloatQuantity(t.TokenTransactions[pair].Value, t.TokenTransactions[pair].TokenDecimal),
			Token1: Token{
				Id:       t.TokenTransactions[token1].TokenSymbol,
				Address:  t.TokenTransactions[token1].ContractAddress,
				Decimals: t.TokenTransactions[token1].TokenDecimal,
			},
			Token1InitialQuantity: -parseTokenFloatQuantity(t.TokenTransactions[token1].Value, t.TokenTransactions[token1].TokenDecimal),
			Token2: Token{
				Id:       t.TokenTransactions[token2].TokenSymbol,
				Address:  t.TokenTransactions[token2].ContractAddress,
				Decimals: t.TokenTransactions[token2].TokenDecimal,
			},
			Token2InitialQuantity: -parseTokenFloatQuantity(t.TokenTransactions[token2].Value, t.TokenTransactions[token2].TokenDecimal),
			InitialDate:           t.Date,
		}
		positions = append(positions, p)
	}
	return positions
}

func normalizeAndRemoveSwaps(ts Transactions) Transactions {
	swapsRemoved := Transactions{}
	for i, t := range ts {
		tokenTransactions := []TokenTransaction{}
		for _, tt := range t.TokenTransactions {
			exists := false
			factor := 1.0
			if tt.SendOrReceive == send {
				factor = -1.0
			}
			value := factor * tt.Value
			for j, a := range tokenTransactions {
				if a.ContractAddress == tt.ContractAddress {
					exists = true
					tokenTransactions[j].Value += value
					break
				}
			}
			if !exists {
				tt.Value = value
				tt.SendOrReceive = receive
				tokenTransactions = append(tokenTransactions, tt)
			}
		}
		// Only add liquidity pool actions (exactly 3 tokenTransactions)
		if len(tokenTransactions) == 3 {
			ts[i].TokenTransactions = tokenTransactions
			swapsRemoved = append(swapsRemoved, ts[i])
		}
	}

	return swapsRemoved
}

func processInternalTransactions(us *UniswapSummaryRequest, ts Transactions, r EtherscanInternalTransactionsResponse) Transactions {
	for i, t := range ts {
		for _, tt := range r.Result {
			if t.Hash == tt.Hash {
				if tt.From == UNISWAP_CONTRACT_ADDRESS {
					tokenTransaction := TokenTransaction{
						TokenSymbol:              "WETH",
						TokenDecimal:             18,
						ContractAddress:          TOKEN_WETH.Address,
						Value:                    toFloat(tt.Value),
						SendOrReceive:            receive,
						IsLiquidityProviderToken: false,
					}
					ts[i].TokenTransactions = append(ts[i].TokenTransactions, tokenTransaction)
				}
			}
		}
	}
	return ts
}

func processTokenTransactions(us *UniswapSummaryRequest, ts Transactions, r EtherscanTokenTransactionsResponse) Transactions {
	for i, t := range ts {
		for _, tt := range r.Result {
			if t.Hash == tt.Hash {

				isLpToken := false
				if tt.TokenSymbol == LIQUIDITY_PROVIDER_TOKEN_SYMBOL {
					isLpToken = true
				}

				sendOrReceive := send
				if icaseCompare(tt.To, us.UserAddress) {
					sendOrReceive = receive
				} else if !icaseCompare(tt.From, us.UserAddress) {
					fmt.Println(tt.To)
					fmt.Println(tt.From)
					panic("ERC20 token transaction TO or FROM should be equal the user wallet address")
				}

				tokenTransaction := TokenTransaction{
					TokenSymbol:              tt.TokenSymbol,
					TokenDecimal:             toInt(tt.TokenDecimal),
					ContractAddress:          tt.ContractAddress,
					Value:                    toFloat(tt.Value),
					SendOrReceive:            sendOrReceive,
					IsLiquidityProviderToken: isLpToken,
				}

				ts[i].TokenTransactions = append(ts[i].TokenTransactions, tokenTransaction)
			}
		}
	}
	return ts
}

func processNormalTransactions(r EtherscanNormalTransactionsResponse) Transactions {
	ts := Transactions{}
	for _, t := range r.Result {
		if t.IsError == "0" && t.TxReceiptStatus == "1" {
			if t.To == UNISWAP_CONTRACT_ADDRESS {
				tokenTransactions := []TokenTransaction{}
				if t.Value != "0" {
					tokenTransaction := TokenTransaction{
						TokenSymbol:              "WETH",
						TokenDecimal:             18,
						ContractAddress:          TOKEN_WETH.Address,
						Value:                    toFloat(t.Value),
						SendOrReceive:            send,
						IsLiquidityProviderToken: false,
					}
					tokenTransactions = append(tokenTransactions, tokenTransaction)
				}
				transaction := Transaction{
					Hash:              t.Hash,
					GasUsed:           toFloat(t.GasUsed),
					GasPrice:          toFloat(t.GasPrice),
					Date:              toTime(t.TimeStamp),
					TokenTransactions: tokenTransactions,
				}
				ts = append(ts, transaction)
			}
		}

	}
	return ts
}

func toInt(str string) int {
	f, err := strconv.ParseInt(str, 10, 64)
	handleError(err)
	return int(f)
}

func toFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	handleError(err)
	return f
}

func toTime(timestamp string) time.Time {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	handleError(err)
	return time.Unix(i, 0)
}

func fetchAllInternalTransactions(us *UniswapSummaryRequest) EtherscanInternalTransactionsResponse {
	endpoint := fmt.Sprintf(us.EtherscanInternalTransactionsEndpoint, us.EtherscanApiKey, us.UserAddress)
	responseBody := callEndpoint(endpoint)
	var response EtherscanInternalTransactionsResponse
	err := json.Unmarshal([]byte(responseBody), &response)
	handleError(err)
	return response
}

func fetchAllTokenTransactions(us *UniswapSummaryRequest) EtherscanTokenTransactionsResponse {
	endpoint := fmt.Sprintf(us.EtherscanTokenTransactionsEndpoint, us.EtherscanApiKey, us.UserAddress)
	responseBody := callEndpoint(endpoint)
	var response EtherscanTokenTransactionsResponse
	err := json.Unmarshal([]byte(responseBody), &response)
	handleError(err)
	return response
}

func fetchAllNormalTransactions(us *UniswapSummaryRequest) EtherscanNormalTransactionsResponse {
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

func icaseCompare(a, b string) bool {
	return strings.ToLower(a) == strings.ToLower(b)
}

func parseTokenFloatQuantity(quantity float64, decimals int) float64 {
	return quantity * math.Pow(10, -float64(decimals))
}
