package unisummary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"
)

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
		throttleRequest(attempts)
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
var THROTTLE_DURATION = 1000 * time.Millisecond
var MAX_ATTEMPTS = 3

func throttleRequest(attempts int) {
	ellapsed := time.Now().UnixNano() - LAST_FAILURE_TIME
	left := THROTTLE_DURATION - time.Duration(ellapsed)
	if left < 0 {
		return
	}
	wait := float64(left) * math.Pow(2.0, float64(attempts))
	log(fmt.Sprintf("Sleeping for %0.2f milliseconds ", wait/1e6))
	// Exponential backoff
	time.Sleep(time.Duration(wait) * time.Nanosecond)
	log("Finished sleeping...")
}

func shouldRetry(attempts int) bool {
	LAST_FAILURE_TIME = time.Now().UnixNano()
	if attempts < MAX_ATTEMPTS {
		return true
	}
	return false
}
