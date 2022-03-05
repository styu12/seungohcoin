package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/styu12/seungohcoin/db"
	"github.com/styu12/seungohcoin/utils"
)

type Block struct {
	Hash string `json:"hash"`
	PrevHash string	`json:"prevHash,omitempty"`
	Height int 	`json:"height"`
	Difficulty int	`json:"difficulty"`
	Nonce int	`json:"nonce"`
	Timestamp int `json:"timestamp"`
	Transactions []*Tx	`json:"transactions"`
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		fmt.Printf("\n\n\nTarget: %s\nHash: %s\nNonce: %d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Timestamp = int(time.Now().Unix())
			b.Hash = hash
			break
		}	else {
			b.Nonce++
		}
	}
	
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))	
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash: "",
		PrevHash: prevHash,
		Height: height,
		Difficulty: difficulty(b),
		Nonce: 0,
		Timestamp: 0,
	}
	block.mine()
	block.Transactions = Mempool.TxToConfirm()
	block.persist()
	return block
}  
var ErrNotFound = errors.New("Block Not Found.")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}