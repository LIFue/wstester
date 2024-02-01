package controller

import "wstester/pkg/log"

type Args struct {
	A, B int
}

type Arith struct {
}

func NewArith() *Arith {
	a := &Arith{}
	controllerRegister.AddService(a)
	return a
}

func (t *Arith) Multiply(args *Args, reply *int) error {

	log.Info("multiply")
	*reply = args.A * args.B
	return nil
}
