package ssh

import (
	"fmt"
	"wstester/internal/base/data"
	"wstester/internal/entity"
)

type SshHostRepo struct {
	d *data.Data
}

func NewSshHostRepo(d *data.Data) *SshHostRepo {
	return &SshHostRepo{
		d: d,
	}
}

func (s *SshHostRepo) Insert(sshHost *entity.SshHost) error {
	if err := s.d.DB.Create(sshHost).Error; err != nil {
		return err
	}
	return nil
}

func (s *SshHostRepo) Update(sshHost *entity.SshHost) error {
	if err := s.d.DB.Updates(sshHost).Error; err != nil {
		return err
	}
	return nil
}

func (s *SshHostRepo) Delete(sshHost *entity.SshHost) error {
	if err := s.d.DB.Delete(sshHost).Error; err != nil {
		return err
	}
	return nil
}

func (s *SshHostRepo) Fetch(sshHost *entity.SshHost, pageIndex, pageSize int) ([]entity.SshHost, error) {
	var out []entity.SshHost
	tx := s.d.DB.Model(&entity.SshHost{})
	if len(sshHost.Ip) > 0 {
		tx = tx.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", sshHost.Ip))
	}
	if pageSize > 0 && pageIndex > 0 {
		tx = tx.Limit(pageSize).Offset(pageIndex)
	}
	if err := tx.Find(&out).Error; err != nil {
		return out, err
	}

	return out, nil
}
