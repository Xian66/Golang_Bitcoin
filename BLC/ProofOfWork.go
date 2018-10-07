package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

//0000 0000 0000 0000 1001 0001 ...0001
//希望256位Hash里面前面至少要有16个零
const targetBit = 12
type ProofOfWork struct{
	Block *Block //当前要验证的区块
	//diff int64  //难度值
	target  *big.Int    //代表数据难度  big.Int大数据存储 解决数据溢出的问题
}

func (pow *ProofOfWork) prepareData(nonce int)[]byte  {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)
	return data

}

func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	//1.把所有的Block的属性拼接成字节数组
	//2.生成hash
	//3.判断hash有效性， 如果满足条件， 跳出循环
	nonce := 0
	var hashInt big.Int  //存储新生成的hash
	var hash [32]byte
	for{
		//准备数据
		dataBytes := proofOfWork.prepareData(nonce)
		//生成hash
		//32位的
		hash = sha256.Sum256(dataBytes)
		//%x在原有数据的地方刷新
		fmt.Printf("%x\r", hash)
		//将hash存储到hashInt
		//估计是要传地址才能改变值， 传地址的话不会修改原有的值
		hashInt.SetBytes(hash[:])
		//判断hashInt 是否小于Block里面的target
		//Cmp compares x and y returns:
		//-1 if x < y
		//0 if x == y
		//1 x > y
		if proofOfWork.target.Cmp(&hashInt) == 1 {
				break
		}
		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}
//创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	//1. big.Int对象
	//0000 0001 希望生成hash前面至少有2个零
	//左移8-2=6位
	//0100 0000  64 就为target
	//生成的hash 前两个为0  0010 0000 32 比target小
	//1.创建一个初始值为1的target
	target := big.NewInt(1)
	//2.左移 256 - targetBit = target
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{block, target}
}

