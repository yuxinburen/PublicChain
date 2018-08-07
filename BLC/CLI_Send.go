package BLC

import (
	"fmt"
	"os"
)

//执行转账交易的业务层方法
func (cli *CLI) Send(fromArgs []string, toArgs []string, amountArgs []string) {

	//思路:
	//1.先拿到区块链对象
	//2.如果区块链对象为nil,说明没有区块链，提示后直接结束运行
	//3.如果区块蓝对象不为空,则执行转账交易
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Printf("BlockChain未创建，无法实现转账\n")
		os.Exit(0)
	}
	defer bc.DB.Close()
	bc.Send(fromArgs, toArgs, amountArgs)
}
