// Package configs
// @Author      : lilinzhen
// @Time        : 2022/5/7 23:49:34
// @Description :
package configs

import "blueblog/pkg/errors"

const (
	ProjectName = "blueblog"

	AppNameForAdmin     = ProjectName + "-admin"
	AppNameForInterface = ProjectName + "-interface"
	AppNameForJob       = ProjectName + "-job"
	AppNameForService   = ProjectName + "-service"
	AppNameForTask      = ProjectName + "-task"

	CacheKey4Article = "blues:" + ProjectName + ":article:%s:%s"

	MQExchange = "blues:" + ProjectName
	MQRoutingKey4Article = MQExchange + ":article:%s"

	// HeaderSignToken 签名验证 Token，Header 中传递的参数
	HeaderSignToken = "Authorization"

	// HeaderSignTokenDate 签名验证 Date，Header 中传递的参数
	HeaderSignTokenDate = "Authorization-Date"

	// ZhCN 简体中文 - 中国
	ZhCN = "zh-cn"

	// EnUS 英文 - 美国
	EnUS = "en-us"
)

var ErrNotExist = errors.New("Not Exist")