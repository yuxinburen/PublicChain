package BLC

import (
	"github.com/boltdb/bolt"
	"fmt"
	"log"
	"os"
	"time"
	"math/big"
	"strconv"
)

type BlockChain struct {
	DB  *bolt.DB //对应的数据库对象
	Tip []byte   //存储区块中最后一个块的hash值
}

//打印区块数据
func (chain *BlockChain) PrintChains() {
	//打开数据库
	//循环遍历里面的数据
	//放到数据组中进行返回
	//对结果进行处理（打印或者返回给业务层等操作)
	chainIterator := chain.Iterator()
	for {
		//获取下一个
		block := chainIterator.Next()
		fmt.Printf("第%d个区块的信息:\n", block.Height+1)
		fmt.Printf("\t高度:%d\n", block.Height)
		fmt.Printf("\t上一个区块的hash:%x\n", block.PreBlockHash)
		fmt.Printf("\t当前区块自己的Hash:%x\n", block.Hash)
		fmt.Printf("\t区块的信息:\n")

		for _, tx := range block.Txs {

			fmt.Printf("\t\t交易ID：%x\n", tx.TxID) //[]byte --> 0x...
			fmt.Printf("\t\tVins:\n")
			for _, in := range tx.Vins { //每一个TxInput：TxId,vout,解锁脚本
				fmt.Printf("\t\t\tTxID：%x\n", in.TxID)
				fmt.Printf("\t\t\tVout：%x\n", in.Vout)
				fmt.Printf("\t\t\tScriptSiq：%x\n", in.ScriptSip)
			}

			fmt.Printf("\t\tVouts:\n")
			for _, out := range tx.Vouts {
				fmt.Printf("\t\t\tValue:%d\n", out.Value)
				fmt.Printf("\t\t\tScriptPubKey:%s\n", out.ScriptPubKey)
			}
		}
		fmt.Printf("\t随机数的值:%d\n", block.Nonce)
		fmt.Printf("\t区块生产时间:%s\n", time.Unix(block.TimeStamp, 0).Format("2018-08-01 20:03"))

		//step2.判断是否到了iterator的末尾，即创世区块，如果到了创世区块，则结束循环
		hashInt := new(big.Int)
		hashInt.SetBytes(block.PreBlockHash)
		if big.NewInt(0).Cmp(hashInt) == 0 {
			break
		}
	}
}

//BlockChain的迭代器方法,根据迭代器可以迭代里面的数据
func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.DB, chain.Tip}
}

