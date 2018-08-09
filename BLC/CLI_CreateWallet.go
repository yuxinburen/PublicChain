package BLC

//创建一个新的钱包并返回打印钱包地址
func (cli *CLI) CreateWallet() {
	wallets := NewWallets()
	wallets.CreateNewWallet()
}
