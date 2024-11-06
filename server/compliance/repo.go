package compliance

import (
	"BaselineCheck/server/repository"
	"errors"

	"gorm.io/gorm"
)

type Repo struct {
	repo repository.MyDB
}

func NewRepo(repo repository.MyDB) *Repo {
	return &Repo{
		repo: repo,
	}
}
func (r *Repo) CreateComplianceResult(ins *ComplianceResult) error {
	return r.repo.Create(ins).Error
}

func (r *Repo) CreateComplianceDetails(ins *ComplianceDetails) error {
	return r.repo.Create(ins).Error
}

func (r *Repo) GetComplianceResultByHostname(hostname string) (*ComplianceResult, error) {
	var result ComplianceResult
	if err := r.repo.Where("hostname = ?", hostname).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &result, nil

}

func (r *Repo) GetComplianceDetailsByResultIdandDetailId(result_id, id string) (*ComplianceDetails, error) {
	var result ComplianceDetails
	if err := r.repo.Where("result_id = ? AND id = ?", result_id, id).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &result, nil

}

// 更新
func (r *Repo) UpdateComplianceResultByHostname(ins *ComplianceResult) error {
	return r.repo.Save(ins).Error
}

func (r *Repo) GetComplianceHostList(page, pageSize int) ([]ComplianceResult, error) {
	var result []ComplianceResult
	offset := (page - 1) * pageSize
	limit := pageSize
	if err := r.repo.Offset(offset).Limit(limit).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil

}
func (r *Repo) GetComplianceDetailByhostId(host_id string) ([]ComplianceDetails, error) {
	var result []ComplianceDetails
	if err := r.repo.Debug().Where("result_id = ?", host_id).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil

}
