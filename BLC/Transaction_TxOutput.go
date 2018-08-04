package BLC

type TxOutput struct {
	Value        int64  //金额
	ScriptPubKey string //锁定脚本，公钥
}


//判断TxOutput是否是指定的用户解锁
func (txOutput *TxOutput) UnlockWithAddress(address string) bool {
	return txOutput.ScriptPubKey == address
}
