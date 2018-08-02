package BLC

import (
	"github.com/boltdb/bolt"
	"fmt"
	"log"
	"os"
	"time"
	"math/big"
)

type BlockChain struct {
	DB  *bolt.DB //对应的数据库对象
	Tip []byte   //存储区块中最后一个块的hash值
}

//向区块链中保存数据
func (chain *BlockChain) AddBlockToBlockChain(data string) {
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
			newBlock := NewBlock(data, lastBlock.Hash, lastBlockHeight+1)
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
		fmt.Printf("\t区块的信息:%s\n", block.Data)
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

//create a blockchain,nclude genesis block
func CreateBlockChainWithGenesisBlock(data string) *BlockChain {
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
	fmt.Printf("数据库不存在...")
	//1.创建数据库
	//2.创建创世区块
	//3.存储到数据库中

	genesisBlock := CreateGenesisBlock(data)

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
