package controller

import (
	"wstester/internal/entity"
	"wstester/internal/schema"
	"wstester/internal/service/platform"
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
	resp.Status = "success"

	err = p.platformService.ConnectToPlatform("", &req.Platform)
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
