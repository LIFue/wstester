package controller

import (
	"net/rpc"
	"wstester/pkg/log"
)

type ControllerRegister struct {
	controllerServices []interface{}
}

func NewControllerRegister() *ControllerRegister {
	return &ControllerRegister{
		controllerServices: make([]interface{}, 0),
	}
}

func (c *ControllerRegister) AddService(i interface{}) {
	c.controllerServices = append(c.controllerServices, i)
}

func (c *ControllerRegister) Register() error {
	log.Infof("register controller")
	for _, service := range c.controllerServices {
		if err := rpc.Register(service); err != nil {
			return err
		}
	}
	return nil
}
