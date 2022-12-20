package main

import (
	"flag"
	"github.com/Li-giegie/sanJiaoMaoTCPNet"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"log"
)



func main()  {
	addr :=flag.String("address","0.0.0.0:9999","listen address：0.0.0.0:9999")
	key :=flag.String("key","ant-server","唯一标识 key ")
	Authentication :=flag.String("Authentication","","认证密文")
	flag.Parse()
	log.Printf("raddr %v key %v Authentication %v\n",*addr,*key,*Authentication)
	Server(*addr,*key,*Authentication)
}

func Server(addr,key,Authentication string)  {
	srv,err := sanJiaoMaoTCPNet.NewServer(addr,key)
	if err != nil {
		panic(any(err))
	}

	srv.SetAuthentication(func(ip string, key string, data []byte) (bool, string) {
		log.Printf("ip :%v key :%v data :%v\n",ip,key,string(data))
		if Authentication != ""{
			if string(data) != Authentication {
				log.Println("拒绝认证-------")
				return false,"非法链接"
			}
		}
		return true,"success"
	})

	srv.AddHandleFunc("test", func(msg *Message.Message, reply Message.ReplyMessageI) {
		msg.Data = append(msg.Data, []byte("\n---------ok")...)
		err := reply.Bytes(200,msg.Data)
		if err != nil {
			log.Println(err)
		}
	})

	srv.Run()
}
