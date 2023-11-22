package schemas

type UnitParams struct {
	UnitName   string `json:"unit_name" binding:"required"`
	ScanSource string `json:"scan_source"` // 扫描模式，qcc/tyc
}
