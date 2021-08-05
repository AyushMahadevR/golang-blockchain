package controller

import (
	"encoding/json"
	"fmt"
	"go-blockchain/blockchain"
	"io/ioutil"
	"net/http"
)

func BlockChainController(w http.ResponseWriter, r *http.Request, b *blockchain.Blockchain) {
	bChainJSON, _ := json.Marshal(b)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bChainJSON)
}

func TranscationController(w http.ResponseWriter, r *http.Request, b *blockchain.Blockchain) {
	transaction := &blockchain.Transaction{}
	fmt.Printf("%+v", transaction)
	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, transaction)
	fmt.Printf("%+v", transaction)
	if err != nil {
		fmt.Println(err)
	}
	b.AddToPendingTransaction(transaction)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"message\":\"Transaction successfully added!\"}"))
}

func Mine(w http.ResponseWriter, r *http.Request, b *blockchain.Blockchain, nodeAddress string) {
	b.Mine(nodeAddress)
	bChainJSON, _ := json.Marshal(b)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bChainJSON)
}
