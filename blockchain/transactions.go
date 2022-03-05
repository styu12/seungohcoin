package blockchain

import (
	"errors"
	"time"

	"github.com/styu12/seungohcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id string	`json:"id"`
	Timestamp int	`json:"timestamp"`
	TxIns []*TxIn	`json:"txIns"`
	TxOuts []*TxOut	`json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	TxId string 	`json:"txId"`
	Index int 	`json:"index"`
	Owner string	`json:"owner"`
}

type TxOut struct {
	Owner string	`json:"owner"`
	Amount int	`json:"amount"`
}

type UTxOut struct {
	TxId string		`json:"txId"`
	Index int 		`json:"index"`
	Amount int  	`json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", 0, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxIns: txIns,
		TxOuts: txOuts,
	}
	tx.getId()
	return &tx
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
	Outer:
	for _, Tx := range Mempool.Txs {
		for _, input := range Tx.TxIns {
			if input.TxId == uTxOut.TxId && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}


func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("Not Enough Money.")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxId, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}	
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxIns: txIns,
		TxOuts: txOuts,
	}
	tx.getId()
	return tx, nil

}

func (m *mempool) AddTx(to string, amount int) error {
	newTx, err := makeTx("Seungoh", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, newTx)
	return err
}

func (m *mempool) TxToConfirm() []*Tx {
	txs := m.Txs
	coinbase := makeCoinbaseTx("Seungoh")
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}