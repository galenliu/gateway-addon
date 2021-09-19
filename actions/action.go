package actions

import rpc "github.com/galenliu/gateway-grpc"

type Action struct {
}

func NewActionFormMessage(action *rpc.Action) *Action {
	a := Action{}
	return &a
}
