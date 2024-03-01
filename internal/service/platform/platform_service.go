package platform

import (
	"fmt"
	"strings"
	"time"
	"wstester/internal/base/code"
	"wstester/internal/base/encrypt"
	"wstester/internal/base/ws"
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
	encrypt        *encrypt.Encrypt

	wsManager *ws.WsManager
}

func NewPlatformService(
	platformRepo *platform.PlatformRepo,
	loginService *login.LoginService,
	videoService *video.VideoService,
	userService *user.UserService,
	messageService *message.MessageService,
	encrypt *encrypt.Encrypt,
	wsManager *ws.WsManager) *PlatformService {
	return &PlatformService{
		platformRepo:   platformRepo,
		loginService:   loginService,
		videoService:   videoService,
		userService:    userService,
		messageService: messageService,
		encrypt:        encrypt,
		wsManager:      wsManager,
	}
}

func (p *PlatformService) ConnectToPlatform(serverID int64, platformInfo *entity.Platform) error {
	go func() {
		if err := p.updatePlatformDB(platformInfo); err != nil {
			log.Errorf("update platform error: %s", err.Error())
		}
		platformInfo.User.PlatformID = platformInfo.ID
		if err := p.userService.AddUser(&platformInfo.User); err != nil {
			log.Errorf("add user error: %s", err.Error())
		}
	}()
	if !strings.HasPrefix(platformInfo.Ip, "10.35") {
		platformInfo.IsPublic = true
	}
	return p.connectToPlatform(serverID, platformInfo)
}

func (p *PlatformService) connectToPlatform(wsId int64, platformInfo *entity.Platform) error {

	wsUrl, err := p.loginService.Login(platformInfo)
	if err != nil {
		return err
	}
	// log.Infof("fetch ws url: %s", wsUrl)
	// platformLoginSign := p.encrypt.Encode(code.EncryptCodeMD5, fmt.Sprintf("%s:%s", platformInfo.Ip, platformInfo.User.Uid))
	// if err = p.videoService.RegisterWsUrl(wsId, platformLoginSign, wsUrl); err != nil {
	// 	return err
	// }
	platformID := fmt.Sprintf("ip:%s:uid:%s", platformInfo.Ip, platformInfo.User.Uid)

	return p.wsManager.InitAndRegisterClient(wsId, platformID, wsUrl, platformInfo.IsPublic)
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

func (p *PlatformService) SendMessage(wsId int64, message entity.MessageEntity) (string, error) {
	go func() {
		if err := p.messageService.UpdateMessage(&message); err != nil {
			log.Errorf("update message info error: %s", err.Error())
		}
	}()
	// return p.videoService.SendMessage(wsId, message.Message)

	c, err := p.wsManager.SendMessage(wsId, "", message.Message)
	if err != nil {
		return "", err
	}

	timeout := time.NewTicker(10 * time.Second)
	for {
		select {
		case data := <-c:
			return string(data), nil
		case <-timeout.C:
			return "", code.ERR_TIMEOUT
		}
	}
}
