package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/styu12/seungohcoin/db"
	"github.com/styu12/seungohcoin/utils"
)

type Block struct {
	Data string	`json:"data"`
	Hash string `json:"hash"`
	PrevHash string	`json:"prevHash,omitempty"`
	Height int 	`json:"height"`
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))	
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data: data,
		Hash: "",
		PrevHash: prevHash,
		Height: height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.Hash = hash
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