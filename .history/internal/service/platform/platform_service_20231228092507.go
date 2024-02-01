package platform

import (
	"wstester/internal/entity"
	"wstester/internal/repo/platform"
	"wstester/internal/service/login"
	"wstester/internal/service/video"
	"wstester/pkg/log"
)

type PlatformService struct {
	platformRepo *platform.PlatformRepo
	loginService *login.LoginService
	videoService *video.VideoService
}

func NewPlatformService(platformRepo *platform.PlatformRepo) *PlatformService {
	return &PlatformService{
		platformRepo: platformRepo,
	}
}

func (p *PlatformService) ConnectToPlatform(token string, platformInfo *entity.Platform) error {
	go func() { _ = p.updatePlatformDB(platformInfo) }()

	return p.connectToPlatform(platformInfo)
}

func (p *PlatformService) connectToPlatform(platformInfo *entity.Platform) error {
	wsUrl, err := p.loginService.Login(platformInfo)
	if err != nil {
		return err
	}
	log.Infof("fetch ws url: %s", wsUrl)
	p.videoService.RegisterWsUrl(wsUrl)
	return nil
}

// 判断数据库是否存在一致的平台
// 如果不存在的话，将其插入到数据库中
func (p *PlatformService) updatePlatformDB(platformInfo *entity.Platform) error {
	isExist, err := p.platformRepo.IsExistSamePlatfrom(platformInfo)
	if err != nil {
		return err
	}

	if !isExist {
		if err = p.platformRepo.Insert(platformInfo); err != nil {
			return err
		}
	}
	return nil
}

func (p *PlatformService) GetPlatformList(platformInfo *entity.Platform, pageIndex, pageSize int) ([]entity.Platform, error) {
	return p.platformRepo.QueryPlatformList(platformInfo, pageIndex, pageSize)
}
