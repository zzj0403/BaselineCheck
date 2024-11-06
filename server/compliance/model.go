package compliance

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ComplianceResult struct {
	IP            string `json:"ip" gorm:"comment:IP地址"`
	Hostname      string `json:"hostname" gorm:"comment:主机名"`
	BaselineCount int    `json:"baseline_count" gorm:"comment:基线检查项数"`

	gorm.Model
}

func (u *ComplianceResult) TableName() string {
	return "compliance_result"
}

type ComplianceDetails struct {
	Details  datatypes.JSON `json:"result" gorm:"type:json;comment:检查结果"`
	ResultId uint           `json:"result_id" gorm:"comment:检查结果ID"`
	gorm.Model
}

func (u *ComplianceDetails) TableName() string {
	return "compliance_details"
}
