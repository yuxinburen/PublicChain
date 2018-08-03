package BLC

type TxOutput struct {
	Value        int64    //金额
	ScriptPubKey string //锁定脚本，公钥
}
