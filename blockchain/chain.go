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

var b *blockchain
var once sync.Once


func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	persistBlockchain(b)
}

func UTxOutsByAddress(address string,b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	// dictionary 만드는 법 {string : bool, string : bool, ...}
	creatorTxs := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxId] = true
				}
			}

			for index, output := range tx.TxOuts {
				if output.Owner == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	total := 0
	txOuts := UTxOutsByAddress(address, b)
	for _, txOut := range txOuts {
		total += txOut.Amount
	}
	return total
}

func recalculateDifficulty(b *blockchain) int {
	blocks := Blocks(b)
	newestBlock := blocks[0]
	pointerBlock := blocks[difficultyInterval - 1]
	actualTime := (newestBlock.Timestamp/60) - (pointerBlock.Timestamp/60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	}	else if actualTime > (expectedTime - allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty 
}

func difficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	}	else if b.Height % difficultyInterval == 0 {
		return recalculateDifficulty(b)
	}	else {
		return b.CurrentDifficulty
	}
}

func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func Blocks(b *blockchain) []*Block {
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