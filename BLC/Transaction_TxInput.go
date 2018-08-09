package BLC

import "bytes"

type TxInput struct {
	TxID []byte //引用的TxOutput所在的交易ID
	Vout int    //txoutput 下标

	//ScriptSip string //解锁脚本
	Signature []byte //数字签名
	PublicKey []byte //原始公钥,钱包里的公钥
}

//判断TxInput是否是指定用户的消费
func (txInput *TxInput) UnlockWithAddress(pubKeyHash []byte) bool {
	pubKeyHash2 := PubKeyHash(txInput.PublicKey)
	return bytes.Compare(pubKeyHash, pubKeyHash2) == 0
}
