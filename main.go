package main

import (
	"github.com/styu12/seungohcoin/wallet"
)

func main() {
	wallet.Wallet()
}

// func main() {
// 	defer db.CloseDB()
// 	cli.Start()
	
// }

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