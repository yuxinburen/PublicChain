package BLC

import "bytes"

type TxOutput struct {
	Value int64 //金额
	//ScriptPubKey string //锁定脚本，公钥
	PubKeyHash []byte //公钥哈希
}

//判断TxOutput是否是指定的用户解锁
func (txOutput *TxOutput) UnlockWithAddress(address string) bool {
	//return txOutput.ScriptPubKey == address
	full_payload := Base58Encode([]byte(address))
	pubKeyHash := full_payload[1 : len(full_payload)-addressCheckSumLen]
	return bytes.Compare(pubKeyHash, txOutput.PubKeyHash) == 0
}

//根据转账的目标账户和转账的数量,生成一个TxOutput
func NewTxOutput(amount int64, toAddress string) *TxOutput {
	txOutput := &TxOutput{amount, nil}
	txOutput.Lock(toAddress)
	return txOutput
}

//根据转账地址进行计算出公钥哈希并赋值给txtOutput对象
func (output *TxOutput) Lock(address string) {
	full_payload := Base58Encode([]byte(address))
	output.PubKeyHash = full_payload[1 : len(full_payload)-addressCheckSumLen]
}
