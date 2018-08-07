package BLC

import (
	"fmt"
	"os"
)

//打印出所有区块的信息
func (cli *CLI) PrintChains() {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Printf("没有blockchain，无法打印任何区块数据\n")
		os.Exit(1)
	}
	defer bc.DB.Close()
	//调用bc的打印数据的方法
	bc.PrintChains()
}
