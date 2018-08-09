package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/elliptic"
)

//转账交易中的交易对象
type Transaction struct {
	TxID []byte //交易的ID
	//交易的输入方
	Vins  []*TxInput
	Vouts []*TxOutput //输出方，找零
}

//验证签名的方法
//验证签名的原理:公钥 + 要签名的数据 验证 签名：rs
func (tx *Transaction) Vertify(prevTxs map[string]*Transaction) bool {
	//1.如果coinbase交易，则不需要验签
	if tx.IsCoinBaseTransaction() {
		return true
	}

	//prevTxs
	for _, input := range prevTxs {
		if prevTxs[hex.EncodeToString(input.TxID)] == nil {
			log.Panic("当前的Input没有找到对应的Transaction，无法验证签名")
		}
	}

	//验证签名
	txCopy := tx.TrimmedCopy()
	curev := elliptic.P256() //曲线变量

	for index, input := range tx.Vins {

		//原理：再次获取 要签名的数据 ＋ 公钥哈希 ＋ 签名
		/**
		*验证签名的有效性:
		*第一个参数:公钥
		*第二个参数：签名的数据
		 *第三，四个参数；签名：R，S
		 */

	}

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

//签名:对一笔交易进行签名.
func (transaction *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTxsmap map[string]*Transaction) {
	//1.如果是coinbase交易，则不需要进行签名
	if transaction.IsCoinBaseTransaction() {
		return
	}

	//2.获取input对应的output所在的tx，如果不存在,无法进行签名
	for _, input := range transaction.Vins {
		if prevTxsmap[hex.EncodeToString(input.TxID)] == nil {
			log.Panic("当前的Input，没有找到对应的output所在的Transaction，无法签名..")
		}
	}

	//即将进行签名：私钥，要签名的数据
	txCopy := transaction.TrimmedCopy()
}

//创建一个CoinBase交易
func NewCoinBaseTransaction(address string) *Transaction {
	txInput := &TxInput{[]byte{}, -1, nil, nil}
	//txOutput := &TxOutput{10, address}
	txOutput := NewTxOutput(10, address)

	txCoinBaseTransaction := &Transaction{[]byte{}, []*TxInput{txInput}, []*TxOutput{txOutput}}
	//设置交易ID
	txCoinBaseTransaction.SetID()

	return txCoinBaseTransaction
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
	total, spentableUTXO := chain.FindSpentableUTXOs(from, amount, txs) //map[TxID] --> []int{index}的形式

	//获取钱包集合
	wallets := NewWallets()
	wallet := wallets.WalletMap[from]

	for txID, indexArray := range spentableUTXO {
		txIDBytes, _ := hex.DecodeString(txID)
		for _, index := range indexArray {
			txInput := &TxInput{txIDBytes, index, nil, wallet.Publickey}
			txInputs = append(txInputs, txInput)
		}
	}

	//idBytes, _ := hex.DecodeString("2a2303271c75c8b9bae50d73404cf36f15b3ebb0abee9a8cc4132df57c901c1f")
	//txInput := &TxInput{idBytes, 1, from}
	//txInputs = append(txInputs, txInput)

	//3.创建out
	//txOutput := &TxOutput{amount, to}
	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)

	//找零
	//txOutput2 := &TxOutput{total - amount, from}
	txOutput2 := NewTxOutput(total-amount, from)
	txOutputs = append(txOutputs, txOutput2)

	//创建交易
	tx := &Transaction{[]byte{}, txInputs, txOutputs}

	//设置交易的ID
	tx.SetID()

	//4.创建完完整的交易以后，要对交易进行签名
	chain.SignTransaction(tx, wallet.PrivateKey)

	return tx
}

//判断tx是否是CoinBase交易
func (tx *Transaction) IsCoinBaseTransaction() bool {
	return len(tx.Vins[0].TxID) == 0 && tx.Vins[0].Vout == -1
}

//交易的副本中包含的数据
//包含了原本traansaction中的输入和输出。
//输入中：sign, publickey。
func (transaction *Transaction) TrimmedCopy() *Transaction {
	var inputs [] *TxInput
	var outputs [] *TxOutput

	for _, in := range transaction.Vins {
		inputs = append(inputs, &TxInput{in.TxID, in.Vout, nil, nil})
	}

	for _, out := range transaction.Vouts {
		outputs = append(outputs, &TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := &Transaction{transaction.TxID, inputs, outputs}
	return txCopy
}
