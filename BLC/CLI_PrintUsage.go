package BLC

import "fmt"

//打印程序用法说明
func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreatewallet --创建钱包")
	fmt.Println("\tprintwalletlist --打印输出系统的钱包列表")
	fmt.Println("\tcreateblockchain -address address --创建创世区块")
	fmt.Println("\tsend -from from -to to -amount amount --转账给他人")
	fmt.Println("\tprintchain --打印所有区块")
	fmt.Println("\tgetbalance -address address --查询账户余额")
	fmt.Println("\tgetunutxos --获取所有未花费的utxo集合")
}
