package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha1"
)

//创建一个block结构体
type Block struct {
	Height       int64
	PreBlockHash []byte
	Txs          []*Transaction
	TimeStamp    int64
	Hash         []byte
	//Nonce
	Nonce int64
}

//step:提供一个函数用于创建一个区块
func NewBlock(txs []*Transaction, prevBlockHash []byte, height int64) *Block {
	//创建区块
	block := &Block{Height: height,
		PreBlockHash: prevBlockHash,
		Txs: txs,
		TimeStamp: time.Now().Unix(),
	}

	pow := NewProfOfWork(block)
	hash, nonce := pow.Run() //pow算法开始工作
	block.Hash = hash
	block.Nonce = nonce

	//返回新创建的区块
	return block
}

//生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(txs, make([]byte, 32, 32), 0)
}

//区块的序列化方法：将对象数据序列化后编程数组形式
func (bl *Block) Serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(bl)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

//将一个byte形式的数据反序列化为对象数据
func Deserialize(blockBytes []byte) *Block {
	var block *Block
	reader := bytes.NewReader(blockBytes)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return block
}

//将区块内的所有的交易都hash化
func (block *Block) HashTransactions() []byte {

	//1.创建一个二维数组,存储每笔交易的txID
	var txsHashes [][]byte
	//2.遍历
	for _, tx := range block.Txs {
		txsHashes = append(txsHashes, tx.TxID)
	}
	//生成hash
	hash := sha1.Sum(bytes.Join(txsHashes, []byte{}))
	return hash[:]
}
