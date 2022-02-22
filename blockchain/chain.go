package blockchain

import (
	"fmt"
	"sync"

	"github.com/styu12/seungohcoin/db"
	"github.com/styu12/seungohcoin/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height int `json:"height"`
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

// b를 대신 내보내는 함수
func Blockchain() *blockchain {
	if b == nil {
		// 여러 개의 goRoutine이 동시에 첫 블록체인 생성을 요구할 수도 있으니 더욱 확실하게 한번만 실행!
		once.Do(func() {
			b = &blockchain{"", 0}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				fmt.Println("checkpoint is nil")
				b.AddBlock("Genesis")
			}	else {
				b.restore(checkpoint)
			}
		})
	}
	fmt.Printf("Newest Hash : %s\n", b.NewestHash)
	fmt.Printf("Height : %d\n", b.Height)
	return b
} 