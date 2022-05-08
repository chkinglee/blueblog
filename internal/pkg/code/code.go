package code

import "blueblog/internal/pkg/configs"

// Failure 错误时返回结构
type Failure struct {
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 描述信息
}

// 业务码
// 1 : 系统级错误
// 2 : 普通错误，多用于接口错误码的定义
// 3 : 任务状态机

// 1 : 系统级错误
const (
	ServerError        = 10101
	TooManyRequests    = 10102
	ParamBindError     = 10103
	AuthorizationError = 10104
	UrlSignError       = 10105
	CacheSetError      = 10106
	CacheGetError      = 10107
	CacheDelError      = 10108
	CacheNotExist      = 10109
	ResubmitError      = 10110
	HashIdsEncodeError = 10111
	HashIdsDecodeError = 10112
	RBACError          = 10113
	RedisConnectError  = 10114
	MySQLConnectError  = 10115
	WriteConfigError   = 10116
	SendEmailError     = 10117
	MySQLExecError     = 10118
	GoVersionError     = 10119
	IdNotFound         = 10120
	NameDuplicate      = 10121
	JsonUnmarshalError = 10122
	UnsupportedError   = 10123

	MetadataGetError     = 10201
	MetadataGetSuccess   = 10202
	SSHError             = 10203
	SSHDownloadError     = 10204
	SSHDownloadSuccess   = 10205
	SSHUploadError       = 10206
	SSHUploadSuccess     = 10207
	ShellError           = 10208
	ReadFileError        = 10209
	ReadFileSuccess      = 10210
	WriteFileError       = 10211
	WriteFileSuccess     = 10212
	ScriptsNotConfig     = 10213
	GrpcConnectError     = 10214
	GrpcConnectSuccess   = 10215
	GrpcCloseError       = 10216
	GrpcCloseSuccess     = 10217
	StateMachineState    = 10218
	ImportPerfClientInfo = 10219
	StateMachineSkip     = 10220
)

// 2 : 普通错误，多用于接口错误码的定义
const (
	AuthorizedCreateError    = 20101
	AuthorizedListError      = 20102
	AuthorizedDeleteError    = 20103
	AuthorizedUpdateError    = 20104
	AuthorizedDetailError    = 20105
	AuthorizedCreateAPIError = 20106
	AuthorizedListAPIError   = 20107
	AuthorizedDeleteAPIError = 20108

	AdminCreateError             = 20201
	AdminListError               = 20202
	AdminDeleteError             = 20203
	AdminUpdateError             = 20204
	AdminResetPasswordError      = 20205
	AdminLoginError              = 20206
	AdminLogOutError             = 20207
	AdminModifyPasswordError     = 20208
	AdminModifyPersonalInfoError = 20209
	AdminMenuListError           = 20210
	AdminMenuCreateError         = 20211
	AdminOfflineError            = 20212
	AdminDetailError             = 20213

	MenuCreateError       = 20301
	MenuUpdateError       = 20302
	MenuListError         = 20303
	MenuDeleteError       = 20304
	MenuDetailError       = 20305
	MenuCreateActionError = 20306
	MenuListActionError   = 20307
	MenuDeleteActionError = 20308

	CronCreateError  = 20401
	CronUpdateError  = 20402
	CronListError    = 20403
	CronDetailError  = 20404
	CronExecuteError = 20405

	ArticleCreateError = 20501
	ArticleUpdateError = 20502
	ArticleDeleteError = 20503
	ArticleListError   = 20504
	ArticleDetailError = 20505
)

func Text(code int) string {
	lang := configs.Get().Language.Local

	if lang == configs.ZhCN {
		return zhCNText[code]
	}

	if lang == configs.EnUS {
		return enUSText[code]
	}

	return zhCNText[code]
}
