package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hailongz/kk-go-account/plugin"
	kkdb "github.com/hailongz/kk-go-db/kk"
	"github.com/hailongz/kk-go-task/task"
	"github.com/hailongz/kk-go/kk"
	"log"
	"os"
	"time"
)

func help() {
	fmt.Println("kk-account <name> <0.0.0.0:8080> <url> <prefix>")
}

func main() {

	log.SetFlags(log.Llongfile | log.LstdFlags)

	var args = os.Args
	var name string = ""
	var address string = ""
	var url string = ""
	var prefix string = ""

	if len(args) > 4 {
		name = args[1]
		address = args[2]
		url = args[3]
		prefix = args[4]
	} else {
		help()
		return
	}

	var cli *kk.TCPClient = nil
	var cli_connect func() = nil
	var db, err = sql.Open("mysql", url)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	_, err = db.Exec("SET NAMES utf8mb4")

	if err != nil {
		log.Fatal(err)
		return
	}

	db.SetMaxIdleConns(6)
	db.SetMaxOpenConns(200)

	err = kkdb.DBInit(db)

	if err != nil {
		log.Fatal(err)
		return
	}

	var context = task.NewContext()

	context.Set("db", db)
	context.Set("prefix", prefix)

	plugin.Load(context)

	cli_connect = func() {
		log.Println("connect " + address + " ...")
		cli = kk.NewTCPClient(name, address, map[string]interface{}{"exclusive": true})
		cli.OnConnected = func() {
			log.Println(cli.Address())
		}
		cli.OnDisconnected = func(err error) {
			log.Println("disconnected: " + cli.Address() + " error:" + err.Error())
			kk.GetDispatchMain().AsyncDelay(cli_connect, time.Second)
		}
		cli.OnMessage = func(message *kk.Message) {

			if message.Method != "REQUEST" {
				log.Println(message)
				return
			}

			var apiname = message.To[len(name):]
			var tk = context.NewAPITask(apiname)

			if tk == nil {
				if cli != nil {
					var v = kk.Message{"NOIMPLEMENT", message.To, message.From, "text", []byte(apiname)}
					log.Println(v)
					cli.Send(&v, nil)
				}
				return
			} else if message.Type == "text/json" {
				var err = json.Unmarshal(message.Content, tk)
				if err != nil {
					var b, _ = json.Marshal(&plugin.Result{plugin.ERRNO_ACCOUNT, err.Error()})
					var v = kk.Message{message.Method, message.To, message.From, "text/json", b}
					cli.Send(&v, nil)
					return
				}
			}

			log.Println(tk)

			go func() {
				var err = context.Handle(tk)
				if err != nil && err != task.ERROR_BREAK {
					var b, _ = json.Marshal(&plugin.Result{plugin.ERRNO_ACCOUNT, err.Error()})
					var v = kk.Message{message.Method, message.To, message.From, "text/json", b}
					kk.GetDispatchMain().Async(func() {
						if cli != nil {
							cli.Send(&v, nil)
						}
					})
					return
				} else {
					var rs, ok = tk.(plugin.IResultTask)
					if ok {
						var b, _ = json.Marshal(rs.GetResult())
						var v = kk.Message{message.Method, message.To, message.From, "text/json", b}
						kk.GetDispatchMain().Async(func() {
							if cli != nil {
								cli.Send(&v, nil)
							}
						})
					} else {
						var v = kk.Message{message.Method, message.To, message.From, "text/json", []byte("{}")}
						kk.GetDispatchMain().Async(func() {
							if cli != nil {
								cli.Send(&v, nil)
							}
						})
					}
				}
			}()

		}
	}

	cli_connect()

	kk.DispatchMain()

}
