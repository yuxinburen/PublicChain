package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"bytes"
)

const version = byte(0x00)
const addressCheckSumLen = 4

//钱包对象
type Wallet struct {
	PrivateKey ecdsa.PrivateKey //私钥
	Publickey  []byte           //公钥
}

//根据给定的地址,判断地址是否合法
//判断地址是否合法的思路:
//1.先将地址进行base58解码
//2.截取checkSum
//3.将剩余的部分进行两次hhash，然后截取4个字节，和第2步截取的进行比较
//4.获取对你结果，一致说明合法，不一致则不合法
func IsValidAddress(address []byte) bool {
	//1.将地址进行base8解码
	full_payload := Base58Decode(address) //25个字节
	//2.获取地址中的CheckSum
	checkSumBytes := full_payload[len(full_payload)-AddressCheckSumLen:]
	//版本号和公钥哈希共同的数据内容
	version_payload := full_payload[:len(full_payload)-AddressCheckSumLen]

	//3.将第二步分离出来的数据进行两次sha256哈希计算
	checkSumBytes2 := CheckSum(version_payload)

	//4.将两个结果进行对比。获取对比结果
	return bytes.Compare(checkSumBytes, checkSumBytes2) == 0
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

//根据公钥获取公钥哈希
func PubKeyHash(publicKey []byte) []byte {
	//一次sha256
	hasher := sha256.New()
	hasher.Write(publicKey)
	hasher1 := hasher.Sum(nil)

	//再一次160
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

func GetAddressByPubKeyHash(pubKeyHash []byte) []byte {
	//step2：添加版本号：
	versioned_payload := append([]byte{version}, pubKeyHash...)

	//step3：根据versioned_payload-->两次sha256,取前4位，得到checkSum
	checkSumBytes := CheckSum(versioned_payload)

	//step4：拼接全部数据
	full_payload := append(versioned_payload, checkSumBytes...)
	//fmt.Println("full_payload:", full_payload, ",len:", len(full_payload))
	//step5：Base58编码
	address := Base58Encode(full_payload)
	return address
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
