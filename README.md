# starChain
* 把所有的区块存储到数据库

使用bolt数据库 key-value

install

```
 go get github.com/boltdb/bolt/...
```

* 区块存储的格式

```go
Block1{
    preHash
    Hash
}
Block2{
    preHash
    Hash
}
```

区块是一个对象 存储的时候  

1.先把区块序列化成字节数组[]byte

2. 以当前区块的hash作为key， 序列化后的字节数组为value存储到数据库中

```
dp.put(block.Hash, block.Seri())
```

数据在数据库中有点像map是无序的

通过block.Hash作为key获取序列化后的block 然后通过反序列化 获取block

##### 在数据库中遍历整个所有区块

为了能够找到上一个区块 需要在数据库中存储最新区块的hash  以key-value的形式存储

key可以自定义 但是key不改变 ，  遍历区块的时候 通过该key找到最新区块的hash值 然后获取最新区块 该区块中存储了上一个区块的hash 以此类推 知道当前区块的上一个区块的hash为0 就到了创世区块 这样就能遍历整个区块

##### 数据库

存数据的时候先判断数据库中表是否存在， 不存在就创建一个

默认会在当前可执行文件的同级目录生成.db文件  权限的最大值为777

```go
//数据库操作
func main() {
	//my.db  
	//数据库不存在的话会创建一个，  mode:权限 最大为777
	//创建或打开数据库
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//关闭数据库
	defer db.Close()
	//	 更新表数据库
	err = db.View(func(tx *bolt.Tx) error {
		//获取BlockBucket表
		b := tx.Bucket([]byte("BlockBucket"))
		if b != nil {
			data := b.Get([]byte("l"))
			fmt.Printf("%s", data)
			data = b.Get([]byte("ll"))
			fmt.Printf("%s", data)

			//返回nil， 以便数据库处理相应操作
		}
		return nil
	})
	if err != nil {
		//更新数据失败
		log.Panic(err)
	}
}
```



##### 区块数据存放在数据库后的区块的遍历

区块存储到数据库中之后 Blockchain中不需要block数组保存区块 而是添加一个Tip字段用来存储最新的区块hash 每次添加区块的时候就把hash值保存到Tip中



##### 区块的遍历

创世区块中的PrevBlockHash是个全是0的数组 如果当前区块的PrevBlockHash为0 则当前区块是创世区块  可以退出循环

```go
var hashInt big.Int		
hashInt.SetBytes(block.PrevBlockHash)
if big.NewInt(0).Cmp(&hashInt) == 0{
			break;
		}
```

##### 解析命令行

* 添加区块： bc addBlock -data "abc"
* 打印区块 ： bc printchain



#### Transaction-UTXO

```go
type Block struct {
	//	1.区块高度 int64
	Height int64
	//2.上一个区块的hash 字节数组
	PrevBlockHash []byte
	//3.交易数据  会有多笔交易
	Txs []*Transaction
	//4. 时间戳 int64
	Timestamp int64
	//5.Hash 字节数组
	Hash []byte
	//6.Nonce值
	Nonce int64
}
```



```go
//交易信息 UTXO
type Transaction struct{

//	1.交易hash
	TxHash []byte
//2. 输入


//3.输出 
	Vouts []*TXOutput
}
```

**UTXO**

状态1： 已花费

状态2： 未花费

遍历所有区块 先找到未花费的钱和已花费的钱   如果加起来足够转账的时候就不再搜索了

**coinbase**：

在创建创世区块的时候比较特殊

Transaction 创建分两种情况：
	1.创世区块创建时的Transaction
	2.转账时产生的Transaction



* 假设区块里只有一笔交易的时候  transaction只需要扫描整个数据库 找出所有的utxo  进行交易就行 
* 但是 真实情况下 区块里不止一笔区块 还没打包的transaction 也会对其它的transaction产生影响
* 同时打包多笔交易 还未打包到区块里的交易也要考虑到对其它交易的影响



JSON 解析

```
json.Unmarshal([]byte(jsonStr), &sArr)
```

在golang的Terminal中编译执行会报错 invalid character '\'' looking for beginning of value[]

 切换到powershell或者bash下就正常



#### 遍历

传入名字之后 返回该名字的所有未花费的交易输出



##### 发送交易的格式

./main.exe send -from '["bob"]' -to '["alice"]' -amount '["2"]'