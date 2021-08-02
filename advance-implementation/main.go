package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Transaction struct {
	amount   float32
	sender   []byte
	receiver []byte
}

type Block struct {
	index       uint32
	nonce       uint32
	transcation []*Transaction
	prevHash    []byte
	hash        []byte
	timestamp   uint32
}

type Blockchain struct {
	blocks              []*Block
	pendingTransactions []*Transaction
}

func main() {
	bitcoin := InitBlockChain()
	transaction := createTransaction("Ayush", "Tesla", 5000.25)
	bitcoin.addToPendingTransaction(transaction)

	http.HandleFunc("/blockchain", func(w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(*transaction)
		fmt.Println(err)
		w.Header().Add("content-type", "text/json")
		w.Write(j)
		fmt.Println(string(j))
	})

	log.Println("Go!")

	http.ListenAndServe(":8080", nil)

	// bitcoin.mine()
	// fmt.Println(bitcoin.blocks[1])
	// fmt.Printf("%+v", bitcoin.pendingTransactions[0])
	// fmt.Println(proofOfWork([]byte("test"), bitcoin.blocks[0]))
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&Transaction{
		amount:   t.amount,
		sender:   t.sender,
		receiver: t.receiver,
	})
}

func InitBlockChain() *Blockchain {
	genesis := createGenesisBlock()
	return &Blockchain{blocks: []*Block{genesis}}
}

func createGenesisBlock() *Block {
	prevHash := sha256.Sum256([]byte("Golang Genesis Previous"))
	hash := sha256.Sum256([]byte("Golang Genesis Hash"))
	return &Block{index: 0, prevHash: prevHash[:], hash: hash[:], timestamp: getTimestamp()}
}

func getTimestamp() uint32 {
	return uint32(time.Now().Unix())
}

func (b *Blockchain) mine() {
	noOfBlocks := len(b.blocks)
	previousBlock := b.blocks[noOfBlocks-1]
	previousHash := previousBlock.hash
	currentBlock := &Block{transcation: b.pendingTransactions}
	nonce, currentBlockHash := proofOfWork(previousHash, currentBlock)
	b.createNewBlock(noOfBlocks, nonce, string(previousHash), currentBlockHash)
}

func (b *Blockchain) createNewBlock(index int, nonce int32, previousHash string, blockHash string) {
	block := &Block{
		index:       uint32(index),
		nonce:       uint32(nonce),
		transcation: b.pendingTransactions,
		prevHash:    []byte(previousHash),
		hash:        []byte(blockHash),
		timestamp:   getTimestamp(),
	}
	b.pendingTransactions = []*Transaction{createTransaction("00", "Miner#id", 12.5)}
	b.blocks = append(b.blocks, block)
}

func (b *Blockchain) addToPendingTransaction(newTransaction *Transaction) {
	b.pendingTransactions = append(b.pendingTransactions, newTransaction)
}

func proofOfWork(previousHash []byte, currentBlock *Block) (int32, string) {
	nonce := 0
	hash, memory := generateBlockHash(nonce, previousHash, currentBlock, []byte{})
	for hash[:5] != "00000" {
		nonce++
		hash, _ = generateBlockHash(nonce, previousHash, currentBlock, memory)
	}
	// fmt.Println(stringHash)
	// fmt.Println(stringHash[:5] != "00000")
	return int32(nonce), hash
}

func generateBlockHash(nonce int, previousHash []byte, currentBlock *Block, memory []byte) (string, []byte) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(nonce))
	byteBlock := []byte{}
	if len(memory) > 0 {
		byteBlock = memory
	} else {
		byteBlock = getByteTransactions(currentBlock.transcation)
	}
	combinedByteData := append(append(b, previousHash...), byteBlock...)
	// fmt.Printf("Nonce value: %v\n", nonce)
	return fmt.Sprintf("%x", sha256.Sum256(combinedByteData)), byteBlock
}

func getByteTransactions(data []*Transaction) []byte {
	allTransaction := []byte{}
	for _, transaction := range data {
		allTransaction = append(allTransaction, transaction.getByteTransaction()...)
	}
	return allTransaction
}

func (t *Transaction) getByteTransaction() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	_ = encoder.Encode(t.amount)
	_ = encoder.Encode(t.receiver)
	_ = encoder.Encode(t.sender)
	return w.Bytes()
}

func createNewBlock() *Block {
	return &Block{}
}

func (b *Block) addTranscationToBlock(Transaction *Transaction) {
	b.transcation = append(b.transcation, Transaction)
}

func createTransaction(sender string, receiver string, amount float32) *Transaction {
	return &Transaction{amount: amount, sender: []byte(sender), receiver: []byte(receiver)}
}
