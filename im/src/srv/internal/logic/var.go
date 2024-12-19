package logic

import "errors"

var (
	ErrMobileRegisted = errors.New("手机号已注册")
	ErrUserNotFind    = errors.New("未找到该用户")
)
