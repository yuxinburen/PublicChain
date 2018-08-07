package main

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"io/ioutil"
	"os"
	"crypto/elliptic"
)

//多个钱包
type Wallets struct {
	WalletMap map[string]*Wallet //钱包集合
}

//生成一个钱包集合对象
//思路:
//1.读取本地钱包文件，如果文件存在，直接获取
//2.如果文件不存在,创建钱包对象
func NewWallets() *Wallets {
	//wallets := &Wallets{}
	//wallets.WalletMap = make(map[string]*Wallet)
	//return wallets

	//1.判断钱包文件是否存在
	if _, err := os.Stat(WalletFileName); os.IsNotExist(err) {
		fmt.Printf("钱包文件不存在.\n")
		wallets := &Wallets{}
		wallets.WalletMap = make(map[string]*Wallet)
		return wallets
	}

	//2.包文件存在：读取本地的钱包把文件内的数据－－>钱包对象，反序列化
	wsBytes, err := ioutil.ReadFile(WalletFileName)
	if err != nil {
		log.Panic(err)
	}
	gob.Register(elliptic.P256()) //Curve

	var wallets Wallets
	reader := bytes.NewReader(wsBytes)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	return &wallets
}

//创建一个新的钱包
func (ws *Wallets) CreateNewWallet() {
	wallet := NewWallet()

	address := wallet.GetAddress()
	fmt.Printf("创建的钱包的地址:%s\n", address)
	ws.WalletMap[string(address)] = wallet

	//将钱包集合,存入到本地文件中
	ws.saveFile()
}

//将钱包对象存入文件
func (wallets *Wallets) saveFile() {

	//注意:
	//序列化的过程中:被序列化的对象中如果包含了接口， 则接口需要进行注册

	//1.object --> byte[]
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(wallets)
	if err != nil {
		log.Panic(err)
	}

	wsBytes := buf.Bytes()

	//2.将数据存储到文件中
	//注意:该方法的实现:ioutil.WriteFile,覆盖写数据
	err = ioutil.WriteFile(WalletFileName, wsBytes, 0644)
	if err != nil {
		log.Panic(err)
	}
}
