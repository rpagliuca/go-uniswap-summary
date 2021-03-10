# Go Uniswap Summary

# About
* This module provides the package `unisummary` to be used within other Golang apps
* It calculates useful information about Uniswap financial positions for liquidity providers:
    * Yearly profit rate
    * Divergence (impermanent) loss
    * Accrued fees
    * Current balance
* Blockchain data is fetched from Etherscan

# Usage
* Importing the package
```
import "github.com/rpagliuca/go-uniswap-summary/pkg/unisummary"
```
* See example at `cmd/example/example.go`
