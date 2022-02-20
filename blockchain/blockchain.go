package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data string	`json:"data"`
	Hash string `json:"hash"`
	PrevHash string	`json:"prevHash,omitempty"`
	Height int 	`json:"height"`
}

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

type blockchain struct {
	// blocks slice가 너무 커질 수 있으니 pointer만 저장하자!
	blocks []*Block
}

var b *blockchain
var once sync.Once

func getLastHash() string {
	totalBlocks := GetBlockchain().blocks
	if len(totalBlocks) == 0 {
		return ""
	}
	return totalBlocks[len(totalBlocks) - 1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().blocks) + 1}
	newBlock.calculateHash()
	// blockchain = pointer slice이므로 new Block 내보낼 때도 Block pointer를 내보내야지
	return &newBlock
}

// b를 대신 내보내는 함수
func GetBlockchain() *blockchain {
	if b == nil {
		// 여러 개의 goRoutine이 동시에 첫 블록체인 생성을 요구할 수도 있으니 더욱 확실하게 한번만 실행!
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
} 

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
} 

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}

func (b *blockchain) GetBlock(height int) *Block {
	return b.blocks[height - 1]
}