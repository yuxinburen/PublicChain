package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const TargetBit = 16 //目标hash的0的个数,16,20,24,28....

//工作量证明算法结构
type ProfOfWork struct {
	Block  *Block   //要验证的区块对象
	Target *big.Int //目标hash
}

//创建一个新的PoW算法对象,然后设置其属性并返回
func NewProfOfWork(block *Block) *ProfOfWork {

	//1.创建pow对象
	pow := &ProfOfWork{}
	//设置属性值
	pow.Block = block
	target := big.NewInt(1) //目标hash值，初始值为1
	target.Lsh(target, 256-TargetBit)
	pow.Target = target

	return pow
}

//PoW算法挖矿的过程方法,目标是得到复合挖矿算法的nonce值，并返回hash
func (pow *ProfOfWork) Run() ([]byte, int64) {
	//挖矿的过程及目标:更改nonnce的值,计算hash,直到小于目标hash
	//挖矿思路:
	//1.设置nonce值
	//2.block的属性字段拼接，并计算hash值
	//3.比较实际的hash值和目标hash
	var nonce int64 = 0
	var hash [32]byte
	for { //由于需要多次尝试nonce的值，因此是个死循环
		//1.根据nonce值，将整个块的数据转换成[]byte形式的数据
		data := pow.PrepareData(nonce)
		//2.生成hash
		hash = sha256.Sum256(data)
		fmt.Printf("\r%d,%x", nonce, hash)
		//3.验证：将计算得到的hash和目标hash进行比较
		hashInt := new(big.Int)
		hashInt.SetBytes(hash[:])
		if (pow.Target.Cmp(hashInt)) == 1 {
			break
		}
		nonce++
	}
	fmt.Println()
	return hash[:], nonce
}

//根据变化的nonce，返回整个块的hash计算值
func (work *ProfOfWork) PrepareData(nonce int64) []byte {
	//1.根据nonce,生成pow中要验证的block数组
	data := bytes.Join([][]byte{
		IntToHex(work.Block.Height),
		work.Block.PreBlockHash,
		IntToHex(work.Block.TimeStamp),
		work.Block.Data,
		IntToHex(nonce),
		IntToHex(TargetBit),
	}, []byte{})

	return data
}
