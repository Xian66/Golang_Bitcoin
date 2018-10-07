package BLC

type TXOutput struct {
	//代表钱包有多少钱
	Value int64
	// 用户名 公钥
	ScriptPubKey string
}



//解锁
func (txOutput *TXOutput) UnLockScriptPubKeyWithAddress(address string) bool {

	return txOutput.ScriptPubKey == address
}
