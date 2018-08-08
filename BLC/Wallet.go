package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressCheckSumLen = 4

//钱包对象
type Wallet struct {
	PrivateKey ecdsa.PrivateKey //私钥
	Publickey  []byte           //公钥
}

//获取得到钱包的地址
func (wallet *Wallet) GetAddress() []byte {
	//思路:
	//1.原始公钥-->sha256--->sha256 --->公钥哈希
	//2.版本号+公钥哈希-->校验码
	//3.版本号+公钥哈希+校验码-->Base58编码

	//原始公钥－－>公钥哈希
	pubKeyHash := PubKeyHash(wallet.Publickey)

	//2.添加版本号
	version_payload := append([]byte{version}, pubKeyHash...)

	//两次sha256 并获得checksum
	checkSumBytes := CheckSum(version_payload)

	//拼接全部数据
	full_payload := append(version_payload, checkSumBytes...)
	fmt.Printf("full_payload:", full_payload, ",len:", len(full_payload))

	//Base58编码
	address := Base58Encode(full_payload)
	return address
}

//两次sha256,产生校验码
func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:addressCheckSumLen]
}

//sha256
func PubKeyHash(publickKey []byte) []byte {

	hasher := sha256.New()
	hasher.Write(publickKey)
	hasher1 := hasher.Sum(nil)

	//再一次sha256
	hasher2 := ripemd160.New()
	hasher2.Write(hasher1)
	hashs := hasher2.Sum(nil)

	return hashs
}

//创建钱包对象
func NewWallet() *Wallet {
	//思路:
	//1.准一对私钥和公钥
	//2.构造wallet对象进行返回
	pritvateKey, publicKey := NewPaireKey()
	wallet := &Wallet{PrivateKey: pritvateKey, Publickey: publicKey}
	return wallet
}

func NewPaireKey() (ecdsa.PrivateKey, []byte) {

	curve := elliptic.P256() //生成一个椭圆
	//privateKey
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	//publicKey
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	//返回私钥，公钥
	return *privateKey, publicKey
}
