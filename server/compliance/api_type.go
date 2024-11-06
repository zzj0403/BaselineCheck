package compliance

import "BaselineCheck/client/baselinelinux"

type RegisterRequest struct {
	baselinelinux.Result
}

type ComplianceHostnameListRequest struct {
	Page `json:"page"`
}

type ComplianceHostnameListRespone struct {
	Page  `json:"page"`
	Total int                `json:"total"`
	List  []ComplianceResult `json:"list"`
}

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func DefaultPage() Page {
	return Page{
		Page:     1,
		PageSize: 10,
	}
}
