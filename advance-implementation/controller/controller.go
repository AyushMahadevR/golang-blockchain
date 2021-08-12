package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-blockchain/blockchain"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
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

func RegisterAndBroadcastNewNode(w http.ResponseWriter, r *http.Request, b *blockchain.Blockchain) {
	request := make(map[string]string)
	reqBody := getBodyAsBytes(r.Body)
	fmt.Print(string(reqBody))
	_ = json.Unmarshal(reqBody, &request)
	newNodeUrl := request["newNodeUrl"]
	addNewNodeIfNotExists(newNodeUrl, b)
	postBody, _ := json.Marshal(map[string]string{"newNodeUrl": newNodeUrl})
	responseBody := bytes.NewBuffer(postBody)
	waitGroup := sync.WaitGroup{}
	for _, nodesUrl := range b.NetworkNodes {
		waitGroup.Add(1)
		go func(nUrl string) {
			fmt.Println("url ", nUrl)
			resp, err := http.Post(nUrl+"/register-node", "application/json", responseBody)
			if err == nil {
				defer resp.Body.Close()
				reqBody := getBodyAsBytes(resp.Body)
				fmt.Println(string(reqBody))
			}
			waitGroup.Done()
		}(nodesUrl)
	}
	waitGroup.Wait()
	_postBody, _ := json.Marshal(map[string][]string{"allNetworkNodes": append(b.NetworkNodes, b.CurrentNodeUrl)})
	fmt.Println(string(_postBody))
	fmt.Println("network nodes " + strings.Join(b.NetworkNodes, ", "))
	_responseBody := bytes.NewBuffer(_postBody)
	_resp, _err := http.Post(newNodeUrl+"/register-nodes-bulk", "application/json", _responseBody)
	fmt.Println("broadcast report send")
	if _err == nil {
		defer _resp.Body.Close()
		_reqBody := getBodyAsBytes(_resp.Body)
		fmt.Println("register-nodes-bulk response: " + string(_reqBody))
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"message\":\"Bulk registration successfully performed!\"}"))
}

func RegisterNewNode(w http.ResponseWriter, r *http.Request, b *blockchain.Blockchain) {
	request := make(map[string]string)
	reqBody := getBodyAsBytes(r.Body)
	_ = json.Unmarshal(reqBody, &request)
	newNodeUrl := request["newNodeUrl"]
	addNewNodeIfNotExists(newNodeUrl, b)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"message\":\"Node: " + newNodeUrl + " successfully added!\"}"))
}

func RegisterNewNodesBulk(w http.ResponseWriter, r *http.Request, b *blockchain.Blockchain) {
	request := make(map[string][]string)
	reqBody := getBodyAsBytes(r.Body)
	_ = json.Unmarshal(reqBody, &request)
	allNetworkNodes := request["allNetworkNodes"]
	for _, newNodeUrl := range allNetworkNodes {
		addNewNodeIfNotExists(newNodeUrl, b)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("{\"message\":\"Nodes: " + strings.Join(allNetworkNodes, ", ") + " successfully added!\"}"))
}

func addNewNodeIfNotExists(newNodeUrl string, blockchain *blockchain.Blockchain) {
	isNotCurrentNodeUrl := sanitizedString(newNodeUrl) != blockchain.CurrentNodeUrl
	newNodeUrlNotPresent := !stringInSlice(newNodeUrl, blockchain.NetworkNodes)
	if isNotCurrentNodeUrl && newNodeUrlNotPresent {
		blockchain.NetworkNodes = append(blockchain.NetworkNodes, newNodeUrl)
	}
}

func sanitizedString(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
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
		if sanitizedString(b) == sanitizedString(a) {
			return true
		}
	}
	return false
}
