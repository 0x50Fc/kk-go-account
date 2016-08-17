package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hailongz/kk-go-account/account"
	kkdb "github.com/hailongz/kk-go-db/kk"
	"github.com/hailongz/kk-go-task/task"
	"github.com/hailongz/kk-go/kk"
	"log"
	"os"
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

	var replay func(message *kk.Message) bool = nil

	replay, _ = kk.TCPClientConnect(name, address, map[string]interface{}{"exclusive": true}, func(message *kk.Message) {

		if message.Method != "REQUEST" {
			log.Println(message)
			return
		}

		var apiname = message.To[len(name):]
		var tk = context.NewAPITask(apiname)

		if tk == nil {
			var v = kk.Message{"NOIMPLEMENT", message.To, message.From, "text", []byte(apiname)}
			log.Println(v)
			replay(&v)
			return
		} else if message.Type == "text/json" {
			var err = json.Unmarshal(message.Content, tk)
			if err != nil {
				var b, _ = json.Marshal(&account.Result{account.ERRNO_ACCOUNT, err.Error()})
				var v = kk.Message{message.Method, message.To, message.From, "text/json", b}
				replay(&v)
				return
			}
		}

		log.Println(tk)

		go func() {
			var err = context.Handle(tk)
			if err != nil && err != task.ERROR_BREAK {
				var b, _ = json.Marshal(&account.Result{account.ERRNO_ACCOUNT, err.Error()})
				var v = kk.Message{message.Method, message.To, message.From, "text/json", b}
				kk.GetDispatchMain().Async(func() {
					replay(&v)
				})
				return
			} else {
				var rs, ok = tk.(account.IResultTask)
				if ok {
					var b, _ = json.Marshal(rs.GetResult())
					var v = kk.Message{message.Method, message.To, message.From, "text/json", b}
					kk.GetDispatchMain().Async(func() {
						replay(&v)
					})
				} else {
					var v = kk.Message{message.Method, message.To, message.From, "text/json", []byte("{}")}
					kk.GetDispatchMain().Async(func() {
						replay(&v)
					})
				}
			}
		}()

	})

	context.Set("replay", replay)

	account.Load(context)

	kk.DispatchMain()

}
