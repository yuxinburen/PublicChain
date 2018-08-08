package BLC

import (
	"fmt"
	"os"
)

//查询账户地址的余额数据
func (cli *CLI) GetBalance(address string) {
	blockChain := GetBlockChainObject()
	if blockChain == nil {
		fmt.Printf("没有数据库，无法查询账户余额.请先创建区块链数据库.\n")
		os.Exit(1)
	}
	defer blockChain.DB.Close()
	totalBalance := blockChain.GetBalance(address, []*Transaction{})
	fmt.Printf("账户%s的余额为:%d\n", address, totalBalance)
}
