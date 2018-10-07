package BLC

import (
	"fmt"
	"os"
	"log"
	"flag"
)

type CLI struct {
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\tcreateblockchain -address --交易数据")
	fmt.Println("\tsend -from FROM -to TO -amount Amount -交易明细。。")
	fmt.Println("\tprintchain -- 输出区块信息。。。")
	fmt.Println("\tgetbalance -address -- 输出区块信息。。。")
}

func isValidArgs() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}



func (cli *CLI) Run() {
	isValidArgs()
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账地址来源")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址")
	flageAmount := sendBlockCmd.String("amount", "", "转账的金额")

	flagCreateBlockchainWithAddress := createBlockchainCmd.String("address", "", "创建创世区块的地址")
	getbalanceWithAddress := getbalanceCmd.String("address", "", "查询某个账号的余额")

	switch os.Args[1] {
	case "send":
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:

		printUsage()
		os.Exit(1)
	}

	if sendBlockCmd.Parsed() {

		//如果输入为空 就打印提示信息
		if *flagFrom == "" || *flagTo == "" || *flageAmount == "" {
			printUsage()
			os.Exit(1)
		}
		//三个数组的个数应该相同

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)
		amount := JSONToArray(*flageAmount)
		cli.send(from, to, amount)

	}

	if printChainCmd.Parsed() {
		//fmt.Println("输出所有区块的数据。。。")
		cli.printchain()
	}

	if createBlockchainCmd.Parsed() {
		if *flagCreateBlockchainWithAddress == "" {
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		//调用创世区块的方法
		cli.createGenesisBlockchain(*flagCreateBlockchainWithAddress)

	}

	if getbalanceCmd.Parsed() {
		if *getbalanceWithAddress == "" {
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		//调用创世区块的方法
		cli.getBalance(*getbalanceWithAddress)

	}
}