//Model层的转账交易
func (chain *BlockChain) Send(fromArgs []string, toArgs []string, amountArgs []string) {
	//构建交易对象
	//将交易构建到数据区块中
	//将带有交易信息的数据保存到数据库中

	//1.创建交易对象
	var txs []*Transaction
	amountInt, _ := strconv.ParseInt(amountArgs[0], 10, 64)
	tx := NewSimpleTransaction(fromArgs[0], toArgs[0], amountInt)
	txs = append(txs, tx)

	//2.构造新的区块
	newBlock := new(Block)
	err := chain.DB.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(BlockBucketName))
		if bk != nil {
			//读取数据库
			bytes := bk.Get(chain.Tip)
			lastBlock := Deserialize(bytes)
			newBlock = NewBlock(txs, lastBlock.Hash, lastBlock.Height+1)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//将构造的新区块保存到数据库中
	err = chain.DB.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(BlockBucketName))
		if bk != nil {
			//向数据库中存入数据
			bk.Put(newBlock.Hash, newBlock.Serialize())
			//更新数据库中的表示最新的hash的值的标志数据
			bk.Put([]byte("l"), newBlock.Hash)
			chain.Tip = newBlock.Hash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

//create a blockchain,nclude genesis block
func CreateBlockChainWithGenesisBlock(address string) *BlockChain {
	//1.创建创世区块
	//2.创建区块链对象并返回

	//判断数据库是否已经存在，如果不存在，就要创建数据库
	if dbExists() {
		fmt.Printf("数据库已经存在...\n")
		//打开数据库
		db, err := bolt.Open(DBName, 0600, nil)
		if err != nil {
			log.Panic(err)
		}

		var blockchain *BlockChain
		err = db.View(func(tx *bolt.Tx) error {

			//打开bucket, 读取对应的最新的hash
			bk := tx.Bucket([]byte(BlockBucketName))
			if bk != nil {
				//读取最新的hash
				hash := bk.Get([]byte("l"))
				blockchain = &BlockChain{db, hash}
			}
			return nil
		})

		if err != nil {
			log.Panic(err)
		}
		return blockchain
	}

	//数据库不存在
	fmt.Printf("数据库不存在...\n")
	//1.创建数据库
	//2.创建创世区块
	//3.存储到数据库中

	//创建一个txs-->CoinBase
	txCoinBase := NewCoinBaseTransaction(address)

	genesisBlock := CreateGenesisBlock([]*Transaction{txCoinBase})

	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {

		bk, err := tx.CreateBucketIfNotExists([]byte(BlockBucketName))
		if err != nil {
			log.Panic(err)
		}

		if bk != nil {
			//保存创世区块的hash值和创世区块对应的数据
			err := bk.Put(genesisBlock.Hash, genesisBlock.Serialize())

			if err != nil {
				log.Panic(err)
			}
			//保存整个区块链最新的区块的hash值
			bk.Put([]byte("l"), genesisBlock.Hash)
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &BlockChain{db, genesisBlock.Hash}
}

//创建一个CoinBase交易
func NewCoinBaseTransaction(address string) *Transaction {

	txInput := &TxInput{[]byte{}, -1, "Gensis Data"}
	txOutput := &TxOutput{10, address}

	txCoinBaseTransaction := &Transaction{[]byte{}, []*TxInput{txInput}, []*TxOutput{txOutput}}
	//设置交易ID
	txCoinBaseTransaction.SetID()

	return txCoinBaseTransaction

}

//判断区块链数据是否存在
func dbExists() bool {
	//_, err := os.Stat(DBName);
	//return os.IsNotExist(err)
	if _, err := os.Stat(DBName); os.IsNotExist(err) {
		return false
	}
	return true
}

//返回一个区块链结构体对象
func GetBlockChainObject() *BlockChain {

	//思路:
	//1.如果区块链数据库存在,则进行查询最新的hash块数据并构造BlockChain对象
	//2.如果区块链数据库不存在，则表明区块链不存在，可以直接返回nil对象

	if dbExists() { //区块链存在
		//打开数据库
		db, err := bolt.Open(DBName, 0600, nil)
		if err != nil {
			log.Panic(err)
		}

		var blockChain *BlockChain
		err = db.View(func(tx *bolt.Tx) error {
			bk := tx.Bucket([]byte(BlockBucketName))
			if bk != nil {
				lastHash := bk.Get([]byte("l"))
				blockChain = &BlockChain{db, lastHash}
			}
			return nil
		})

		if err != nil {
			log.Panic(err)
		}
		return blockChain
	} else {
		fmt.Printf("数据库不存在,无法获取BlockChain对象\n")
		return nil
	}
}

//向区块链中保存数据
func (chain *BlockChain) AddBlockToBlockChain(txs []*Transaction) {
	//思路：
	//1.根据要保存的数据构建一个blockchain对象
	//2.添加到数据库中
	//3.修改数据库的最新的hash的值
	err := chain.DB.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(BlockBucketName))

		if bk != nil {
			//获取chain的tip就是最新的hash,从数据库中读取最后一个block:hash,height
			blockBytes := bk.Get(chain.Tip)
			lastBlock := Deserialize(blockBytes)
			lastBlockHeight := lastBlock.Height
			//构造一个新的区块
			newBlock := NewBlock(txs, lastBlock.Hash, lastBlockHeight+1)
			//存入到数据库中
			err := bk.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//更新数据库中的存储最新hash值的l的对应的值
			bk.Put([]byte("l"), newBlock.Hash)
			//更新BlockChain对象的最新的tips
			chain.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//根据用户输入的地址查询给定的地址账户的余额
func (chain *BlockChain) GetBalance(s string) int {
	fmt.Printf("查询账户余额功能...")
	return 0
}
