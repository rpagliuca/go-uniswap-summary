package unisummary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UniswapSummaryRequest struct {
	EtherscanApiKey          string
	EtherscanSupplyEndpoint  string
	EtherscanBalanceEndpoint string
	UserAddress              string
	LiquidityProviderTokens  []LiquidityProviderToken
}

func NewUniswapSummaryRequest(key string, userAddress string, lpTokens []LiquidityProviderToken) *UniswapSummaryRequest {
	return &UniswapSummaryRequest{
		EtherscanApiKey:          key,
		EtherscanSupplyEndpoint:  ETHERSCAN_ENDPOINT_SUPPLY,
		EtherscanBalanceEndpoint: ETHERSCAN_ENDPOINT_BALANCE,
		UserAddress:              userAddress,
		LiquidityProviderTokens:  lpTokens,
	}
}

// Global client for HTTP keep-alive
var client = &http.Client{}

type Token struct {
	Id       string
	Address  string
	Decimals int
}

type LiquidityProviderToken struct {
	Pair                  Token
	Token1                Token
	Token1InitialQuantity float64
	Token2                Token
	Token2InitialQuantity float64
}

type UniswapSummaryResponse struct {
	Token               LiquidityProviderToken
	Balance             float64
	Supply              float64
	Liquidity1          float64
	Liquidity2          float64
	TotalK              float64
	MyK                 float64
	InitialK            float64
	Token1FinalQuantity float64
	Token2FinalQuantity float64
	Token1Increase      float64
	Token2Increase      float64
	Token1Fee           float64
	Token2Fee           float64
	RatioK              float64
	PercentageFees      float64
}

func (us UniswapSummaryRequest) Do() []UniswapSummaryResponse {
	var wg sync.WaitGroup
	results := make([]UniswapSummaryResponse, len(us.LiquidityProviderTokens))
	for i, t := range us.LiquidityProviderTokens {
		wg.Add(1)
		go func(index int, thisT LiquidityProviderToken) {

			var balance, supply, liquidity1, liquidity2 float64

			var wg2 sync.WaitGroup

			wg2.Add(1)
			go func() {
				balance = parseTokenQuantity(getBalance(us, thisT.Pair.Address, us.UserAddress), thisT.Pair.Decimals)
				wg2.Done()
			}()

			wg2.Add(1)
			go func() {
				supply = parseTokenQuantity(getSupply(us, thisT.Pair.Address), thisT.Pair.Decimals)
				wg2.Done()
			}()

			wg2.Add(1)
			go func() {
				liquidity1 = parseTokenQuantity(getBalance(us, thisT.Token1.Address, thisT.Pair.Address), thisT.Token1.Decimals)
				wg2.Done()
			}()

			wg2.Add(1)
			go func() {
				liquidity2 = parseTokenQuantity(getBalance(us, thisT.Token2.Address, thisT.Pair.Address), thisT.Token2.Decimals)
				wg2.Done()
			}()

			wg2.Wait()

			tokenState := UniswapSummaryResponse{
				Token:               thisT,
				Balance:             balance,
				Supply:              supply,
				Liquidity1:          liquidity1,
				Liquidity2:          liquidity2,
				InitialK:            thisT.Token1InitialQuantity * thisT.Token2InitialQuantity,
				Token1FinalQuantity: balance / supply * liquidity1,
				Token2FinalQuantity: balance / supply * liquidity2,
				Token1Increase:      balance/supply*liquidity1 - thisT.Token1InitialQuantity,
				Token2Increase:      balance/supply*liquidity2 - thisT.Token2InitialQuantity,
			}

			totalK, myK := calculateK(tokenState)

			tokenState.TotalK = totalK
			tokenState.MyK = myK
			tokenState.RatioK = tokenState.MyK / tokenState.InitialK
			tokenState.PercentageFees = (math.Pow(tokenState.RatioK, 0.5) - 1.0) * 100.0

			tokenState.Token1Fee = tokenState.Token1FinalQuantity * (1.0 - 1.0/math.Pow(tokenState.RatioK, 0.5))
			tokenState.Token2Fee = tokenState.Token2FinalQuantity * (1.0 - 1.0/math.Pow(tokenState.RatioK, 0.5))

			results[index] = tokenState
			wg.Done()
		}(i, t)
	}
	wg.Wait()
	return results
}

func getBalance(us UniswapSummaryRequest, tokenAddress string, walletAddress string) string {
	endpoint := fmt.Sprintf(us.EtherscanBalanceEndpoint, us.EtherscanApiKey, tokenAddress, walletAddress)
	result := getResult(endpoint)
	return result
}

func getSupply(us UniswapSummaryRequest, tokenAddress string) string {
	endpoint := fmt.Sprintf(us.EtherscanSupplyEndpoint, us.EtherscanApiKey, tokenAddress)
	result := getResult(endpoint)
	return result
}

func getResult(endpoint string) string {
	attempts := 0
	var result string
	for {
		throttleRequest()
		log(fmt.Sprintf("Fetching endpoint %s...", endpoint))
		resp, err := client.Get(endpoint)
		handleError(err)
		body, err := ioutil.ReadAll(resp.Body)
		handleError(err)
		var data map[string]string
		err = json.Unmarshal(body, &data)
		handleError(err)
		if status, ok := data["status"]; ok && status != "1" {
			if shouldRetry(attempts) {
				attempts++
				continue
			}
			panic(fmt.Sprintf("Status for endpoint %s should be 1", endpoint))
		}
		if _, ok := data["result"]; !ok {
			if shouldRetry(attempts) {
				attempts++
				continue
			}
			panic(fmt.Sprintf("Expecting property `result` for endpoint %s", endpoint))
		}
		result = data["result"]
		break
	}
	return result
}

var LAST_FAILURE_TIME = int64(0)
var THROTTLE_DURATION = 1500 * time.Millisecond
var MAX_ATTEMPTS = 3

func throttleRequest() {
	ellapsed := time.Now().UnixNano() - LAST_FAILURE_TIME
	left := THROTTLE_DURATION - time.Duration(ellapsed)
	if left < 0 {
		return
	}
	log(fmt.Sprintf("Sleeping for %s nanoseconds...", left))
	time.Sleep(time.Duration(left) * time.Nanosecond)
	log("Finished sleeping...")
}

func shouldRetry(attempts int) bool {
	LAST_FAILURE_TIME = time.Now().UnixNano()
	if attempts < MAX_ATTEMPTS {
		return true
	}
	return false
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseTokenQuantity(quantity string, decimals int) float64 {
	q, err := strconv.ParseFloat(quantity, 64)
	handleError(err)
	return q * math.Pow(10, -float64(decimals))
}

func calculateK(tokenState UniswapSummaryResponse) (float64, float64) {
	totalK := tokenState.Liquidity1 * tokenState.Liquidity2
	myK := math.Pow(tokenState.Balance/tokenState.Supply, 2) * totalK
	return totalK, myK
}

func log(i ...interface{}) {
	if false {
		fmt.Println(i...)
	}
}
