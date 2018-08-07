package BLC

import (
	"fmt"
	"os"
)

//注：在添加了send转账功能以后，该方法就被废弃掉，不再使用，在功能中也不再体现了

//添加新的区块数据到区块链中
func (cli *CLI) AddBlockToBlockChain(txs []*Transaction) {
	blockChain := GetBlockChainObject()
	if blockChain == nil {
		fmt.Printf("没有数据库,无法添加新的区块\t")
		os.Exit(1)
	}
	defer blockChain.DB.Close()
	blockChain.AddBlockToBlockChain(txs)
}
