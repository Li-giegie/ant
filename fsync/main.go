package main

import (
	"encoding/json"
	"flag"
	"github.com/Li-giegie/sanJiaoMaoTCPNet"
	"github.com/Li-giegie/sanJiaoMaoTCPNet/Message"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
)

type AntFile struct {
	FileName string
	Data []byte
	Root int
}

func main(){
	raddr :=flag.String("radd","","remote address")
	key :=flag.String("key","ant-fsync","唯一标识 key ")
	Authentication :=flag.String("Authentication","","认证密文")
	flag.Parse()
	log.Printf("raddr %v key %v Authentication %v\n",*raddr,*key,*Authentication)
	Fsync(*raddr,*key,*Authentication)
}

// 被同步端
func Fsync(radd,key,Authentication string){
	cli:= sanJiaoMaoTCPNet.NewClient(radd,key)
	cli.AddHandlerFunc("sync", func(msg *Message.Message, reply Message.ReplyMessageI) {
		var Anfs []AntFile
		err := json.Unmarshal(msg.Data,&Anfs)
		if err != nil {
			reply.String(201,"json uerr" + err.Error())
			return
		}
		var errInfo string
		for _, anf := range Anfs {
			dir,fn := path.Split(strings.ReplaceAll(anf.FileName,`\`,"/"))
			err = os.MkdirAll(dir,fs.FileMode(anf.Root))
			if err != nil {
				log.Println("mkdir err：",err)
				continue
			}

			err :=os.WriteFile(dir+fn,anf.Data,fs.FileMode(anf.Root))
			if err != nil {

				errInfo = errInfo +anf.FileName +" write err:" + err.Error()
			}
		}
		if errInfo == ""{ errInfo = "success" }
		err = reply.String(200,"sync>>"+errInfo)
		if err != nil {
			log.Println("fsync err:",err)
		}

	})

	err := cli.Connect("sanjiaomao")
	if err != nil {
		log.Fatalln(err)
	}

	cli.Run()
}