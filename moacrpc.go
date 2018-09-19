package moacrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
)

// MoacError - moac error
type MoacError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err MoacError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type ethResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *MoacError       `json:"error"`
}

type ethRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// MoacRPC - Moacereum rpc client
type MoacRPC struct {
	url    string
	client httpClient
	log    logger
	Debug  bool
}

// New create new rpc client with given url
func New(url string, options ...func(rpc *MoacRPC)) *MoacRPC {
	rpc := &MoacRPC{
		url:    url,
		client: http.DefaultClient,
		log:    log.New(os.Stderr, "", log.LstdFlags),
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// NewMoacRPC create new rpc client with given url
func NewMoacRPC(url string, options ...func(rpc *MoacRPC)) *MoacRPC {
	return New(url, options...)
}

func (rpc *MoacRPC) call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// Call returns raw response of method call
func (rpc *MoacRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if rpc.Debug {
		rpc.log.Println(fmt.Sprintf("%s\nRequest: %s\nResponse: %s\n", method, body, data))
	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// RawCall returns raw response of method call (Deprecated)
func (rpc *MoacRPC) RawCall(method string, params ...interface{}) (json.RawMessage, error) {
	return rpc.Call(method, params...)
}

// Chain3ClientVersion returns the current client version.
func (rpc *MoacRPC) Chain3ClientVersion() (string, error) {
	var clientVersion string

	err := rpc.call("chain3_clientVersion", &clientVersion)
	return clientVersion, err
}

// Chain3Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (rpc *MoacRPC) Chain3Sha3(data []byte) (string, error) {
	var hash string

	err := rpc.call("chain3_sha3", &hash, fmt.Sprintf("0x%x", data))
	return hash, err
}

// NetVersion returns the current network protocol version.
func (rpc *MoacRPC) NetVersion() (string, error) {
	var version string

	err := rpc.call("net_version", &version)
	return version, err
}

// NetListening returns true if client is actively listening for network connections.
func (rpc *MoacRPC) NetListening() (bool, error) {
	var listening bool

	err := rpc.call("net_listening", &listening)
	return listening, err
}

// NetPeerCount returns number of peers currently connected to the client.
func (rpc *MoacRPC) NetPeerCount() (int, error) {
	var response string
	if err := rpc.call("net_peerCount", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacProtocolVersion returns the current moac protocol version.
func (rpc *MoacRPC) MoacProtocolVersion() (string, error) {
	var protocolVersion string

	err := rpc.call("mc_protocolVersion", &protocolVersion)
	return protocolVersion, err
}

// MoacSyncing returns an object with data about the sync status or false.
func (rpc *MoacRPC) MoacSyncing() (*Syncing, error) {
	result, err := rpc.RawCall("mc_syncing")
	if err != nil {
		return nil, err
	}
	syncing := new(Syncing)
	if bytes.Equal(result, []byte("false")) {
		return syncing, nil
	}
	err = json.Unmarshal(result, syncing)
	return syncing, err
}

// MoacCoinbase returns the client coinbase address
func (rpc *MoacRPC) MoacCoinbase() (string, error) {
	var address string

	err := rpc.call("mc_coinbase", &address)
	return address, err
}

// MoacMining returns true if client is actively mining new blocks.
func (rpc *MoacRPC) MoacMining() (bool, error) {
	var mining bool

	err := rpc.call("mc_mining", &mining)
	return mining, err
}

// MoacHashrate returns the number of hashes per second that the node is mining with.
func (rpc *MoacRPC) MoacHashrate() (int, error) {
	var response string

	if err := rpc.call("mc_hashrate", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacGasPrice returns the current price per gas in wei.
func (rpc *MoacRPC) MoacGasPrice() (big.Int, error) {
	var response string
	if err := rpc.call("mc_gasPrice", &response); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// MoacAccounts returns a list of addresses owned by client.
func (rpc *MoacRPC) MoacAccounts() ([]string, error) {
	accounts := []string{}

	err := rpc.call("mc_accounts", &accounts)
	return accounts, err
}

// MoacBlockNumber returns the number of most recent block.
func (rpc *MoacRPC) MoacBlockNumber() (int, error) {
	var response string
	if err := rpc.call("mc_blockNumber", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacGetBalance returns the balance of the account of given address in wei.
func (rpc *MoacRPC) MoacGetBalance(address, block string) (big.Int, error) {
	var response string
	if err := rpc.call("mc_getBalance", &response, address, block); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// MoacGetStorageAt returns the value from a storage position at a given address.
func (rpc *MoacRPC) MoacGetStorageAt(data string, position int, tag string) (string, error) {
	var result string

	err := rpc.call("mc_getStorageAt", &result, data, IntToHex(position), tag)
	return result, err
}

// MoacGetTransactionCount returns the number of transactions sent from an address.
func (rpc *MoacRPC) MoacGetTransactionCount(address, block string) (int, error) {
	var response string

	if err := rpc.call("mc_getTransactionCount", &response, address, block); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacGetBlockTransactionCountByHash returns the number of transactions in a block from a block matching the given block hash.
func (rpc *MoacRPC) MoacGetBlockTransactionCountByHash(hash string) (int, error) {
	var response string

	if err := rpc.call("mc_getBlockTransactionCountByHash", &response, hash); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacGetBlockTransactionCountByNumber returns the number of transactions in a block from a block matching the given block
func (rpc *MoacRPC) MoacGetBlockTransactionCountByNumber(number int) (int, error) {
	var response string

	if err := rpc.call("mc_getBlockTransactionCountByNumber", &response, IntToHex(number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacGetUncleCountByBlockHash returns the number of uncles in a block from a block matching the given block hash.
func (rpc *MoacRPC) MoacGetUncleCountByBlockHash(hash string) (int, error) {
	var response string

	if err := rpc.call("mc_getUncleCountByBlockHash", &response, hash); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacGetUncleCountByBlockNumber returns the number of uncles in a block from a block matching the given block number.
func (rpc *MoacRPC) MoacGetUncleCountByBlockNumber(number int) (int, error) {
	var response string

	if err := rpc.call("mc_getUncleCountByBlockNumber", &response, IntToHex(number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// MoacGetCode returns code at a given address.
func (rpc *MoacRPC) MoacGetCode(address, block string) (string, error) {
	var code string

	err := rpc.call("mc_getCode", &code, address, block)
	return code, err
}

// MoacSign signs data with a given address.
// Calculates an Moacereum specific signature with: sign(keccak256("\x19Moacereum Signed Message:\n" + len(message) + message)))
func (rpc *MoacRPC) MoacSign(address, data string) (string, error) {
	var signature string

	err := rpc.call("mc_sign", &signature, address, data)
	return signature, err
}

// MoacSendTransaction creates new message call transaction or a contract creation, if the data field contains code.
func (rpc *MoacRPC) MoacSendTransaction(transaction T) (string, error) {
	var hash string

	err := rpc.call("mc_sendTransaction", &hash, transaction)
	return hash, err
}

// MoacSendRawTransaction creates new message call transaction or a contract creation for signed transactions.
func (rpc *MoacRPC) MoacSendRawTransaction(data string) (string, error) {
	var hash string

	err := rpc.call("mc_sendRawTransaction", &hash, data)
	return hash, err
}

// MoacCall executes a new message call immediately without creating a transaction on the block chain.
func (rpc *MoacRPC) MoacCall(transaction T, tag string) (string, error) {
	var data string

	err := rpc.call("mc_call", &data, transaction, tag)
	return data, err
}

// MoacEstimateGas makes a call or transaction, which won't be added to the blockchain and returns the used gas, which can be used for estimating the used gas.
func (rpc *MoacRPC) MoacEstimateGas(transaction T) (int, error) {
	var response string

	err := rpc.call("mc_estimateGas", &response, transaction)
	if err != nil {
		return 0, err
	}

	return ParseInt(response)
}

func (rpc *MoacRPC) getBlock(method string, withTransactions bool, params ...interface{}) (*Block, error) {
	var response proxyBlock
	if withTransactions {
		response = new(proxyBlockWithTransactions)
	} else {
		response = new(proxyBlockWithoutTransactions)
	}

	err := rpc.call(method, response, params...)
	if err != nil {
		return nil, err
	}
	block := response.toBlock()

	return &block, nil
}

// MoacGetBlockByHash returns information about a block by hash.
func (rpc *MoacRPC) MoacGetBlockByHash(hash string, withTransactions bool) (*Block, error) {
	return rpc.getBlock("mc_getBlockByHash", withTransactions, hash, withTransactions)
}

// MoacGetBlockByNumber returns information about a block by block number.
func (rpc *MoacRPC) MoacGetBlockByNumber(number int, withTransactions bool) (*Block, error) {
	return rpc.getBlock("mc_getBlockByNumber", withTransactions, IntToHex(number), withTransactions)
}

func (rpc *MoacRPC) getTransaction(method string, params ...interface{}) (*Transaction, error) {
	transaction := new(Transaction)

	err := rpc.call(method, transaction, params...)
	return transaction, err
}

// MoacGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (rpc *MoacRPC) MoacGetTransactionByHash(hash string) (*Transaction, error) {
	return rpc.getTransaction("mc_getTransactionByHash", hash)
}

// MoacGetTransactionByBlockHashAndIndex returns information about a transaction by block hash and transaction index position.
func (rpc *MoacRPC) MoacGetTransactionByBlockHashAndIndex(blockHash string, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("mc_getTransactionByBlockHashAndIndex", blockHash, IntToHex(transactionIndex))
}

// MoacGetTransactionByBlockNumberAndIndex returns information about a transaction by block number and transaction index position.
func (rpc *MoacRPC) MoacGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("mc_getTransactionByBlockNumberAndIndex", IntToHex(blockNumber), IntToHex(transactionIndex))
}

// MoacGetTransactionReceipt returns the receipt of a transaction by transaction hash.
// Note That the receipt is not available for pending transactions.
func (rpc *MoacRPC) MoacGetTransactionReceipt(hash string) (*TransactionReceipt, error) {
	transactionReceipt := new(TransactionReceipt)

	err := rpc.call("mc_getTransactionReceipt", transactionReceipt, hash)
	if err != nil {
		return nil, err
	}

	return transactionReceipt, nil
}

// MoacGetCompilers returns a list of available compilers in the client.
func (rpc *MoacRPC) MoacGetCompilers() ([]string, error) {
	compilers := []string{}

	err := rpc.call("mc_getCompilers", &compilers)
	return compilers, err
}

// MoacNewFilter creates a new filter object.
func (rpc *MoacRPC) MoacNewFilter(params FilterParams) (string, error) {
	var filterID string
	err := rpc.call("mc_newFilter", &filterID, params)
	return filterID, err
}

// MoacNewBlockFilter creates a filter in the node, to notify when a new block arrives.
// To check if the state has changed, call MoacGetFilterChanges.
func (rpc *MoacRPC) MoacNewBlockFilter() (string, error) {
	var filterID string
	err := rpc.call("mc_newBlockFilter", &filterID)
	return filterID, err
}

// MoacNewPendingTransactionFilter creates a filter in the node, to notify when new pending transactions arrive.
// To check if the state has changed, call MoacGetFilterChanges.
func (rpc *MoacRPC) MoacNewPendingTransactionFilter() (string, error) {
	var filterID string
	err := rpc.call("mc_newPendingTransactionFilter", &filterID)
	return filterID, err
}

// MoacUninstallFilter uninstalls a filter with given id.
func (rpc *MoacRPC) MoacUninstallFilter(filterID string) (bool, error) {
	var res bool
	err := rpc.call("mc_uninstallFilter", &res, filterID)
	return res, err
}

// MoacGetFilterChanges polling method for a filter, which returns an array of logs which occurred since last poll.
func (rpc *MoacRPC) MoacGetFilterChanges(filterID string) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("mc_getFilterChanges", &logs, filterID)
	return logs, err
}

// MoacGetFilterLogs returns an array of all logs matching filter with given id.
func (rpc *MoacRPC) MoacGetFilterLogs(filterID string) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("mc_getFilterLogs", &logs, filterID)
	return logs, err
}

// MoacGetLogs returns an array of all logs matching a given filter object.
func (rpc *MoacRPC) MoacGetLogs(params FilterParams) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("mc_getLogs", &logs, params)
	return logs, err
}

// Moac1 returns 1 moac value (10^18 wei)
func (rpc *MoacRPC) Moac1() *big.Int {
	return Moac1()
}

// Moac1 returns 1 moac value (10^18 wei)
func Moac1() *big.Int {
	return big.NewInt(1000000000000000000)
}
