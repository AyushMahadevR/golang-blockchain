package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

type Transaction struct {
	amount   float32
	sender   []byte
	receiver []byte
}

type Block struct {
	data     []*Transaction
	prevHash []byte
	hash     []byte
}

type Blockchain struct {
	blocks              []*Block
	pendingTransactions []*Transaction
}

func main() {
	bitcoin := InitBlockChain()
	transaction := createTransaction("Ayush", "Tesla", 5000.25)
	bitcoin.addToPendingTransaction(transaction)
	fmt.Println(proofOfWork([]byte("test"), bitcoin.blocks[0]))
}

func InitBlockChain() *Blockchain {
	genesis := createGenesisBlock()
	return &Blockchain{blocks: []*Block{genesis}}
}

func createGenesisBlock() *Block {
	prevHash := sha256.Sum256([]byte("Golang Genesis Previous"))
	hash := sha256.Sum256([]byte("Golang Genesis Hash"))
	data := &Transaction{amount: 51000.85, sender: []byte("Dexter associates"), receiver: []byte("AA Markers")}
	return &Block{prevHash: prevHash[:], data: []*Transaction{data}, hash: hash[:]}
}

func (b *Blockchain) addToPendingTransaction(newTransaction *Transaction) {
	b.pendingTransactions = append(b.pendingTransactions, newTransaction)
}

func proofOfWork(previousHash []byte, currentBlock *Block) int32 {
	nonce := 0
	hash, memory := generateBlockHash(nonce, previousHash, currentBlock, []byte{})
	for hash[:5] != "00000" {
		nonce++
		hash, _ = generateBlockHash(nonce, previousHash, currentBlock, memory)
	}
	// fmt.Println(stringHash)
	// fmt.Println(stringHash[:5] != "00000")
	return int32(nonce)
}

func generateBlockHash(nonce int, previousHash []byte, currentBlock *Block, memory []byte) (string, []byte) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(nonce))
	byteBlock := []byte{}
	if len(memory) > 0 {
		byteBlock = memory
	} else {
		byteBlock = getByteTransactions(currentBlock.data)
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
	b.data = append(b.data, Transaction)
}

func createTransaction(sender string, receiver string, amount float32) *Transaction {
	return &Transaction{amount: amount, sender: []byte(sender), receiver: []byte(receiver)}
}
