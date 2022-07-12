package constant

/* 匹配 */
const (
	// Value 值列表
	Value = "value"
	// IsBlank 字符为空匹配
	IsBlank = "isBlank"
	// IsUnBlank 字符为非空匹配
	IsUnBlank = "isUnBlank"
	// Range 范围匹配
	Range = "range"
	// Model 固定的几个模式匹配
	Model = "model"
	// Condition 条件表达式
	Condition = "condition"
	// Regex 正则表达式
	Regex = "regex"
	// Customize 自定义函数回调
	Customize = "customize"
)

/* 匹配后处理 */
const (
	// ErrMsg 自定义错误异常
	ErrMsg = "errMsg"
	// Accept 匹配后是否接受
	Accept = "accept"
	// Disable 是否启用属性的核查功能
	Disable = "disable"
)

/* tag关键字 */
const (
	EQUAL = "="
	MATCH = "match"
	CHECK = "check"
)

/* range匹配关键字 */
const (
	LeftEqual    = "["
	LeftUnEqual  = "("
	RightUnEqual = ")"
	RightEqual   = "]"

	Now    = "now"
	Past   = "past"
	Future = "future"
)

/* model类别 */
const (
	IdCard     = "id_card"
	Phone      = "phone"
	FixedPhone = "fixed_phone"
	MAIL       = "mail"
	IpAddress  = "ip"
)
