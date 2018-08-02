package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代器数据结构
type BlockChainIterator struct {
	DB          *bolt.DB
	CurrentHash []byte
}

//迭代器的下一个
func (iterator *BlockChainIterator) Next() *Block {
	block := new(Block)
	//根据iterator,操作db对象，读取数据库
	err := iterator.DB.View(func(tx *bolt.Tx) error {
		//自己根据已经赋值的最新的hash获取block，并更新tip的hash串
		bk := tx.Bucket([]byte(BlockBucketName))
		if bk != nil {
			//1.根据当前的hash值获取block串,并进行反序列化
			//2.更新保存hash的变量的值为最新查出来的hash值
			block = Deserialize(bk.Get(iterator.CurrentHash))
			iterator.CurrentHash = block.PreBlockHash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	return block
}
