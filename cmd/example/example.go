package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	us "github.com/rpagliuca/go-uniswap-summary/pkg/unisummary"
)

var etherscanApiKey = os.Getenv("ETHERSCAN_API_KEY")
var userAddress = os.Getenv("USER_ADDRESS")

func main() {

	if etherscanApiKey == "" || userAddress == "" {
		fmt.Println("Environment variables ETHERSCAN_API_KEY and USER_ADDRESS are required.")
		os.Exit(1)
	}

	utc := time.FixedZone("UTC", 0)
	liquidityProviderTokens := []us.LiquidityProviderPosition{
		{us.TOKEN_DAI_WETH_LP, us.TOKEN_DAI, 100, us.TOKEN_WETH, 0.075, time.Date(2021, 1, 28, 18, 29, 0, 0, utc)},
		{us.TOKEN_DAI_USDC_LP, us.TOKEN_DAI, 100, us.TOKEN_USDC, 100, time.Date(2021, 1, 26, 12, 12, 0, 0, utc)},
	}

	req := us.NewUniswapSummaryRequest(
		etherscanApiKey,
		userAddress,
		liquidityProviderTokens,
	)

	resp := req.Do()

	jsonBytes, err := json.MarshalIndent(resp, "", "    ")
	handleError(err)

	fmt.Println(string(jsonBytes))

}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
