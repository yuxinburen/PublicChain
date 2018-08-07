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
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	printWalletListCmd := flag.NewFlagSet("printwalletlist", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	//2.设置命令后的参数对象
	flagCreateBlockChainData := createBlockChainCmd.String("data", "GensisisBlock", "创世区块的信息")
	//flagAddBlockData := addBlockCmd.String("data", "helloworld", "区块的数据")
	//send 的参数对象
	flagSendFromData := sendBlockCmd.String("from", "", "转账源地址")
	flagSendToData := sendBlockCmd.String("to", "", "入账地址")
	flagSendAmountData := sendBlockCmd.String("amount", "", "转账金额")
	//余额功能参数接收
	flagGetBalanceData := getBalanceCmd.String("address", "", "账户地址")

	var err error
	//解析用户的意图命令
	switch os.Args[1] {
	case "createwallet": //创建钱包
		err = createWalletCmd.Parse(os.Args[2:])
	case "printwalletlist":
		err = printWalletListCmd.Parse(os.Args[2:])
	case "createblockchain": //创建创世区块及区块链
		err = createBlockChainCmd.Parse(os.Args[2:])
	case "send": //转账功能参数解析
		err = sendBlockCmd.Parse(os.Args[2:])
	case "printchain": //将区块链数据从数据库中查询出来并打印
		err = printChainCmd.Parse(os.Args[2:])
	case "getbalance":
		err = getBalanceCmd.Parse(os.Args[2:])
	default:
		PrintUsage()
		os.Exit(1)
	}
	if err != nil {
		log.Panic(err)
	}

	//创建钱包
	if createWalletCmd.Parsed() {
		cli.CreateWallet()
	}

	//打印系统中的钱包列表
	if printWalletListCmd.Parsed() {
		cli.PrintWalletList()
	}

	//根据用户在中断输入的命令执行对应的功能
	if createBlockChainCmd.Parsed() {
		if *flagCreateBlockChainData == "" {
			PrintUsage()
			os.Exit(1)
		}
		cli.CreateBlockChain(*flagCreateBlockChainData)
	}

	//根据用户在终端输入的命令添加区块到区块链中的功能
	//if addBlockCmd.Parsed() {
	//	if *flagAddBlockData == "" || *flagSendToData == "" || *flagSendAmountData == "" {
	//		printUsage()
	//		os.Exit(1)
	//	}
	//	//cli.AddBlockToBlockChain(*flagAddBlockData)
	//	from := JSONToArray(*flagSendFromData)
	//	to := JSONToArray(*flagSendToData)
	//	amount := JSONToArray(*flagSendAmountData)
	//	cli.Send(from, to, amount)
	//}

	//转账交易功能send
	if sendBlockCmd.Parsed() {
		if *flagSendFromData == "" || *flagSendToData == "" || *flagSendAmountData == "" {
			fmt.Printf("请指定转账参数\n")
			PrintUsage()
			os.Exit(1)
		}

		//参数正常，解析并开始转账业务
		from := JSONToArray(*flagSendFromData) //[]string 有可能设计到多个转账
		to := JSONToArray(*flagSendToData)
		amount := JSONToArray(*flagSendAmountData)
		cli.Send(from, to, amount)
	}

	//根据用户在终端输入的命令：打印出所有区块蓝数据
	if printChainCmd.Parsed() {
		cli.PrintChains()
	}

	//查询账户余额
	if getBalanceCmd.Parsed() {
		if *flagGetBalanceData == "" {
			fmt.Printf("请输入要查询的账户地址.\n")
			os.Exit(1)
		}
		cli.GetBalance(*flagGetBalanceData)
	}
}

//检查用户参数是否合法
func isValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}
}
