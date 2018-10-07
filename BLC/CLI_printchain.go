package BLC

import (
	"fmt"
	"os"
)

func (cli *CLI) printchain() {
	if !DBExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()
	blockchain.PrintChain()

}

