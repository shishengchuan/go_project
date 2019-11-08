package main

import(
	"fmt"
	"socket.server/util"
	"net"
	"io"
	"socket.server/proto"
	"encoding/json"
)

const network,address string ="tcp", "127.0.0.1:8886"
const bufSize int = 1024
const msgLenSize uint = 2
const protoLenSize uint = 2

var handler map[string]func(net.Conn ,[]byte)
var conns map[net.Conn]net.Conn

func init(){
	handler = make(map[string]func(net.Conn,[]byte))
	handler["fire"] = handleFire


	conns = make(map[net.Conn]net.Conn)
}

func main(){
	fmt.Println("Server Start")

	listen , err := net.Listen(network,address)
	util.CheckError(err)

	defer listen.Close()

	for{
		conn , err := listen.Accept()
		if err != nil {
			continue
		}
		if conns[conn] == nil {
			conns[conn] = conn
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn){
	fmt.Println("local addr " , conn.LocalAddr() , " remote addr " , conn.RemoteAddr())
	allBuf := make([]byte,0)
	buf := make([]byte ,bufSize)
	for{
		n,err := conn.Read(buf)	
		if err == io.EOF {
			conn.Close()
			delete(conns,conn)
			break
		}
		if err != nil {
			conn.Close()
			delete(conns,conn)
			break;
		}
		allBuf = append(allBuf,buf[:n]...)
		dataLen := len(allBuf)
		for dataLen > int(msgLenSize) {
			msgLen := util.ReadUint16(allBuf[0:msgLenSize])
			if int(dataLen) > int(msgLen) - int(msgLenSize) {
				protoNameLen := util.ReadUint16(allBuf[msgLenSize:msgLenSize+protoLenSize])
				protoName := string(allBuf[uint(msgLenSize) + uint(protoLenSize):uint(msgLenSize) + uint(protoLenSize) + uint(protoNameLen)])
				fmt.Println("Proto name = " , protoName) 
				f , ok := handler[protoName]
				if ok {
					bodyStart := uint(msgLenSize) +uint(protoLenSize) + uint(protoNameLen)
					f(conn,allBuf[bodyStart:uint(msgLen) + uint(msgLenSize)])
				}
				allBuf = allBuf[uint(msgLen) + uint(msgLenSize):]
				dataLen = len(allBuf)
			}else{
				break
			}
		}
	}
}

func handleFire(conn net.Conn ,data []byte){
	var fire proto.Fire
	json.Unmarshal(data,&fire)

	var fireBack proto.Fire
	fireBack.X = 100.0
	fireBack.Y = 200.0
	fireBack.Z = 300.0
	fireBack.ProtoName ="MsgFireGo"
	msgData , err := json.Marshal(&fireBack)
	if err == nil {
		protoName := "MsgFireGo"
		nameByte := []byte(protoName)
		nameLen := len(nameByte)
		dataLen := len(msgData)
		totalLen := nameLen + dataLen +2
		sendData := make([]byte,4)
		util.WriteUint16(sendData[0:2],uint16(totalLen))
		util.WriteUint16(sendData[2:4],uint16(nameLen))
		sendData = append(sendData,nameByte...)
		sendData = append(sendData,msgData...)
		sendN , err := sendMsg(conn,sendData)
		if err !=nil {
			fmt.Println("send error" ,err)
		}else{
			fmt.Println("Send num = " , sendN)
		}

	}else{
		fmt.Println("error " , err)
	}
}

func sendMsg(conn net.Conn , data []byte)(int,error){
	fmt.Println("send str " , string(data))
	n ,err := conn.Write(data)
	return n ,err
}



