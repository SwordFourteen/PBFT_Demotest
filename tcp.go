package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

var num = 0
var average int64 = 0
var replyCount = (nodeCount-1)/3 + 1

//客户端使用的tcp监听
func clientTcpListen() {
	listen, err := net.Listen("tcp", clientAddr)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close() //结束监听

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		clientmespool[num] = BytesToInt(b)
		println(clientmespool[num])
		num++

		if num == nodeCount*LoopNum-1 {

			for i := replyCount - 1; i < nodeCount*LoopNum; i += nodeCount {
				average += clientmespool[i]
			}
			average /= 20
			fmt.Printf("20次的平均时间是%d\n", average)
			average = 0
		}

	}
}

//节点使用的tcp监听
func (p *pbft) tcpListen() {
	listen, err := net.Listen("tcp", p.node.addr)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("节点开启监听，地址：%s\n", p.node.addr)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		p.handleRequest(b)
	}

}

//错误节点的tcp发送
func tcpWrong(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect error", err)
		return
	}

	_, err = conn.Write([]byte("我是作恶节点"))
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}

//使用tcp发送消息
func tcpDial(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect error", err)
		return
	}
	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
func BytesToInt(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}
