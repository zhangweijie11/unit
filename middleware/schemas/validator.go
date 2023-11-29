package schemas

var taskValidatorErrorMessage = map[string]string{
	"ScanSourcerequired": "缺少扫描源",
}

// RegisterValidatorRule 注册参数验证错误消息, Key = e.StructNamespace(), value.key = e.Field()+e.Tag()
var RegisterValidatorRule = map[string]map[string]string{
	"UnitParams": taskValidatorErrorMessage,
}

const (
	WorkKeyIDErr       = "关键词或 PID 需至少存在一个"
	WorkScanSourceErr  = "无效的扫描源"
	WorkResultFieldErr = "无效的结果字段"
)

func ValidParamsExist(param1, param2 string) bool {
	return param1 != "" || param2 != ""
}
