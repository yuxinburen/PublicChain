package BLC

import (
	"encoding/hex"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
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
	if err == nil {
		log.Panic(err)
	}
	//byte[] --> hash
	hash := sha256.Sum256(buf.Bytes())
	transaction.TxID = hash[:]
}

//创建一个单笔交易
func NewSimpleTransaction(from, to string, amount int64) *Transaction {

	//1. 定义Input和Output的数据
	var txInputs []*TxInput
	var txOutputs []*TxOutput

	//创建Input
	/*
		创世区块中交易ID：c16d3ad93450cd532dcd7ef53d8f396e46b2e59aa853ad44c284314c7b9db1b4
 	*/
	idBytes, _ := hex.DecodeString("c16d3ad93450cd532dcd7ef53d8f396e46b2e59aa853ad44c284314c7b9db1b4")
	txInput := &TxInput{idBytes, 1, from}
	txInputs = append(txInputs, txInput)

	//创建out
	txOutput := &TxOutput{amount, to}
	txOutputs = append(txOutputs, txOutput)

	//找零
	txOutput2 := &TxOutput{6 - amount, from}
	txOutputs = append(txOutputs, txOutput2)

	//创建交易
	tx := &Transaction{[]byte{}, txInputs, txOutputs}

	//设置交易的ID
	tx.SetID()
	return tx
}
