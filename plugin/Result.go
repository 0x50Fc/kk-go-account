package plugin

type IResultTask interface {
	GetResult() interface{}
}

type Result struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}

/**
 * 账号错误
 */
const ERRNO_ACCOUNT = 0x1000

/**
 * 未找到账号名
 */
const ERRNO_NOT_FOUND_NAME = ERRNO_ACCOUNT + 1

/**
 * 未找到UID
 */
const ERRNO_NOT_FOUND_UID = ERRNO_ACCOUNT + 2

/**
 * 未找到账号
 */
const ERRNO_NOT_FOUND_ACCOUNT = ERRNO_ACCOUNT + 3

/**
 * 账号名已存在
 */
const ERRNO_EXISTS_NAME = ERRNO_ACCOUNT + 4
