package main

import "fmt"

//多个钱包
type Wallets struct {
	WalletMap map[string]*Wallet //钱包集合
}

//生成一个钱包集合对象
func NewWallets() *Wallets {
	wallets := &Wallets{}
	wallets.WalletMap = make(map[string]*Wallet)
	return wallets
}

//创建一个新的钱包
func (ws *Wallets) CreateNewWallet() {
	wallet := NewWallet()

	address := wallet.GetAddress()
	fmt.Printf("创建的钱包的地址:%s\n", address)
	ws.WalletMap[string(address)] = wallet

}
