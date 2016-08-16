package plugin

import (
	"github.com/hailongz/kk-go-task/task"
)

type AccountSetTaskResult struct {
	Result
	User *User `json:"user,omitempty"`
}

/**
 * 修改账号
 */
type AccountSetTask struct {
	task.Task
	Uid      int64  `json:"uid"`
	Password string `json:"password"`
	Result   AccountSetTaskResult
}

func (T *AccountSetTask) API() string {
	return "account/set"
}

func (T *AccountSetTask) GetResult() interface{} {
	return &T.Result
}
