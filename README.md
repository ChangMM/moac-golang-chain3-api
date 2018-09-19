# Moac-Golang-Chain3-Api
[![Build Status](https://travis-ci.org/onrik/ethrpc.svg?branch=master)](https://travis-ci.org/onrik/ethrpc)
[![Coverage Status](https://coveralls.io/repos/github/onrik/ethrpc/badge.svg?branch=master)](https://coveralls.io/github/onrik/ethrpc?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/onrik/ethrpc)](https://goreportcard.com/report/github.com/onrik/ethrpc)
[![GoDoc](https://godoc.org/github.com/onrik/ethrpc?status.svg)](https://godoc.org/github.com/onrik/ethrpc)


Golang client for MOAC [Chain3 JSON RPC API](https://github.com/MOACChain/chain3).

- [x] chain3_clientVersion
- [x] chain3_sha3
- [x] net_version
- [x] net_peerCount
- [x] net_listening
- [x] mc_protocolVersion
- [x] mc_syncing
- [x] mc_coinbase
- [x] mc_mining
- [x] mc_hashrate
- [x] mc_gasPrice
- [x] mc_accounts
- [x] mc_blockNumber
- [x] mc_getBalance
- [x] mc_getStorageAt
- [x] mc_getTransactionCount
- [x] mc_getBlockTransactionCountByHash
- [x] mc_getBlockTransactionCountByNumber
- [x] mc_getUncleCountByBlockHash
- [x] mc_getUncleCountByBlockNumber
- [x] mc_getCode
- [x] mc_sign
- [x] mc_sendTransaction
- [x] mc_sendRawTransaction
- [x] mc_call
- [x] mc_estimateGas
- [x] mc_getBlockByHash
- [x] mc_getBlockByNumber
- [x] mc_getTransactionByHash
- [x] mc_getTransactionByBlockHashAndIndex
- [x] mc_getTransactionByBlockNumberAndIndex
- [x] mc_getTransactionReceipt
- [x] mc_getCompilers
- [x] mc_newFilter
- [x] mc_newBlockFilter
- [x] mc_newPendingTransactionFilter
- [x] mc_uninstallFilter
- [x] mc_getFilterChanges
- [x] mc_getFilterLogs
- [x] mc_getLogs


##### Usage:
```go
package main

import (
    "fmt"
    "log"

    "github.com/dacelee/moac-golang-chain3-api"
)

func main() {
    client := moacrpc.New("http://127.0.0.1:8545")

    version, err := client.Chain3ClientVersion()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(version)
    ```
    ----------------
#####  Golang Usage:
    // Send 1 MOAC （demo）
   ``` txid, err := client.MoacSendTransaction(moacrpc.T{
        From:  "0x6247cf0412c6462da2a51d05139e2a3c6c630f0a",
        To:    "0xcfa202c4268749fbb5136f2b68f7402984ed444b",
        Value: moacrpc.Moac1(),
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(txid)
}
```
----------------
##Donations

MOAC: 0x53be4cb8f27152893b448f9f569624afd1a97e0c


