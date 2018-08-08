package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

//转账交易中的交易对象
type Transaction struct {
	TxID []byte //交易的ID
	//交易的输入方
	Vins  []*TxInput
	Vouts []*TxOutput //输出方，找零
}

//根据tx,生成一个hash
func (transaction *Transaction) SetID() {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(transaction)
	if err != nil {
		log.Panic(err)
	}
	//byte[] --> hash
	hash := sha256.Sum256(buf.Bytes())
	transaction.TxID = hash[:]
}

//创建一个单笔交易
func NewSimpleTransaction(from, to string, amount int64, chain *BlockChain, txs []*Transaction) *Transaction {

	//1. 定义Input和Output的数据
	var txInputs []*TxInput
	var txOutputs []*TxOutput

	//创建Input
	/*
		创世区块中交易ID：c16d3ad93450cd532dcd7ef53d8f396e46b2e59aa853ad44c284314c7b9db1b4
 	*/
	/**
 	*TxID应该是自动寻找创建的
	 */
	//首先获取此次转账需要用到的output
	total, spentableUTXO := chain.FindSpentableUTXOs(from, to, amount, txs) //map[TxID] --> []int{index}的形式
	for txID, indexArray := range spentableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TxInput{txIDBytes, index, from}
			txInputs = append(txInputs, txInput)
		}
	}

	//idBytes, _ := hex.DecodeString("2a2303271c75c8b9bae50d73404cf36f15b3ebb0abee9a8cc4132df57c901c1f")
	//txInput := &TxInput{idBytes, 1, from}
	//txInputs = append(txInputs, txInput)

	//3.创建out
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput2 := &TxOutput{total - amount, from}
	txOutputs = append(txOutputs, txOutput2)

	//创建交易
	tx := &Transaction{[]byte{}, txInputs, txOutputs}

	//设置交易的ID
	tx.SetID()
	return tx
}

//判断tx是否是CoinBase交易
func (tx *Transaction) IsCoinBaseTransaction() bool {
	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}
