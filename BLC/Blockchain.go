package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"math/big"
	"time"
	"os"
	"strconv"
	"encoding/hex"
)

const dbName = "blockchain.db"  //数据库名字
const blockTableName = "blocks" //表名
type Blockchain struct {
	Tip []byte //存储最新的区块的hash
	DB  *bolt.DB
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

//判断数据库是否存在 要不要创建创世区块
func DBExists() bool {
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}

//遍历输出所有区块的信息
func (blc *Blockchain) PrintChain() {
	//
	blockchainIterattor := blc.Iterator()
	for {
		block := blockchainIterattor.Next()

		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("prevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println("Txs: ")

		for _, tx := range block.Txs {
			fmt.Printf("%x\n", tx.TxHash)
			fmt.Println("Vins: ")

			for _, in := range tx.Vins {
				fmt.Printf("%x\n", in.TxHash)
				//
				fmt.Printf("%d\n", in.Vout)
				fmt.Printf("%s\n", in.ScriptSig)
			}

			fmt.Println("Vouts: ")

			for _, out := range tx.Vouts {
				fmt.Println(out.Value)
				fmt.Println(out.ScriptPubKey)
			}
		}
		fmt.Println("------------------------------------------")

		var hasInt big.Int
		hasInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hasInt) == 0 {
			break;
		}
	}
}

//增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {

	//	在创建创世区块的时候有了数据库 现在直接使用就好
	err := blc.DB.Update(func(tx *bolt.Tx) error {

		//1.获取表
		b := tx.Bucket([]byte(blockTableName))
		//2.创建新区快
		if b != nil {
			//先获取最新区块
			blockBytes := b.Get(blc.Tip)
			//反序列化
			block := DeserializeBlock(blockBytes)
			newBlock := NewBlock(txs, block.Height+1, block.Hash)
			//3.将区块序列化并且存储到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//4.更新数据库了里面“l"对应的hash值
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//5.更新blockchain的Tip
			blc.Tip = newBlock.Hash
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

//首先要有创建区块链的方法
//1.创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock(address string) *Blockchain {
	//判断数据库是否存在 如果存在直接退出
	if DBExists() {
		fmt.Println("创世区块已经存在。。。")
		os.Exit(1)
	}

	//创建数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var genesisHash []byte

	err = db.Update(func(tx *bolt.Tx) error {

		//创世区块的时候还么有表创建数据库表
		fmt.Println("正在创建创世区块")
		b, err := tx.CreateBucket([]byte(blockTableName))

		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			//创建创世区块
			//创建coinbase Transaction
			txConinbase := NewCoinbaseTransaction(address)
			genesisBlock := CreateGenesisBlock([]*Transaction{txConinbase})
			/*将创世区块序列化存储到表中*/
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//	存储最新的区块hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			genesisHash = genesisBlock.Hash

		}
		return nil
	})
	return &Blockchain{genesisHash, db}
}

//返回Blockchain对象
func BlockchainObject() *Blockchain {
	//创建Blockchain对象并返回
	//创建数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//获取最新区块的hash
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &Blockchain{tip, db}
}

//若果一个地址对应的TXOutput未花费， 那么把这个output添加到数组中返回
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {
	//存储未花费的transaction
	var unUTXOs []*UTXO

	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs{


		if tx.IsCoinbaseTransaction() ==false{

			for _, in := range tx.Vins {
				//是否能够解锁
				if in.UnLockWithAddress(address) {
					//转换成字符串
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	}

	for _,tx := range txs {
		//判断out是否存在spentTXOutput[key]中

		Work1:
		for index, out := range tx.Vouts {

			if out.UnLockScriptPubKeyWithAddress(address){

				if len(spentTXOutputs) == 0 {

					utxo := &UTXO{tx.TxHash, index,out}
					unUTXOs = append(unUTXOs, utxo)

				} else{
					for hash, indexArray := range spentTXOutputs{

						txHashStr := hex.EncodeToString(tx.TxHash)

						if hash == txHashStr {
							var isUnSpentUTXO bool
							for _, outIndex := range indexArray {

								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}

								if isUnSpentUTXO ==false {
									utxo := &UTXO{tx.TxHash, index,out}
									unUTXOs = append(unUTXOs, utxo)
								}
							}
						}else {
							utxo := &UTXO{tx.TxHash, index,out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}


			}

		}
	}



	//通过迭代器取遍历
	blockIterator := blockchain.Iterator()
	for {
		block := blockIterator.Next()
		fmt.Println(block)
		fmt.Println()

		//判断tx是否有效
		//有问题要处理

		for /*i := len(block.Txs) - 1; i >= 0; i--*/_, tx := range block.Txs {
			//	txHash
			//Vins
			if tx.IsCoinbaseTransaction() ==false{

				for _, in := range tx.Vins {
					//是否能够解锁
					if in.UnLockWithAddress(address) {
						//转换成字符串
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}

			//Vouts
			work:
			for index, out := range tx.Vouts {

				if out.UnLockScriptPubKeyWithAddress(address) {
					//	判断output是否被消费
					if spentTXOutputs != nil {

						if len(spentTXOutputs) != 0 {
							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {

								for _, i := range indexArray {
									//说明这笔钱已经被花费
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}

							if isSpentUTXO ==false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}
			}

		}

		fmt.Println("未花费的输出", spentTXOutputs)
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}
	}
	return unUTXOs
}

//转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	//	1.获取所有的UTXO
	utxos := blockchain.UnUTXOs(from, txs)

	spendAbleUTXO := make(map[string][]int)
	//遍历utxos
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.Output.Value

		hash := hex.EncodeToString(utxo.TxHash)
		spendAbleUTXO[hash] = append(spendAbleUTXO[hash], utxo.Index)
		if value >= int64(amount) {
			break
		}
	}

	if value < int64(amount) {
		fmt.Printf("%s's fund is not enough", from)
		os.Exit(1)
	}

	return value, spendAbleUTXO
}

//产生新的区块
func (blockchain *Blockchain) MineNewBlock(from, to, amount []string) {
	//建立一笔交易

	fmt.Println("From: ", from)
	fmt.Println("To: ", to)
	fmt.Println("Amount: ", amount)

	var txs []*Transaction
	for index, address := range from{
		value, _ := strconv.Atoi(amount[index])
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)
		txs = append(txs, tx)
		//fmt.Println(tx)
	}

	//1. 通过相关算法建立Transaction数组

	var block *Block
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	//2.构造新的区块 新区块的高度和前一个区块hash可以从数据库中读取
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区快存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		//先获取到表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())

			b.Put([]byte("l"), block.Hash)
			//把最新区块的hash存到Tip
			blockchain.Tip = block.Hash
		}
		return nil
	})
}

//查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {

	utxos := blockchain.UnUTXOs(address, []*Transaction{})

	var amount int64
	for _, utxo := range utxos {
		amount = amount + utxo.Output.Value
	}

	return amount
}
