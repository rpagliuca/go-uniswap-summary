package unisummary

const ETHERSCAN_ENDPOINT_SUPPLY = "https://api.etherscan.io/api?module=stats&apikey=%s&action=tokensupply&contractaddress=%s"
const ETHERSCAN_ENDPOINT_BALANCE = "https://api.etherscan.io/api?module=account&apikey=%s&action=tokenbalance&contractaddress=%s&address=%s&tag=latest"
const ETHERSCAN_WALLET_ERC20_TRANSACTIONS = "https://api.etherscan.io/api?module=account&apikey=%s&action=tokentx&address=%s&startblock=0&endblock=999999999&sort=asc"
const ETHERSCAN_WALLET_NORMAL_TRANSACTIONS = "https://api.etherscan.io/api?module=account&apikey=%s&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc"
const ETHERSCAN_WALLET_INTERNAL_TRANSACTIONS = "https://api.etherscan.io/api?module=account&apikey=%s&action=txlistinternal&address=%s&startblock=0&endblock=99999999&sort=asc"

var TOKEN_WETH = Token{"WETH", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", 18}

const UNISWAP_CONTRACT_ADDRESS = "0x7a250d5630b4cf539739df2c5dacb4c659f2488d"
const LIQUIDITY_PROVIDER_TOKEN_SYMBOL = "UNI-V2"
