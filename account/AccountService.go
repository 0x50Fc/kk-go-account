package account

import (
	"encoding/json"
	"fmt"
	kkdb "github.com/hailongz/kk-go-db/kk"
	"github.com/hailongz/kk-go-task/task"
	"log"
	"math/rand"
	"time"
)

type AccountService struct {
	task.Service
}

func (S *AccountService) Handle(task task.ITask) error {
	return S.ReflectHandle(task, S)
}

/**
 * 创建账号
 */
func (S *AccountService) HandleAccountCreateTask(task *AccountCreateTask) error {

	var plugin = S.Plugin().(*Plugin)
	var db = plugin.Db

	if task.Name == "" {
		task.Result.Errno = ERRNO_NOT_FOUND_NAME
		task.Result.Errmsg = "未找到账号名"
		return nil
	}

	var password = task.Password

	if password == "" {
		password = EncodePassword(fmt.Sprintf("%d %d", time.Now().Unix(), rand.Intn(100000)))
	} else {
		password = EncodePassword(task.Password)
	}

	var user = User{0, task.Name, password, 0, 0, time.Now().Unix()}

	var _, err = kkdb.DBInsert(db, &UserTable, plugin.Prefix, &user)

	if err != nil {
		task.Result.Errno = ERRNO_EXISTS_NAME
		task.Result.Errmsg = "账号名已存在"
		return nil
	}

	task.Result.User = &user

	log.Println("AccountService.HandleAccountCreateTask")
	log.Println(user)

	return nil
}

/**
 * 修改账号
 */
