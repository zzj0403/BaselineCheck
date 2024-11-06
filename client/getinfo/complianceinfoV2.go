package getinfo

type SingleComplianceInfoV2 struct {
	Name   string `json:"name"`   // 检查项名称
	Level  string `json:"level"`  // 检查项级别
	Action string `json:"action"` // 检查项操作
	Status string `json:"status"` // 检查项状态

}
