package addon

import (
	"github.com/galenliu/gateway-grpc"
	"google.golang.org/grpc"
)

type RpcClint struct {
	conn *grpc.ClientConn
}

func NewRpcClint(addr string) (*RpcClint, error) {
	c := &RpcClint{}
	gateway_grpc.NewPluginServerClient(nil)
	var err error
	c.conn, err = grpc.Dial(addr)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c RpcClint) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
