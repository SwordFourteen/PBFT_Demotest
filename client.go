package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

var clientmespool [nodeCount * LoopNum]int64

func clientSendMessageAndListen() {
	//开启客户端的本地监听（主要用来接收节点的reply信息）
	go clientTcpListen()
	fmt.Printf("客户端开启监听，地址：%s\n", clientAddr)

	fmt.Println(`---------------------------------------------------------------------------------
					已进入权重制PBFT测试客户端，

					---------------------------------------------------------------------------------
					请在下方输入要存入节点的信息：`)
	//首先通过命令行获取用户输入
	stdReader := bufio.NewReader(os.Stdin) //监听输入区
	for {
		data, err := stdReader.ReadString('\n') //以回车作为结束符进行标记
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		for i := 0; i < LoopNum; i++ {
			r := new(Request)                    //新建请求
			r.Timestamp = time.Now().UnixMicro() //获取当前时间微秒为单位
			r.ClientAddr = clientAddr
			r.Message.ID = getRandom() //随机，但不知道为什么
			//消息内容就是用户的输入
			r.Message.Content = strings.TrimSpace(data) //去空格
			br, err := json.Marshal(r)
			if err != nil {
				log.Panic(err)
			}
			fmt.Println(string(br))
			content := jointMessage(cRequest, br) //将命令阶段和消息链接（但我不知道为什么不直接append）
			//默认N0为主节点，直接把请求信息发送至N0
			println("发送request\n")
			tcpDial(content, nodeTable["N0"])
			time.Sleep(30 * time.Second)
		}
	}
}

//返回一个十位数的随机数，作为msgid
func getRandom() int {
	x := big.NewInt(10000000000)
	for {
		result, err := rand.Int(rand.Reader, x)
		if err != nil {
			log.Panic(err)
		}
		if result.Int64() > 1000000000 {
			return int(result.Int64())
		}
	}
}
