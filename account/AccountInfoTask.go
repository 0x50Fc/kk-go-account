package account

import (
	"github.com/hailongz/kk-go-task/task"
)

type AccountInfoTaskResult struct {
	Result
	Value map[string]interface{} `json:"value,omitempty"`
}

/**
 * 获取账号信息
 */
type AccountInfoTask struct {
	task.Task
	Uid    int64                 `json:"uid"`
	Name   string                `json:"name"`
	Result AccountInfoTaskResult `json:"-"`
}

func (T *AccountInfoTask) API() string {
	return "account/info/get"
}

func (T *AccountInfoTask) GetResult() interface{} {
	return &T.Result
}
