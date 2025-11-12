package common

// 统一错误码信息

type ErrorCode struct {
	Code int
	Msg  string
}

var (
	// =======================  系统相关 ========================

	UNAUTHORIZED = NewErrorCode(401, "未授权")
	BadRequest   = NewErrorCode(400, "请求错误")
	ServerError  = NewErrorCode(500, "服务器错误")
	Forbidden    = NewErrorCode(403, "无权限访问")
	// =======================  用户相关 ========================

	NotLogin             = NewErrorCode(10001, "未登录")
	UserNotExist         = NewErrorCode(10002, "用户不存在")
	UserExist            = NewErrorCode(10003, "用户已存在")
	UserPasswordError    = NewErrorCode(10004, "用户名不存在或密码错误")
	CaptchaError         = NewErrorCode(10005, "验证码错误")
	GenerateCaptchaError = NewErrorCode(10006, "生成验证码错误")
	CaptchaNotExist      = NewErrorCode(10007, "验证码不存在")
	// =======================  主机相关 ========================
	HostNotExist    = NewErrorCode(20001, "主机不存在")
	HostUnreachable = NewErrorCode(20002, "主机不可达")
)

// NewErrorCode 创建错误码，方便后续业务调用
func NewErrorCode(code int, msg string) *ErrorCode {
	return &ErrorCode{
		Code: code,
		Msg:  msg,
	}
}
