package platform

import (
	"fmt"

	"wstester/internal/base/data"
	"wstester/internal/entity"
	"wstester/pkg/log"

	"gorm.io/gorm"
)

type PlatformRepo struct {
	data *data.Data
}

func NewPlatformRepo(data *data.Data) *PlatformRepo {
	return &PlatformRepo{
		data: data,
	}
}

func (p *PlatformRepo) Insert(platform *entity.Platform) error {
	if err := p.data.DB.Create(&platform).Error; err != nil {
		log.Errorf("create platform error: %s", err.Error())
		return err
	}
	return nil
}

func (p *PlatformRepo) IsExistSamePlatfrom(platform *entity.Platform) (bool, error) {
	var platformCount int64
	if err := p.data.DB.Where(platform).Count(&platformCount).Error; err != nil {
		return false, err
	}

	return platformCount > 0, nil
}

func (p *PlatformRepo) QueryPlatformList(platformInfo *entity.Platform, pageIndex, pageSize int) ([]entity.Platform, error) {
	var platformList []entity.Platform
	tx := p.data.DB.Model(&entity.Platform{})
	if len(platformInfo.Ip) > 0 {
		tx = tx.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", platformInfo.Ip))
	}
	if err := tx.Limit(pageSize).Offset(pageIndex).Find(&platformList).Error; err != nil {
		log.Errorf("query platform list by ip: %s pageIndex: %d pageSize: %d error: %s", platformInfo.Ip, pageIndex, pageSize, err.Error())
		return platformList, err
	}

	return platformList, nil
}

func (p *PlatformRepo) FetchPlatform(platformID uint) (entity.Platform, error) {
	var platform entity.Platform
	platform.ID = platformID

	if err := p.data.DB.Model(&entity.Platform{}).Preload("User").First(&platform).Error; err != nil && err != gorm.ErrRecordNotFound {
		return platform, err
	}

	return platform, nil
}
