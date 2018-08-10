package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
)

type UTXOSet struct {
	BlockChain *BlockChain
}

//查询block块中所有的未花费的txo：执行FindUnspentUTXOMap --> map
func (utxoset *UTXOSet) ResetUTXOSet() {
	err := utxoset.BlockChain.DB.Update(func(tx *bolt.Tx) error {
		//1.utxoset表存在,删除
		b := tx.Bucket([]byte(UtxoSettable))
		if b != nil {
			err := tx.DeleteBucket([]byte(UtxoSettable))
			if err != nil {
				log.Panic("重置时,删除表失败..")
			}
		}

		//2.创建utxoset
		b, err := tx.CreateBucket([]byte(UtxoSettable))
		if err != nil {
			log.Panic("重置时，创建表失败...")
		}
		if b != nil {

			//3.将map数据-->表
			unUTXOMap := utxoset.BlockChain.FindUnspentUTXOMap()
			for txIDstr, outs := range unUTXOMap {
				txID, _ := hex.DecodeString(txIDstr)
				b.Put(txID, outs.Serialize())
			}
			fmt.Printf("保存所有的未输出的交易成功")
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}
