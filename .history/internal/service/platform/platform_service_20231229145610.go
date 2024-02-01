package platform

import (
	"wstester/internal/base/encrypt"
	"wstester/internal/entity"
	"wstester/internal/repo/platform"
	"wstester/internal/service/login"
	"wstester/internal/service/message"
	"wstester/internal/service/user"
	"wstester/internal/service/video"
	"wstester/pkg/log"
)

type PlatformService struct {
	platformRepo   *platform.PlatformRepo
	loginService   *login.LoginService
	videoService   *video.VideoService
	userService    *user.UserService
	messageService *message.MessageService
}

func NewPlatformService(
	platformRepo *platform.PlatformRepo,
	loginService *login.LoginService,
	videoService *video.VideoService,
	userService *user.UserService,
	messageService *message.MessageService,
	encrypt *encrypt.Encrypt) *PlatformService {
	return &PlatformService{
		platformRepo:   platformRepo,
		loginService:   loginService,
		videoService:   videoService,
		userService:    userService,
		messageService: messageService,
	}
}

func (p *PlatformService) ConnectToPlatform(wsId string, platformInfo *entity.Platform) error {
	go func() {
		if err := p.updatePlatformDB(platformInfo); err != nil {
			log.Errorf("update platform error: %s", err.Error())
		}
		platformInfo.User.PlatformID = platformInfo.ID
		if err := p.userService.AddUser(&platformInfo.User); err != nil {
			log.Errorf("add user error: %s", err.Error())
		}
	}()

	return p.connectToPlatform(wsId, platformInfo)
}

func (p *PlatformService) connectToPlatform(wsId string, platformInfo *entity.Platform) error {

	wsUrl, err := p.loginService.Login(platformInfo)
	if err != nil {
		return err
	}
	log.Infof("fetch ws url: %s", wsUrl)
	// p.videoService.CloseOldConnection(wsId)

	if err = p.videoService.RegisterWsUrl(wsId, wsUrl); err != nil {
		return err
	}
	return nil
}

// 判断数据库是否存在一致的平台
// 如果不存在的话，将其插入到数据库中
func (p *PlatformService) updatePlatformDB(platformInfo *entity.Platform) error {
	isExist, err := p.platformRepo.IsExistSamePlatfrom(platformInfo)
	if err != nil {
		log.Errorf("try to judge exist same platform error: %s", err.Error())
		return err
	}

	if !isExist {
		if err = p.platformRepo.Insert(platformInfo); err != nil {
			log.Errorf("try to insert platform error: %s", err.Error())
			return err
		}
	}
	return nil
}

func (p *PlatformService) GetPlatformList(platformInfo *entity.Platform, pageIndex, pageSize int) ([]entity.Platform, error) {
	return p.platformRepo.QueryPlatformList(platformInfo, pageIndex, pageSize)
}

func (p *PlatformService) SendMessage(wsId string, message entity.MessageEntity) (string, error) {
	go func() {
		if err := p.messageService.UpdateMessage(&message); err != nil {
			log.Errorf("update message info error: %s", err.Error())
		}
	}()
	return p.videoService.SendMessage(wsId, message.Message)
}
