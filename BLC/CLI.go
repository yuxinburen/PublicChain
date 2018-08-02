package BLC

import (
	"os"
	"fmt"
	"flag"
	"log"
)

type CLI struct {
}

//Cli的运行
func (cli *CLI) Run() {

	//判断命令行参数，不符合规则 输出随用规则及用法
	isValidArgs()

	//处理参数及相关的命令对应的业务逻辑
	//1.创建flagset命令对象
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//2.设置命令后的参数对象
	flagCreateBlockChainData := createBlockChainCmd.String("data", "GensisisBlock", "创世区块的信息")
	flagAddBlockData := addBlockCmd.String("data", "helloworld", "区块的数据")

	var err error
	//解析用户的意图命令
	switch os.Args[1] {
	case "createblockchain": //创建创世区块及区块链
		err = createBlockChainCmd.Parse(os.Args[2:])
	case "addblock": //向区块链中增加新的区块
		err = addBlockCmd.Parse(os.Args[2:])
	case "printchain": //将区块链数据从数据库中查询出来并打印
		err = printChainCmd.Parse(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	//根据用户在中断输入的命令执行对应的功能
	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockChainData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.CreateBlockChain(*flagCreateBlockChainData)
	}

	//根据用户在终端输入的命令添加区块到区块链中的功能
	if addBlockCmd.Parsed() {
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		cli.AddBlockToBlockChain(*flagAddBlockData)
	}

	//根据用户在终端输入的命令：打印出所有区块蓝数据
	if printChainCmd.Parsed() {
		cli.PrintChains()
	}

}

//检查用户参数是否合法
func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}

//打印程序用法说明
func printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("\tcreateblockchain -data DATA --创建创世区块\n")
	fmt.Printf("\taddblock -data DATA --添加区块\n")
	fmt.Printf("\tprintchain --打印所有区块\n")
}

//创建区块链
func (cli *CLI) CreateBlockChain(data string) {
	//创建创世区块
	CreateBlockChainWithGenesisBlock(data)
}

//添加新的区块数据到区块链中
func (cli *CLI) AddBlockToBlockChain(data string) {
	blockChain := GetBlockChainObject()
	if blockChain == nil {
		fmt.Printf("没有数据库,无法添加新的区块\t")
		os.Exit(1)
	}
	blockChain.AddBlockToBlockChain(data)
}

//打印出所有区块的信息
func (cli *CLI) PrintChains() {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Printf("没有blockchain，无法打印任何区块数据")
		os.Exit(1)
	}
	//调用bc的打印数据的方法
	bc.PrintChains()
}
