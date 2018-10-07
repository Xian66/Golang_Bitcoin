package BLC

//未花费的模型

type UTXO struct {
	TxHash []byte
	Index int
	Output *TXOutput
}
