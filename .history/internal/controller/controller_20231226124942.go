package controller

import "net/rpc"

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
	for _, service := range c.controllerServices {
		if err := rpc.Register(service); err != nil {
			return err
		}
	}
	return nil
}
