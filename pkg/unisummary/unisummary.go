package unisummary

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UniswapSummaryRequest struct {
	EtherscanApiKey                       string
	EtherscanSupplyEndpoint               string
	EtherscanBalanceEndpoint              string
	EtherscanNormalTransactionsEndpoint   string
	EtherscanTokenTransactionsEndpoint    string
	EtherscanInternalTransactionsEndpoint string
	UserAddress                           string
	LiquidityProviderTokens               []LiquidityProviderPosition
}

func NewUniswapSummaryRequest(key string, userAddress string, lpTokens []LiquidityProviderPosition) *UniswapSummaryRequest {
	return &UniswapSummaryRequest{
		EtherscanApiKey:                       key,
		EtherscanSupplyEndpoint:               ETHERSCAN_ENDPOINT_SUPPLY,
		EtherscanBalanceEndpoint:              ETHERSCAN_ENDPOINT_BALANCE,
		EtherscanNormalTransactionsEndpoint:   ETHERSCAN_WALLET_NORMAL_TRANSACTIONS,
		EtherscanInternalTransactionsEndpoint: ETHERSCAN_WALLET_INTERNAL_TRANSACTIONS,
		EtherscanTokenTransactionsEndpoint:    ETHERSCAN_WALLET_ERC20_TRANSACTIONS,
		UserAddress:                           userAddress,
		LiquidityProviderTokens:               lpTokens,
	}
}

// Global client for HTTP keep-alive
var client = &http.Client{}

type Token struct {
	Id       string
	Address  string
	Decimals int
}

type LiquidityProviderPosition struct {
	Pair                  Token
	Token1                Token
	Token1InitialQuantity float64
	Token2                Token
	Token2InitialQuantity float64
	InitialDate           time.Time
}

type UniswapSummaryResponse struct {
	Token               LiquidityProviderPosition
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
	InitialPrice        float64
	FinalPrice          float64
	DivergenceLoss      float64
	AccruedProfit       float64
	DaysEllapsed        float64
	YearlyProfit        float64
}

func (us UniswapSummaryRequest) Do() []UniswapSummaryResponse {
	var wg sync.WaitGroup
	results := make([]UniswapSummaryResponse, len(us.LiquidityProviderTokens))
	for i, t := range us.LiquidityProviderTokens {
		wg.Add(1)
		go func(index int, thisT LiquidityProviderPosition) {

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

			results[index] = makeResponse(thisT, balance, supply, liquidity1, liquidity2)

			wg.Done()
		}(i, t)
	}
	wg.Wait()
	return results
}

func daysSince(start time.Time) float64 {
	end := time.Now()
	return end.Sub(start).Hours() / 24.0
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

func log(i ...interface{}) {
	if true {
		fmt.Println(i...)
	}
}

func makeResponse(thisT LiquidityProviderPosition, balance, supply, liquidity1, liquidity2 float64) UniswapSummaryResponse {

	token1FinalQuantity := balance / supply * liquidity1
	token2FinalQuantity := balance / supply * liquidity2
	initialK := thisT.Token1InitialQuantity * thisT.Token2InitialQuantity
	totalK := liquidity1 * liquidity2
	myK := math.Pow(balance/supply, 2) * totalK
	ratioK := myK / initialK
	percentageFees := (math.Pow(ratioK, 0.5) - 1.0) * 100.0
	token1Fee := token1FinalQuantity * (1.0 - 1.0/math.Pow(ratioK, 0.5))
	token2Fee := token2FinalQuantity * (1.0 - 1.0/math.Pow(ratioK, 0.5))
	initialPrice := thisT.Token1InitialQuantity / thisT.Token2InitialQuantity
	finalPrice := token1FinalQuantity / token2FinalQuantity
	priceRatio := finalPrice / initialPrice
	divergenceLoss := (2.0*math.Sqrt(priceRatio)/(1.0+priceRatio) - 1.0) * 100.0
	accruedProfit := percentageFees + divergenceLoss
	daysEllapsed := daysSince(thisT.InitialDate)
	yearlyProfit := (math.Pow(1.0+accruedProfit/100.0, 365.0/daysEllapsed) - 1.0) * 100.0

	response := UniswapSummaryResponse{
		Token:               thisT,
		Balance:             balance,
		Supply:              supply,
		Liquidity1:          liquidity1,
		Liquidity2:          liquidity2,
		InitialK:            initialK,
		Token1FinalQuantity: token1FinalQuantity,
		Token2FinalQuantity: token2FinalQuantity,
		Token1Increase:      token1FinalQuantity - thisT.Token1InitialQuantity,
		Token2Increase:      token2FinalQuantity - thisT.Token2InitialQuantity,
		InitialPrice:        initialPrice,
		FinalPrice:          finalPrice,
		DivergenceLoss:      divergenceLoss,
		AccruedProfit:       accruedProfit,
		DaysEllapsed:        daysEllapsed,
		YearlyProfit:        yearlyProfit,
		Token1Fee:           token1Fee,
		Token2Fee:           token2Fee,
		PercentageFees:      percentageFees,
		RatioK:              ratioK,
		MyK:                 myK,
		TotalK:              totalK,
	}

	return response
}
