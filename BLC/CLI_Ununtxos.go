package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) GetUnuntxos() {
	bc := GetBlockChainObject()
	if bc == nil {
		fmt.Println("没有数据库，无法后去utxos..")
		os.Exit(1)
	}
	defer bc.DB.Close()

	unSpentUTXOsMap := bc.FindUnspentUTXOMap()

	for txIDStr, txoutputs := range unSpentUTXOsMap {
		fmt.Println("交易ID：", txIDStr)
		for _, utxo := range txoutputs.UTXOs {
			fmt.Println("\t金额:", utxo.Output.Value)
			fmt.Println("\t地址:", GetAddressByPubKeyHash(utxo.Output.PubKeyHash))
			fmt.Println("----------------------------------------------")
		}
	}

	utxoSet := &UTXOSet{bc}
	utxoSet.ResetUTXOSet()
}
