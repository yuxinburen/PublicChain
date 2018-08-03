package BLC

type TxInput struct {
	TxID      []byte //引用的TxOutput所在的交易ID
	Vout      int    //txoutput 下标
	ScriptSip string //解锁脚本
}
