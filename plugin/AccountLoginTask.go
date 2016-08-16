package plugin

import (
	"github.com/hailongz/kk-go-task/task"
)

type AccountLoginTaskResult struct {
	Result
	User *User `json:"user,omitempty"`
}

/**
 * 账号登录
 */
type AccountLoginTask struct {
	task.Task
	Name     string                 `json:"name"`
	Password string                 `json:"password"`
	Result   AccountLoginTaskResult `json:"-"`
}

func (T *AccountLoginTask) API() string {
	return "account/login"
}

func (T *AccountLoginTask) GetResult() interface{} {
	return &T.Result
}
