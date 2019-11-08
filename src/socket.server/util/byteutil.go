package util

import(
	"fmt"
	"encoding/binary"
)
func init(){
	fmt.Println("byteutil init")
}

func Hello(){
	fmt.Println("byteutil hello")		
}

func ReadUint16(b []byte) uint16{
	return binary.LittleEndian.Uint16(b)
}

func WriteUint16(b []byte,num uint16){
	binary.LittleEndian.PutUint16(b,num)
} 
