package schemas

var taskValidatorErrorMessage = map[string]string{
	"unit_namerequired": "缺少任务单位名称",
}

// RegisterValidatorRule 注册参数验证错误消息, Key = e.StructNamespace(), value.key = e.Field()+e.Tag()
var RegisterValidatorRule = map[string]map[string]string{
	"UnitParams": taskValidatorErrorMessage,
}
