package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
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
