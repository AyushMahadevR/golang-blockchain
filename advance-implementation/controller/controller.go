package controller

import (
	"encoding/json"
	"fmt"
	"go-blockchain/blockchain"
	"io"
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
	reqBody := getBodyAsBytes(r.Body)
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

func RegisterNewNode(w http.ResponseWriter, r *http.Request, bitcoin *blockchain.Blockchain, currentNodeUrl string) {
	request := make(map[string]string)
	reqBody := getBodyAsBytes(r.Body)
	_ = json.Unmarshal(reqBody, &request)
	newNodeUrl := request["newNodeUrl"]
	isNotCurrentNodeUrl := newNodeUrl != currentNodeUrl
	newNodeUrlNotPresent := !stringInSlice(newNodeUrl, bitcoin.NetworkNodes)
	if isNotCurrentNodeUrl && newNodeUrlNotPresent {
		bitcoin.NetworkNodes = append(bitcoin.NetworkNodes, newNodeUrl)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"message\":\"Node: " + newNodeUrl + " successfully added!\"}"))
}

func getBodyAsBytes(body io.ReadCloser) []byte {
	byteData, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte("{}")
	}
	return byteData
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
