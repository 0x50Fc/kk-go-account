package plugin

import (
	"github.com/hailongz/kk-go-task/task"
)

type AccountCreateTaskResult struct {
	Result
	User *User `json:"user,omitempty"`
}

/**
 * 创建账号
 */
type AccountCreateTask struct {
	task.Task
	Name     string `json:"name"`
	Password string `json:"password"`
	Result   AccountCreateTaskResult
}

func (T *AccountCreateTask) API() string {
	return "account/create"
}

func (T *AccountCreateTask) GetResult() interface{} {
	return &T.Result
}
