package BLC

//创建区块链
func (cli *CLI) CreateBlockChain(address string) {
	//创建创世区块
	CreateBlockChainWithGenesisBlock(address)
}
