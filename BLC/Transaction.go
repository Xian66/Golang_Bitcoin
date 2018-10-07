package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

//tIn = TXInput{[]byte{}, -1, "Genesis Block"}
//txOut = TXOutput{10, "csc"}

//Transaction{"000000", []*Transaction{tIn}, []*TXOutput{txOut}
//交易信息 UTXO
type Transaction struct {
	//	1.交易hash
	TxHash []byte
	//2. 输入
	Vins []*TXInput

	//3.输出
	Vouts []*TXOutput
}

//判断当前的交易是否是Coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {

	//return false
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

//Transaction 创建分两种情况
//1.创世区块创建时的Transaction
func NewCoinbaseTransaction(address string) *Transaction {

	//代表消费
	txInput := &TXInput{[]byte{}, -1, "Genesis data..."}
	//
	txOutput := &TXOutput{10, address}

	txCoinbase := &Transaction{[]byte{}, []*TXInput{txInput}, []*TXOutput{txOutput}}
	//设置hash值
	txCoinbase.HashTransaction()
	return txCoinbase
}

//将交易数据序列化 然后再做hash
func (tx *Transaction) HashTransaction() {
	//缓冲区
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	//将区块序列化
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}

//2.转账时产生的Transaction

func NewSimpleTransaction(from, to string, amount int, blockchain *Blockchain, txs []*Transaction) *Transaction {

	//./main.exe  send -from '["bob"]' -to '["heng"]' -amount '["2"]'
	//from [bob]
	//to [heng]
	//amount [2]

	//1.有一个函数， 返回from这个人所有的未花费交易输出所对应的transaction
	//unUTXOs := blockchain.UnUTXOs(from)

	// 通过一个函数， 返回money， dic
	money, spendableUTXODic := blockchain.FindSpendableUTXOS(from, amount, txs)

	//money 是dic里面hash的总值
	//dic里的内容 两笔交易在同一区块中的时候有一样的哈希 交易的index不同：{hash1:[0,2]}
	//	两笔交易在不同的区块中的时候 ： {hash1:[0], hash2:[2]}

	var txInputs []*TXInput
	var txOutputs []*TXOutput
	//代表消费
	//bytes, _ := hex.DecodeString("6b715cad442465f42473fcde3e3ddc57eca2d048b8b57af4bb7455b6e9c7ada6")

	for txHash, indexArray :=range spendableUTXODic  {

		txHashBytes,_ := hex.DecodeString(txHash)
		for _, index := range indexArray {
			txInput := &TXInput{txHashBytes, index, from}
			//消费
			txInputs = append(txInputs, txInput)
		}
	}


	//转账
	txOutput := &TXOutput{int64(amount), to}
	txOutputs = append(txOutputs, txOutput)
	//找零
	txOutput = &TXOutput{int64(money) - int64(amount), from}
	txOutputs = append(txOutputs, txOutput)

	tx := &Transaction{[]byte{}, txInputs, txOutputs}
	//设置当前的hash
	tx.HashTransaction()
	return tx
}
