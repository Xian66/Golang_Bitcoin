package BLC

import (
	"time"
	"fmt"
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

type Block struct {
	//	1.区块高度 int64
	Height int64
	//2.上一个区块的hash 字节数组
	PrevBlockHash []byte
	//3.交易数据 把Txs里的hash拼接起来
	Txs []*Transaction
	//4. 时间戳 int64
	Timestamp int64
	//5.Hash 字节数组
	Hash []byte
	//6.Nonce值
	Nonce int64
}

//需要将Txs转换成 []byte
func (block *Block) HashTransactions() []byte{

	var txHashes [][]byte
	var txHash [32]byte
	for _, tx := range block.Txs{
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

//将区块序列化成字节数组
func (block *Block) Serialize() []byte {
	//缓冲区
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	//将区块序列化
	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(blockBytes []byte) *Block{
	var block Block
	//反序列化
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	//映射到block中
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//1.创建新的区块
func NewBlock(txs []*Transaction, height int64, prevBlockHash []byte,) *Block{
	//创建区块
	block := &Block{height,prevBlockHash,txs, time.Now().Unix(), nil,0}
	//将PrevBlockHash和Data和Timestamp拼接起来生成hash 再赋值给Hash


	//调用工作量证明的方法 并且 返回有效的Hash和Nonce值
	//挖矿验证
	pow := NewProofOfWork(block)

	//返回符合要求的nonce和hash
	//假设前面有6个零的hash
	hash, nonce := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	fmt.Println()
	return block
}

//生成创世区块
func CreateGenesisBlock(txs []*Transaction) *Block{
	//高度和PrevBlockHash都是已知的
	return NewBlock(txs,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}
