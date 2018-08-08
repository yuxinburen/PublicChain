package BLC

import "fmt"

//打印输出系统中的钱包列表
func (cli *CLI) PrintWalletList() {
	fmt.Printf("所有钱包地址:\n")
	//获取钱包的集合，遍历，依次输出
	wallets := NewWallets()
	for address, _ := range wallets.WalletMap {
		fmt.Println(address)
	}
}