func (S *AccountService) HandleAccountSetTask(task *AccountSetTask) error {

	var plugin = S.Plugin().(*Plugin)
	var db = plugin.Db

	if task.Uid == 0 {
		task.Result.Errno = ERRNO_NOT_FOUND_UID
		task.Result.Errmsg = "未找到用户ID"
		return nil
	}

	var rows, err = kkdb.DBQuery(db, &UserTable, plugin.Prefix, " WHERE uid=?", []interface{}{task.Uid})

	if err != nil {
		task.Result.Errno = ERRNO_ACCOUNT
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	var user = User{}
	var scaner = kkdb.NewDBScaner(&user)

	if rows.Next() {
		err = scaner.Scan(rows)
		if err != nil {
			task.Result.Errno = ERRNO_ACCOUNT
			task.Result.Errmsg = err.Error()
			return nil
		}
	} else {
		task.Result.Errno = ERRNO_NOT_FOUND_ACCOUNT
		task.Result.Errmsg = "未找到账号"
		return nil
	}

	if task.Password == "" {
		user.Password = EncodePassword(fmt.Sprintf("%d %d", time.Now().Unix(), rand.Intn(100000)))
	} else {
		user.Password = EncodePassword(task.Password)
	}

	user.Mtime = time.Now().Unix()

	_, err = kkdb.DBUpdate(db, &UserTable, plugin.Prefix, &user)

	if err != nil {
		task.Result.Errno = ERRNO_ACCOUNT
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.User = &user

	log.Println("AccountService.HandleAccountSetTask")
	log.Println(user)

	return nil
}

/**
 * 账号登录
 */
func (S *AccountService) HandleAccountLoginTask(task *AccountLoginTask) error {

	var plugin = S.Plugin().(*Plugin)
	var db = plugin.Db

	if task.Name == "" {
		task.Result.Errno = ERRNO_NOT_FOUND_NAME
		task.Result.Errmsg = "未找到账号名"
		return nil
	}

	var rows, err = kkdb.DBQuery(db, &UserTable, plugin.Prefix, " WHERE name=? AND password=?", task.Name, EncodePassword(task.Password))

	if err != nil {
		task.Result.Errno = ERRNO_ACCOUNT
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	var user = User{}
	var scaner = kkdb.NewDBScaner(&user)

	if rows.Next() {
		err = scaner.Scan(rows)
		if err != nil {
			task.Result.Errno = ERRNO_ACCOUNT
			task.Result.Errmsg = err.Error()
			return nil
		}
	} else {
		task.Result.Errno = ERRNO_NOT_FOUND_ACCOUNT
		task.Result.Errmsg = "未找到账号"
		return nil
	}

	user.Atime = time.Now().Unix()

	_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET atime=? WHERE uid=?", plugin.Prefix, UserTable.Name), user.Atime, user.Uid)

	if err != nil {
		task.Result.Errno = ERRNO_ACCOUNT
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.User = &user

	log.Println("AccountService.HandleAccountLoginTask")
	log.Println(user)

	return nil
}

/**
 * 设置用户信息
 */
func (S *AccountService) HandleAccountInfoSetTask(task *AccountInfoSetTask) error {

	var plugin = S.Plugin().(*Plugin)
	var db = plugin.Db

	if task.Uid == 0 {
		task.Result.Errno = ERRNO_NOT_FOUND_UID
		task.Result.Errmsg = "未找到用户ID"
		return nil
	}

	if task.Name == "" {
		task.Result.Errno = ERRNO_NOT_FOUND_NAME
		task.Result.Errmsg = "未找到用户信息命名"
		return nil
	}

	var rows, err = kkdb.DBQuery(db, &UserInfoTable, plugin.Prefix, " WHERE name=? AND uid=?", task.Name, task.Uid)

	if err != nil {
		task.Result.Errno = ERRNO_ACCOUNT
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	var userInfo = UserInfo{}
	var scaner = kkdb.NewDBScaner(&userInfo)

	if rows.Next() {
		err = scaner.Scan(rows)
		if err != nil {
			task.Result.Errno = ERRNO_ACCOUNT
			task.Result.Errmsg = err.Error()
			return nil
		}
		var b, _ = json.Marshal(task.Value)
		userInfo.Value = string(b)
		_, err = kkdb.DBUpdate(db, &UserInfoTable, plugin.Prefix, &userInfo)
		if err != nil {
			task.Result.Errno = ERRNO_ACCOUNT
			task.Result.Errmsg = err.Error()
			return nil
		}
	} else {
		var b, _ = json.Marshal(task.Value)
		userInfo.Uid = task.Uid
		userInfo.Name = task.Name
		userInfo.Value = string(b)
		_, err = kkdb.DBInsert(db, &UserInfoTable, plugin.Prefix, &userInfo)
		if err != nil {
			task.Result.Errno = ERRNO_ACCOUNT
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	_, err = db.Exec(fmt.Sprintf("UPDATE %s%s SET mtime=? WHERE uid=?", plugin.Prefix, UserTable.Name), time.Now().Unix(), task.Uid)

	if err != nil {
		task.Result.Errno = ERRNO_ACCOUNT
		task.Result.Errmsg = err.Error()
		return nil
	}

	log.Println("AccountService.HandleAccountInfoSetTask")
	log.Println(task)

	return nil
}

/**
 * 获取用户信息
 */
func (S *AccountService) HandleAccountInfoTask(task *AccountInfoTask) error {

	var plugin = S.Plugin().(*Plugin)
	var db = plugin.Db

	if task.Uid == 0 {
		task.Result.Errno = ERRNO_NOT_FOUND_UID
		task.Result.Errmsg = "未找到用户ID"
		return nil
	}

	if task.Name == "" {
		task.Result.Errno = ERRNO_NOT_FOUND_NAME
		task.Result.Errmsg = "未找到用户信息命名"
		return nil
	}

	var rows, err = kkdb.DBQuery(db, &UserInfoTable, plugin.Prefix, " WHERE name=? AND uid=?", task.Name, task.Uid)

	if err != nil {
		task.Result.Errno = ERRNO_ACCOUNT
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	var userInfo = UserInfo{}
	var scaner = kkdb.NewDBScaner(&userInfo)

	if rows.Next() {
		err = scaner.Scan(rows)
		if err != nil {
			task.Result.Errno = ERRNO_ACCOUNT
			task.Result.Errmsg = err.Error()
			return nil
		}
		var err = json.Unmarshal([]byte(userInfo.Value), &task.Result.Value)
		if err != nil {
			task.Result.Errno = ERRNO_ACCOUNT
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	log.Println("AccountService.HandleAccountInfoTask")
	log.Println(task)

	return nil
}
