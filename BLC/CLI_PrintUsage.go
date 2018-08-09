package BLC

import "fmt"

//打印程序用法说明
func PrintUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("\tcreatewallet --创建钱包\n")
	fmt.Printf("\tprintwalletlist --打印输出系统的钱包列表\n")
	fmt.Printf("\tcreateblockchain -address address --创建创世区块\n")
	fmt.Printf("\tsend -from from -to to -amount amount --转账给他人\n")
	fmt.Printf("\tprintchain --打印所有区块\n")
	fmt.Printf("\tgetbalance -address address --查询账户余额\n")
}
