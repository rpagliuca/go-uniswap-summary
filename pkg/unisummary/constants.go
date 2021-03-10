package unisummary

const ETHERSCAN_ENDPOINT_SUPPLY = "https://api.etherscan.io/api?module=stats&apikey=%s&action=tokensupply&contractaddress=%s"
const ETHERSCAN_ENDPOINT_BALANCE = "https://api.etherscan.io/api?module=account&apikey=%s&action=tokenbalance&contractaddress=%s&address=%s&tag=latest"

var TOKEN_DAI = Token{"DAI", "0x6b175474e89094c44da98b954eedeac495271d0f", 18}
var TOKEN_WETH = Token{"WETH", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", 18}
var TOKEN_USDC = Token{"USDC", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", 6}
var TOKEN_DAI_USDC_LP = Token{"DAI_USDC_LP", "0xae461ca67b15dc8dc81ce7615e0320da1a9ab8d5", 18}
var TOKEN_DAI_WETH_LP = Token{"DAI_WETH_LP", "0xa478c2975ab1ea89e8196811f51a7b7ade33eb11", 18}
