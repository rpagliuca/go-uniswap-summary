package main

import (
	"encoding/json"
	"fmt"
	"os"

	us "github.com/rpagliuca/go-uniswap-summary/pkg/unisummary"
)

var etherscanApiKey = os.Getenv("ETHERSCAN_API_KEY")
var userAddress = os.Getenv("USER_ADDRESS")

func main() {

	if etherscanApiKey == "" || userAddress == "" {
		fmt.Println("Environment variables ETHERSCAN_API_KEY and USER_ADDRESS are required.")
		os.Exit(1)
	}

	req := us.NewUniswapSummaryRequest(
		etherscanApiKey,
		userAddress,
		[]us.LiquidityProviderPosition{},
	)

	resp := us.FromWalletAddress(req)

	jsonBytes, err := json.MarshalIndent(resp, "", "    ")
	handleError(err)

	fmt.Println(string(jsonBytes))

}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
