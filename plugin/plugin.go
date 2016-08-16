package plugin

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"github.com/hailongz/kk-go-db/kk"
	"github.com/hailongz/kk-go-task/task"
)

type User struct {
	Uid      int64  `json:"uid"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Mtime    int64  `json:"mtime"` //修改时间
	Atime    int64  `json:"atime"` //访问时间
	Ctime    int64  `json:"ctime"` //创建时间
}

type UserInfo struct {
	Uiid  int64  `json:"uiid"`
	Uid   int64  `json:"uid"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

var UserTable = kk.DBTable{"user",

	"uid",

	map[string]kk.DBField{"name": kk.DBField{256, kk.DBFieldTypeString},
		"password": kk.DBField{32, kk.DBFieldTypeString},
		"mtime":    kk.DBField{0, kk.DBFieldTypeInt},
		"atime":    kk.DBField{0, kk.DBFieldTypeInt},
		"ctime":    kk.DBField{0, kk.DBFieldTypeInt}},

	map[string]kk.DBIndex{"name": kk.DBIndex{"name", kk.DBIndexTypeAsc, true}}}

var UserInfoTable = kk.DBTable{"userinfo",

	"uiid",

	map[string]kk.DBField{"uid": kk.DBField{0, kk.DBFieldTypeInt64},
		"name":  kk.DBField{64, kk.DBFieldTypeString},
		"value": kk.DBField{0, kk.DBFieldTypeText}},

	map[string]kk.DBIndex{"uid": kk.DBIndex{"uid", kk.DBIndexTypeAsc, false},
		"name": kk.DBIndex{"name", kk.DBIndexTypeAsc, false}}}

func EncodePassword(password string) string {
	m := md5.New()
	m.Write([]byte(password))
	m.Write([]byte("K*)()I<MHGEW@#$%^&*IO"))
	v := m.Sum(nil)
	return hex.EncodeToString(v)
}

type Plugin struct {
	Db     *sql.DB
	Prefix string
}

func Load(context *task.Context) error {

	var db = context.Get("db").(*sql.DB)
	var p = Plugin{db, context.Get("prefix").(string)}

	var err = kk.DBBuild(db, &UserTable, p.Prefix, 1)

	if err != nil {
		return err
	}

	err = kk.DBBuild(db, &UserInfoTable, p.Prefix, 1)

	if err != nil {
		return err
	}

	context.Plugin(&p)(&AccountService{})(&AccountCreateTask{}, &AccountSetTask{}, &AccountInfoTask{}, &AccountInfoSetTask{}, &AccountLoginTask{})

	return nil
}
