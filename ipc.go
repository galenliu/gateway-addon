package addon

import (
	"fmt"
	"github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"net/url"
	"sync"
)

const (
	Disconnect = "Disconnect"
	Connected  = "Connected"
	Registered = "Registered"
)

type UserProfile struct {
	BaseDir        string `validate:"required" json:"base_dir"`
	DataDir        string `validate:"required" json:"data_dir"`
	AddonsDir      string `validate:"required" json:"addons_dir"`
	ConfigDir      string `validate:"required" json:"config_dir"`
	UploadDir      string `validate:"required" json:"upload_dir"`
	MediaDir       string `validate:"required" json:"media_dir"`
	LogDir         string `validate:"required" json:"log_dir"`
	GatewayVersion string
}

type Preferences struct {
	Language string `validate:"required" json:"language"`
	Units    Units  `validate:"required" json:"units"`
}

type Units struct {
	Temperature string `validate:"required" json:"temperature"`
}

const PORT = "9500"

type OnMessage func(data []byte)

//为Plugin提供和gateway Server进行消息的通信

type IpcClient struct {
	ws *websocket.Conn

	url         string
	preferences Preferences
	userProfile UserProfile

	writeCh   chan []byte
	readCh    chan []byte
	closeChan chan interface{}
	reConnect chan interface{}

	gatewayVersion string

	onMessage OnMessage
	mu        *sync.Mutex

	status   string
	pluginID string
	origin   string
	verbose  bool
}

//新建一个Client，注册消息Handler
func NewClient(PluginId string, handler OnMessage) *IpcClient {
	u := url.URL{Scheme: "ws", Host: "localhost:" + PORT, Path: "/"}
	client := &IpcClient{}
	client.pluginID = PluginId
	client.url = u.String()
	client.status = Disconnect
	client.mu = new(sync.Mutex)

	client.closeChan = make(chan interface{})
	client.reConnect = make(chan interface{})

	client.readCh = make(chan []byte)
	client.writeCh = make(chan []byte)

	client.onMessage = handler
	client.Register()
	go client.readLoop()
	return client
}

func (client *IpcClient) onData(data []byte) {

	fmt.Printf("read message : %s \t\n", string(data))

	if json.Get(data, "messageType").ToInt() == PluginRegisterResponse {
		client.preferences.Language = json.Get(data, "data", "preferences", "language").ToString()
		client.preferences.Units.Temperature = json.Get(data, "data", "preferences", "units", "temperature").ToString()
		client.userProfile.AddonsDir = json.Get(data, "data", "user_profile", "addons_dir").ToString()
		client.userProfile.BaseDir = json.Get(data, "data", "user_profile", "base_dir").ToString()
		client.userProfile.ConfigDir = json.Get(data, "data", "user_profile", "config_dir").ToString()
		client.userProfile.DataDir = json.Get(data, "data", "user_profile", "data_dir").ToString()
		client.userProfile.GatewayVersion = json.Get(data, "data", "user_profile", "gateway_version").ToString()
		client.userProfile.LogDir = json.Get(data, "data", "user_profile", "log_dir").ToString()
		client.userProfile.MediaDir = json.Get(data, "data", "user_profile", "media_dir").ToString()
		client.userProfile.UploadDir = json.Get(data, "data", "user_profile", "upload_dir").ToString()
		client.status = Registered
	} else {
		client.onMessage(data)
	}
}

//发送Message Struct
//func (client *IpcClient) sendMessage(data []byte) {
//	log.Printf("send message:  %s\r\n", string(data))
//	if client.ws != nil {
//		client.writeCh <- data
//		return
//	}
//}

func (client *IpcClient) sendMessage(data []byte) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.ws != nil && client.status == Registered {
		err := client.ws.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			fmt.Printf("ipc client write err")
			client.status = Disconnect
		}
		fmt.Printf("ipc client send message: %s \t\n", string(data))
	}
}

func (client *IpcClient) readMessage() {
	if client.ws != nil {
		_, data, err := client.ws.ReadMessage()
		if err != nil {
			fmt.Printf("read faild, websocket err", err.Error())
			client.status = Disconnect
		}
		client.onData(data)
	}
}

func (client *IpcClient) readLoop() {
	for {
		if client.status == Registered && client.ws != nil {
			client.readMessage()
		}
	}
}

//循环往readCh中读取 Message
//func (client *IpcClient) readLoop() {
//
//	for {
//		if client.ws != nil && client.status != Disconnect {
//			_, data, err := client.ws.readMessage()
//			if err == nil {
//				client.onData(data)
//			} else {
//				client.status = Disconnect
//				log.Printf("read loop err : %s", err.Error())
//			}
//		}
//	}
//}

//循环发送writeChan中的Message
//func (client *IpcClient) writeLoop() {
//	defer client.close()
//	for {
//		select {
//		case msg := <-client.writeCh:
//			if client.ws != nil && client.status != Disconnect {
//				err := client.ws.WriteMessage(websocket.BinaryMessage, msg)
//				if err != nil {
//					fmt.Printf("write loop err =%v", err)
//					client.status = Disconnect
//				}
//			}
//
//		case <-client.closeChan:
//			return
//		}
//
//	}
//	fmt.Sprint("exit writ loop")
//}

//client 失去连接后，重新连接
//func (client *IpcClient) runLoop() {
//	for {
//		if client.status == Disconnect {
//			err := client.dial()
//			if err != nil {
//				client.status = Disconnect
//				client.ws = nil
//				fmt.Printf("pluginID:%v, err:%v ,retry after %v second \r\n", client.pluginID, err, 5)
//				time.Sleep(5 * time.Second)
//				continue
//			}
//
//		}
//		client.status = Connected
//	}
//	fmt.Sprint("exit run loop")
//}

//func (client *IpcClient) close() {
//	if client.ws != nil {
//		err := client.ws.Close()
//		if err != nil {
//			fmt.Println("client close-----")
//		}
//	}
//	client.closeChan <- ""
//}

func (client *IpcClient) dial() error {

	var err error = nil
	client.ws, _, err = websocket.DefaultDialer.Dial(client.url, nil)
	if err != nil {
		fmt.Printf("dial err: %s \r\n", err.Error())
		return err
	}
	client.status = Connected
	return nil
}

func (client *IpcClient) Register() {

	if client.status == Disconnect {
		client.dial()
	}

	if client.status == Connected {
		message := struct {
			MessageType int         `json:"messageType"`
			Data        interface{} `json:"data"`
		}{
			MessageType: PluginRegisterRequest,
			Data: struct {
				PluginID string `json:"pluginId"`
			}{PluginID: client.pluginID},
		}

		d, _ := json.MarshalIndent(message, "", " ")
		client.ws.WriteMessage(websocket.BinaryMessage, d)
		_, data, err := client.ws.ReadMessage()
		if err != nil {
			fmt.Printf("read faild, websocket err", err.Error())
			client.status = Disconnect
		}
		client.onData(data)
	}

}

func (client *IpcClient) close() {
	if client.ws != nil {
		client.ws.Close()
	}
}
