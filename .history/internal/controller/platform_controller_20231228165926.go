package controller

import (
	"wstester/internal/entity"
	"wstester/internal/schema"
	"wstester/internal/service/platform"
	"wstester/pkg/log"
)

type Platform struct {
	platformService *platform.PlatformService
}

func NewPlatform(platformService *platform.PlatformService, r *ControllerRegister) *Platform {
	p := &Platform{
		platformService: platformService,
	}

	r.AddService(p)
	return p
}

func (p *Platform) ConnectPlatform(req *schema.ReqConnectPlatform, resp *schema.RespConnectPlatform) (err error) {
	log.Infof("receive a request to connect to platform: %+v", req)
	resp.Status = "success"

	err = p.platformService.ConnectToPlatform(req.JosnID, &req.Platform)
	if err != nil {
		resp.Status = "failed"
	}

	return err
}

func (p *Platform) GetPlatformList(req *schema.ReqGetPlatformList, resp *schema.RespGetPlatformList) (err error) {
	platformInfo := &entity.Platform{
		Ip: req.Ip,
	}
	resp.PlatformList, err = p.platformService.GetPlatformList(platformInfo, req.PageIndex, req.PageSize)
	return
}

func (p *Platform) SendMessage(req *schema.ReqSendMessage, resp *schema.RespSendMessage) (err error) {
	var message string
	message, err = p.platformService.SendMessage(req.JosnID, req.Message)
	if err != nil {
		return
	}

	resp.message = message
	return
}
