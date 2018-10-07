package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"fmt"
)

//将出入的int64转换为字节数组
func IntToHex(num int64) []byte{
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//标准的JSON字符串转数组
func JSONToArray(jsonStr string)[]string {
	//fmt.Println("json string: ", jsonStr)
//	json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonStr), &sArr); err != nil{
		//在golang的Terminal中编译执行会报错
		//invalid character '\'' looking for beginning of value[]
		// 切换到powershell或者bash下就正常
		log.Panic(err)
		fmt.Printf("%v", err)
	}
	return sArr
}
