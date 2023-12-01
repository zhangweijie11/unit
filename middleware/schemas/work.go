package schemas

type Cookies struct {
	Aldzs               string `yaml:"aldzs" json:"aldzs" mapstructure:"aldzs"`
	Xlb                 string `yaml:"xlb" json:"xlb" mapstructure:"xlb"`
	Aiqicha             string `yaml:"aiqicha" json:"aiqicha" mapstructure:"aiqicha"`
	Tianyancha          string `yaml:"tianyancha" json:"tianyancha" mapstructure:"tianyancha"`
	TianyanchaTycid     string `yaml:"tianyanchatycid" json:"tianyanchatycid" mapstructure:"tianyanchatycid"`
	TianyanchaAuthToken string `yaml:"tianyanchaauthtoken" json:"tianyanchaauthtoken" mapstructure:"tianyanchaauthtoken"`
	Qichacha            string `yaml:"qichacha" json:"qichacha" mapstructure:"qichacha"`
	QiMai               string `yaml:"qimai" json:"qimai" mapstructure:"qimai"`
	ChinaZ              string `yaml:"chinaz" json:"chinaz" mapstructure:"chinaz"`
	Veryvp              string `yaml:"veryvp" json:"veryvp" mapstructure:"veryvp"`
}

type UnitParams struct {
	KeyWord     string   `json:"keyword"`    // 关键词 eg 小米
	CompanyID   string   `json:"company_id"` // 公司 ID
	CompanyName string   `json:"company_name"`
	ScanSource  []string `json:"scan_source" binding:"required"` // 扫描模式，qcc/tyc/aqc
	IsMergeOut  bool     `json:"is_merge_out"`                   // 批量查询【取消】合并导出
	ResultField []string `json:"result_field"`                   // 获取字段信息 eg icp
	IsBranch    bool     `json:"is_branch"`                      // 查询分支机构（分公司）信息
	IsHold      bool     `json:"is_hold"`                        // 是否查询控股公司
	IsSupplier  bool     `json:"is_supplier"`                    // 是否查询供应商信息
	InvestNum   float64  `json:"invest_num"`                     // 投资比例 eg 100
	Deep        int      `json:"deep"`                           // 递归搜索n层公司
	Cookies     Cookies  `json:"cookies"`                        // Cookies
}
