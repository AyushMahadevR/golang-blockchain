package main

import (
	"go-blockchain/blockchain"
	"go-blockchain/controller"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var DEFAULT_PORT = "80"
var DEFAULT_NODE = "http://localhost"

func main() {
	nodeAddress := uuid.New().String()
	port, currentNodeUrl := getPortAndNodeUrl(1, 2)
	bitcoin := blockchain.InitBlockChain(currentNodeUrl)
	http.HandleFunc("/blockchain", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("/blockchain called", r.Body)
		controller.BlockChainController(rw, r, bitcoin)
	})
	http.HandleFunc("/transaction", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("/transaction called", r.Body)
		controller.TranscationController(rw, r, bitcoin)
	})
	http.HandleFunc("/transaction/broadcast", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("/transaction/broadcast called", r.Body)
		controller.TranscationBroadcastController(rw, r, bitcoin)
	})
	http.HandleFunc("/mine", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("/mine called", r.Body)
		controller.Mine(rw, r, bitcoin, nodeAddress)
	})
	http.HandleFunc("/register-and-broadcast-node", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("/register-and-broadcast-node", r.Body)
		controller.RegisterAndBroadcastNewNode(rw, r, bitcoin)
	})
	http.HandleFunc("/register-node", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("/register-node called", r.Body)
		controller.RegisterNewNode(rw, r, bitcoin)
	})
	http.HandleFunc("/register-nodes-bulk", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("/register-node called", r.Body)
		controller.RegisterNewNodesBulk(rw, r, bitcoin)
	})

	log.Println("Go! server running on port: " + port)

	http.ListenAndServe(":"+port, nil)

}

func getPortAndNodeUrl(portArgPosition int, NodeUrlArgPosition int) (string, string) {
	port := getArgument(DEFAULT_PORT, portArgPosition)
	currentNodeUrl := getArgument(DEFAULT_NODE, NodeUrlArgPosition) + ":" + port
	return port, currentNodeUrl
}

func getArgument(defaultValue string, position int) string {
	if (position != 0) && len(os.Args) > position {
		return os.Args[position]
	}
	return defaultValue
}
