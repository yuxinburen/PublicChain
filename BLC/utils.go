package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)

//int类型的数值转换成16进制表示
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	//将二进制数写入
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

//json格式的数据转换成字符数组
func JSONToArray(jsonData string) []string {
	var arr []string
	err := json.Unmarshal([]byte(jsonData), &arr)
	if err != nil {
		log.Panic(err)
	}
	return arr
}

//字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
