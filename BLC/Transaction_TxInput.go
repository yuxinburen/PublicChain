package BLC

type TxInput struct {
	TxID      []byte //引用的TxOutput所在的交易ID
	Vout      int    //txoutput 下标
	ScriptSip string //解锁脚本
}

//判断TxInput是否是指定用户的消费
func (TxInput *TxInput) UnlockWithAddress(address string) bool {
	return TxInput.ScriptSip == address
}
