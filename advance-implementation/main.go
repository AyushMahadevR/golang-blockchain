package main

import (
	"encoding/json"
	"go-blockchain/blockchain"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/blockchain", func(w http.ResponseWriter, r *http.Request) {
		bitcoin := blockchain.InitBlockChain()
		transaction := blockchain.CreateTransaction("Ayush", "Tesla", 5000.25)

		bitcoin.AddToPendingTransaction(transaction)
		bitcoin.Mine()
		bitcoinJSON, _ := json.Marshal(bitcoin)
		w.Header().Add("content-type", "text/json")
		w.Write(bitcoinJSON)
	})

	log.Println("Go!")

	http.ListenAndServe(":80", nil)

}
