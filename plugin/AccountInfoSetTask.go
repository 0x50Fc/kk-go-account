package plugin

import (
	"github.com/hailongz/kk-go-task/task"
)

type AccountInfoSetTaskResult struct {
	Result
}

/**
 * 修改账号信息
 */
type AccountInfoSetTask struct {
	task.Task
	Uid    int64                  `json:"uid"`
	Name   string                 `json:"name"`
	Value  map[string]interface{} `json:"value"`
	Result AccountInfoSetTaskResult
}

func (T *AccountInfoSetTask) API() string {
	return "account/info/set"
}

func (T *AccountInfoSetTask) GetResult() interface{} {
	return &T.Result
}
