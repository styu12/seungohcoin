package blockchain

import (
	"fmt"
	"sync"

	"github.com/styu12/seungohcoin/db"
	"github.com/styu12/seungohcoin/utils"
)

const (
	defaultDifficulty int = 2
	difficultyInterval int = 5
	blockInterval int = 2
	allowedRange int = 2
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height int `json:"height"`
	CurrentDifficulty int `json:"currentDifficulty"`
}

func (b *blockchain) txOuts() []*TxOut {
	var txOuts []*TxOut
	blocks := b.Blocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

func (b *blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}
	return ownedTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	total := 0
	txOuts := b.TxOutsByAddress(address)
	for _, txOut := range txOuts {
		total += txOut.Amount
	}
	return total
}

func (b *blockchain) recalculateDifficulty() {
	blocks := b.Blocks()
	newestBlock := blocks[0]
	pointerBlock := blocks[difficultyInterval - 1]
	actualTime := (newestBlock.Timestamp/60) - (pointerBlock.Timestamp/60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime < (expectedTime - allowedRange) {
		b.CurrentDifficulty++
	}	else if actualTime > (expectedTime - allowedRange) {
		b.CurrentDifficulty--
	}
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		b.CurrentDifficulty = defaultDifficulty
	}	else if b.Height % difficultyInterval == 0 {
		b.recalculateDifficulty()
	}
	return b.CurrentDifficulty
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		}	else {
			break
		}
	}
	return blocks
}

// b를 대신 내보내는 함수
func Blockchain() *blockchain {
	if b == nil {
		// 여러 개의 goRoutine이 동시에 첫 블록체인 생성을 요구할 수도 있으니 더욱 확실하게 한번만 실행!
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				fmt.Println("checkpoint is nil")
				b.AddBlock()
			}	else {
				b.restore(checkpoint)
			}
		})
	}
	return b
} 