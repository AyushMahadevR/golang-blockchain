package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"time"
)

type Transaction struct {
	Amount   float32 `json:"amount"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
}

type Block struct {
	Index        uint32         `json:"index"`
	Nonce        uint32         `json:"nonce"`
	Transcations []*Transaction `json:"transcations"`
	PrevHash     string         `json:"prevHash"`
	Hash         string         `json:"hash"`
	Timestamp    uint32         `json:"timestamp"`
}

type Blockchain struct {
	Blocks              []*Block       `json:"blocks"`
	PendingTransactions []*Transaction `json:"pendingTranscations"`
}

func InitBlockChain() *Blockchain {
	genesis := CreateGenesisBlock()
	return &Blockchain{Blocks: []*Block{genesis}}
}

func CreateGenesisBlock() *Block {
	prevHash := sha256.Sum256([]byte("Golang Genesis Previous"))
	hash := sha256.Sum256([]byte("Golang Genesis Hash"))
	return &Block{
		Index:        0,
		Transcations: []*Transaction{},
		PrevHash:     fmt.Sprintf("%x", prevHash[:]),
		Hash:         fmt.Sprintf("%x", hash[:]),
		Timestamp:    GetTimestamp()}
}

func GetTimestamp() uint32 {
	return uint32(time.Now().Unix())
}

func (b *Blockchain) Mine() {
	noOfBlocks := len(b.Blocks)
	previousBlock := b.Blocks[noOfBlocks-1]
	previousHash := []byte(previousBlock.Hash)
	currentBlock := &Block{Transcations: b.PendingTransactions}
	nonce, currentBlockHash := ProofOfWork(previousHash, currentBlock)
	b.CreateNewBlock(noOfBlocks, nonce, string(previousHash), currentBlockHash)
}

func (b *Blockchain) CreateNewBlock(index int, nonce int32, previousHash string, blockHash string) {
	Block := &Block{
		Index:        uint32(index),
		Nonce:        uint32(nonce),
		Transcations: b.PendingTransactions,
		PrevHash:     previousHash,
		Hash:         blockHash,
		Timestamp:    GetTimestamp(),
	}
	b.PendingTransactions = []*Transaction{CreateTransaction("00", "Miner#id", 12.5)}
	b.Blocks = append(b.Blocks, Block)
}

func (b *Blockchain) AddToPendingTransaction(newTransaction *Transaction) {
	b.PendingTransactions = append(b.PendingTransactions, newTransaction)
}

func ProofOfWork(previousHash []byte, currentBlock *Block) (int32, string) {
	nonce := 0
	hash, memory := GenerateBlockHash(nonce, previousHash, currentBlock, []byte{})
	for hash[:5] != "00000" {
		nonce++
		hash, _ = GenerateBlockHash(nonce, previousHash, currentBlock, memory)
	}
	// fmt.Println(stringHash)
	// fmt.Println(stringHash[:5] != "00000")
	return int32(nonce), hash
}

func GenerateBlockHash(nonce int, previousHash []byte, currentBlock *Block, memory []byte) (string, []byte) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(nonce))
	byteBlock := []byte{}
	if len(memory) > 0 {
		byteBlock = memory
	} else {
		byteBlock = GetByteTransactions(currentBlock.Transcations)
	}
	combinedByteData := append(append(b, previousHash...), byteBlock...)
	// fmt.Printf("Nonce value: %v\n", nonce)
	return fmt.Sprintf("%x", sha256.Sum256(combinedByteData)), byteBlock
}

func GetByteTransactions(data []*Transaction) []byte {
	allTransaction := []byte{}
	for _, transaction := range data {
		allTransaction = append(allTransaction, transaction.GetByteTransaction()...)
	}
	return allTransaction
}

func (t *Transaction) GetByteTransaction() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	_ = encoder.Encode(t.Amount)
	_ = encoder.Encode(t.Receiver)
	_ = encoder.Encode(t.Sender)
	return w.Bytes()
}

func CreateNewBlock() *Block {
	return &Block{}
}

func (b *Block) AddTranscationToBlock(Transaction *Transaction) {
	b.Transcations = append(b.Transcations, Transaction)
}

func CreateTransaction(sender string, receiver string, amount float32) *Transaction {
	return &Transaction{Amount: amount, Sender: sender, Receiver: receiver}
}
