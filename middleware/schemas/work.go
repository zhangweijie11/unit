package schemas

type UnitParams struct {
	KeyWord        string   `json:"keyword"`
	CompanyID      string   `json:"company_id"` // 公司 ID
	CompanyName    string   `json:"company_name"`
	ScanSource     []string `json:"scan_source" binding:"required"` // 扫描模式，qcc/tyc/aqc
	IsDetail       bool     `json:"is_detail"`
	IsMergeOut     bool     `json:"is_merge_out"`
	SearchField    []string `json:"search_field"`
	ISKeyPid       bool     `json:"is_key_pid"`
	IsGroup        bool     `json:"is_group"`
	IsGetBranch    bool     `json:"is_get_branch"`
	IsSearchBranch bool     `json:"is_search_branch"`
	IsInvestRd     bool     `json:"is_invest_rd"`
	IsEmailPro     bool     `json:"is_email_pro"`
	IsHold         bool     `json:"is_hold"`
	IsSupplier     bool     `json:"is_supplier"`
	InvestNum      float64  `json:"invest_num"`
	Deep           int      `json:"deep"`
}
