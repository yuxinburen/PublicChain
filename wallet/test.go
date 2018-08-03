package main

import "fmt"

//钱包测试类
func main() {

	wallet := NewWallet()
	fmt.Printf("钱包创建成功....\n")
	fmt.Printf("钱包私钥:\n", wallet.PrivateKey)
	fmt.Printf("钱包公钥:\n", wallet.Publickey)

	address := wallet.GetAddress()
	fmt.Printf(address)
	fmt.Printf(string(address))

}
