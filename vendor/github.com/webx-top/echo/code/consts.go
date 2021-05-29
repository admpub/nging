package code

const (
	// - 系统状态

	SystemUnauthorized Code = -301 // 系统未获得授权
	SystemNotInstalled Code = -300 // 系统未安装

	// - 操作状态

	FrequencyTooFast    Code = -207 //操作频率太快
	OperationProcessing Code = -206 //操作处理中
	RequestFailure      Code = -205 //提交失败
	RequestTimeout      Code = -204 //提交超时
	AbnormalResponse    Code = -203 //响应异常
	OperationTimeout    Code = -202 //操作超时
	Unsupported         Code = -201 //不支持的操作
	RepeatOperation     Code = -200 //重复操作

	// - 数据状态

	InvalidToken Code = -151 //令牌错误
	InvalidAppID Code = -150 //AppID不正确

	DataSizeTooBig      Code = -110 //数据尺寸太大
	DataAlreadyExists   Code = -109 //数据已经存在
	DataFormatIncorrect Code = -108 //数据格式不正确
	DataStatusIncorrect Code = -107 //数据状态不正确
	DataProcessing      Code = -106 //数据未处理中状态
	DataUnavailable     Code = -105 //尚未启用
	DataHasExpired      Code = -104 //数据已经过期
	InvalidType         Code = -103 //类型不正确
	InvalidSignature    Code = -102 //无效的签名
	InvalidParameter    Code = -101 //无效的参数
	DataNotFound        Code = -100 //数据未找到

	// - 用户状态

	CaptchaCodeRequired Code = -11 // captcha code 不能为空
	CaptchaIdMissing    Code = -10 // 缺少captchaId
	CaptchaError        Code = -9  //验证码错误
	BalanceNoEnough     Code = -5  //余额不足
	UserDisabled        Code = -4  //用户被禁用
	UserNotFound        Code = -3  //用户未找到
	NonPrivileged       Code = -2  //无权限
	Unauthenticated     Code = -1  //未登录

	// - 通用

	Failure Code = 0 //操作失败
	Success Code = 1 //操作成功
)
