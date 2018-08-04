package BLC

//UTXO模型
//UTXO:Unspent transaction output //未花费的交易输出
type UTXO struct {
	//1.该output所在的交易ID
	TxID []byte
	//2.该output的下标
	Index int
	//3.未花费的
	Output *TxOutput
}
