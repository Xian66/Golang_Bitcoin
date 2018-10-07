package BLC

import "fmt"

//先用来查询余额
func (cli *CLI) getBalance(address string){

	fmt.Println("地址： "+address)
	blockchain := BlockchainObject()
	defer blockchain.DB.Close()

	amount := blockchain.GetBalance(address)

	fmt.Printf("%s一共有%d个Token\n", address, amount)
	/*txOutputs := blockchain.UnUTXOs(address)

	for _,out := range txOutputs  {
		fmt.Println(out)
	}*/
}
