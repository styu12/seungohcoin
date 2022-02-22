package main

import (
	"github.com/styu12/seungohcoin/blockchain"
	"github.com/styu12/seungohcoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
	// blockchain.Blockchain().AddBlock("Second")
	// blockchain.Blockchain().AddBlock("Third")
	// blockchain.Blockchain().AddBlock("Fourth")
}

// import "github.com/styu12/seungohcoin/explorer"

// func main() {
// 	explorer.Start()
// }

// func main() {
// 	chain := blockchain.GetBlockchain()
// 	chain.AddBlock("Second Block")
// 	chain.AddBlock("Third Block")
// 	chain.AddBlock("Fourth Block")
// 	for _, block := range chain.AllBlocks() {
// 		fmt.Printf("Data : %s\n", block.Data)
// 		fmt.Printf("Hash : %s\n", block.Hash)
// 		fmt.Printf("Prev Hash : %s\n\n", block.PrevHash)
// 	}
// }