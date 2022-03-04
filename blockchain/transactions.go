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
	Owner string	`json:"owner"`
	Amount int	`json:"amount"`
}

type TxOut struct {
	Owner string	`json:"owner"`
	Amount int	`json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
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


func makeTx(from, to string, amount int) (*Tx, error) {
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("Not Enough Money.")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, &txIn)
		total += txIn.Amount
	}
	change := total - amount
	if change > 0 {
		txOut := TxOut{from, change}
		txOuts = append(txOuts, &txOut)
	}
	txOut := TxOut{to, amount}
	txOuts = append(txOuts, &txOut)
	tx := Tx{
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxIns: txIns,
		TxOuts: txOuts,
	}
	tx.getId()
	return &tx, nil
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