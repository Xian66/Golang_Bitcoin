package BLC

import (
	"fmt"
	"os"
)

//转账功能
func (cli *CLI) send(from []string, to []string, amount []string) {
	if !DBExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}

	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	blockchain.MineNewBlock(from, to, amount)
}
