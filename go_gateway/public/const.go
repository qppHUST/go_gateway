package public

const (
	// 验证器key
	ValidatorKey = "ValidatorKey"
	// 翻译器key
	TranslatorKey = "TranslatorKey"
	// amin用户session key
	AdminSessionInfoKey = "AdminSessionInfoKey"

	//service的load方式
	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	RedisFlowDayKey  = "flow_day_count"
	RedisFlowHourKey = "flow_hour_count"

	FlowTotal         = "flow_total"
	FlowServicePrefix = "flow_service_"
	FlowAppPrefix     = "flow_app_"

	JwtSignKey = "my_sign_key"
	JwtExpires = 60 * 60
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)
