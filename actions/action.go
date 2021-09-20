package actions

import rpc "github.com/galenliu/gateway-grpc"

type Action struct {
	Name string `json:"name"`
}

func NewActionFormMessage(action *rpc.Action) *Action {
	a := Action{}
	return &a
}
