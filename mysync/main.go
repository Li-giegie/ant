package main

import (
	"encoding/json"
	"flag"
	"fmt"
	go_reload "github.com/Li-giegie/go-reload"
	"github.com/Li-giegie/sanJiaoMaoTCPNet"
	"time"

	"log"
	"os"
)

type AntFile struct {
	FileName string
	Data []byte
	Root int
}

func main()  {
	raddr :=flag.String("radd","","remote address")
	key :=flag.String("key","ant-sync","唯一标识 key ")
	Authentication :=flag.String("Authentication","","认证密文")
	requestKey :=flag.String("requestKey","ant-fsync","被同步端的key")

	flag.Parse()
	log.Printf("raddr %v key %v Authentication %v\n",*raddr,*key,*Authentication)

	cli := sanJiaoMaoTCPNet.NewClient(*raddr,*key)

	err := cli.Connect(*Authentication)
	if err != nil {
		log.Fatalln(err)
	}

	r := go_reload.New()

	for _, s := range r.GetFileList() {
		buf,err := os.ReadFile(s)
		if err != nil {
			log.Println("os.read err:",err)
			continue
		}
		jbuf,jerr := json.Marshal(&[]AntFile{{
			FileName: s,
			Data:     buf,
			Root:     0666,
		},
		})
		if jerr != nil {
			log.Println("os.read err:",err)
			continue
		}
		fmt.Println(s, len(buf))
		reply,err := cli.SendMessage(*requestKey,"sync",200,jbuf,time.Second*60)
		if err != nil {
			log.Println("request err:",err)
			continue
		}

		log.Println(reply.String(),string(reply.Data))
	}

	r.Run(func(cf []string) {
		fmt.Println("刷新")
		var anf = make([]AntFile,0)
		for _, s := range cf {

			buf,err := os.ReadFile(s)
			if err != nil {
				log.Println(err)
				continue
			}
			anf = append(anf, AntFile{
				FileName: s,
				Data:     buf,
				Root: 0666,
			})

			anfBuf ,err := json.Marshal(anf)
			if err != nil {
				log.Println("json m err :",err)
				continue
			}

			reply,err := cli.SendMessage(*requestKey,"sync",200,anfBuf,time.Second*60)
			if err != nil {
				fmt.Println("请求失败！",err)
				return
			}
			log.Println(reply.String())
			log.Println(string(reply.Data))
		}
	})
}
