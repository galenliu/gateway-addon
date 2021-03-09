package addon

import (
	"fmt"
	json "github.com/json-iterator/go"
	"log"
	"sync"
	"time"
)

type AddonManagerProxy struct {
	*AddonManager
	ipcClient *IpcClient
	adapters  map[string]*AdapterProxy
}

var once sync.Once
var instance *AddonManagerProxy

func NewAddonManagerProxy(packageName string) *AddonManagerProxy {
	once.Do(
		func() {
			instance = &AddonManagerProxy{}
			instance.AddonManager = NewAddonManager(packageName)
			instance.adapters = make(map[string]*AdapterProxy, 10)
			instance.ipcClient = NewClient(packageName, instance.OnMessage)
			instance.Run()
		},
	)
	return instance
}

func (proxy *AddonManagerProxy) handleAdapterAdded(adapter *AdapterProxy) {

	proxy.adapters[adapter.ID] = adapter

	message := struct {
		PluginId    string `json:"pluginId"`
		Name        string `json:"name"`
		PackageName string `json:"packageName"`
		AdapterId   string `json:"adapterId"`
	}{
		PluginId:    proxy.pluginId,
		Name:        adapter.Name,
		AdapterId:   adapter.ID,
		PackageName: proxy.pluginId,
	}
	proxy.send(AdapterAddedNotification, message)
}

func (proxy *AddonManagerProxy) HandleDeviceAdded(device *Device) {

	message := struct {
		PluginId  string  `json:"pluginId"`
		AdapterId string  `json:"adapterId"`
		Device    *Device `json:"device"`
	}{
		PluginId:  proxy.pluginId,
		AdapterId: device.AdapterId,
		Device:    device,
	}
	proxy.send(DeviceAddedNotification, message)

}

func (proxy *AddonManagerProxy) HandleDeviceRemoved(device *Device) {

	message := struct {
		PluginId  string  `json:"pluginId"`
		AdapterId string  `json:"adapterId"`
		Device    *Device `json:"device"`
	}{
		PluginId:  proxy.pluginId,
		AdapterId: device.AdapterId,
		Device:    device,
	}
	proxy.send(AdapterRemoveDeviceRequest, message)

}

func (proxy *AddonManagerProxy) getAdapter(adapterId string) (*AdapterProxy, error) {
	adapter, ok := proxy.adapters[adapterId]
	if !ok {
		return nil, fmt.Errorf("adapter id(%s) invaild", adapterId)
	}
	return adapter, nil
}

func (proxy *AddonManagerProxy) OnMessage(data []byte) {

	var messageType = json.Get(data, "messageType").ToInt()

	switch messageType {
	//卸载plugin
	case PluginUnloadRequest:
		proxy.send(PluginUnloadResponse, struct {
			PluginId string
		}{PluginId: proxy.pluginId})
		proxy.running = false
		var closeFun = func() {
			time.AfterFunc(500*time.Millisecond, func() { proxy.close() })
		}
		go closeFun()
		break
	}

	var adapterId = json.Get(data, "data", "adapterId").ToString()
	adapter, err := proxy.getAdapter(adapterId)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	switch messageType {
	//adapter pairing command
	case AdapterStartPairingCommand:
		timeout := json.Get(data, "data", "timeout").ToFloat64()
		adapter.Pairing(timeout)
		return

	case AdapterCancelPairingCommand:
		go adapter.cancelPairing()
		return

		//adapter unload request

	case AdapterUnloadRequest:
		adapter.Unload()
		unloadFunc := func(proxy *AddonManagerProxy, adapter *AdapterProxy) {
			proxy.send(AdapterUnloadResponse, struct {
				AdapterId string `json:"AdapterId"`
			}{AdapterId: adapter.ID})
		}
		go unloadFunc(proxy, adapter)
		delete(proxy.adapters, adapter.ID)
		break
	}

	var deviceId = json.Get(data, "data", "deviceId").ToString()
	device, err := adapter.FindDevice(deviceId)
	if err != nil {
		log.Println(err.Error())
		return
	}

	switch messageType {
	case AdapterCancelRemoveDeviceCommand:
		adapter := proxy.adapters[adapterId]
		log.Printf(adapter.ID)

	case DeviceSavedNotification:
		adapter := proxy.adapters[adapterId]
		log.Fatal(adapter.ID)
		return

		//adapter remove devices request

	case AdapterRemoveDeviceRequest:
		//go adapter.removeDevice(deviceId)

		//devices set properties command

	case DeviceSetPropertyCommand:
		propName := json.Get(data, "data", "propertyName").ToString()
		newValue := json.Get(data, "data", "propertyValue").GetInterface()
		prop, err := device.FindProperty(propName)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		propChanged := func(newValue interface{}) error {
			prop.UpdateValue(newValue)

			return nil
		}
		e := propChanged(newValue)
		if e != nil {
			log.Printf(e.Error())
			return
		}
		proxy.sendPropertyChangedNotification(prop)

		//devices pin

	case DeviceSetPinRequest:
		pin := json.Get(data, "data", "pin").GetInterface()
		if pin == nil {
			log.Fatal("DeviceSetPinRequest: not find pin form message")
			return
		}
		messageId := json.Get(data, "data", "message_id").ToInt()
		if messageId == 0 {
			log.Fatal("DeviceSetPinRequest:  non  messageId")
		}
		handleFunc := func() {
			err := device.SetPin(pin)
			if err == nil {
				proxy.send(DeviceSetPinResponse, struct {
					PluginId  string  `json:"pluginId"`
					AdapterId string  `json:"adapterId"`
					MessageId int     `json:"messageId"`
					DeviceId  string  `json:"deviceId"`
					Device    *Device `json:"device"`
					Success   bool    `json:"success"`
				}{
					PluginId:  proxy.pluginId,
					AdapterId: adapterId,
					MessageId: messageId,
					DeviceId:  deviceId,
					Device:    device,
					Success:   true,
				})

			} else {
				proxy.send(DeviceSetPinResponse, struct {
					PluginId  string  `json:"pluginId"`
					AdapterId string  `json:"adapterId"`
					MessageId int     `json:"messageId"`
					DeviceId  string  `json:"deviceId"`
					Device    *Device `json:"device"`
					Success   bool    `json:"success"`
				}{
					PluginId:  proxy.pluginId,
					AdapterId: adapterId,
					MessageId: messageId,
					DeviceId:  deviceId,
					Device:    device,
					Success:   false,
				})
			}
		}
		go handleFunc()

	case DeviceSetCredentialsRequest:
		messageId := json.Get(data, "data", "messageId").ToInt()
		username := json.Get(data, "data", "username").ToString()
		password := json.Get(data, "data", "password").ToString()

		handleFunc := func() {
			err := device.SetCredentials(username, password)
			if err != nil {
				fmt.Printf(err.Error())
				proxy.send(DeviceSetCredentialsResponse, struct {
					PluginId  string
					AdapterId string
					MessageId int
					DeviceId  string
					Device    *Device
					Success   bool
				}{
					PluginId:  proxy.pluginId,
					AdapterId: adapterId,
					MessageId: messageId,
					DeviceId:  deviceId,
					Device:    device,
					Success:   false,
				})
				return
			}
			proxy.send(DeviceSetCredentialsResponse, struct {
				PluginId  string
				AdapterId string
				MessageId int
				DeviceId  string
				Device    *Device
				Success   bool
			}{
				PluginId:  proxy.pluginId,
				AdapterId: adapterId,
				MessageId: messageId,
				DeviceId:  deviceId,
				Device:    device,
				Success:   true,
			})
		}
		go handleFunc()
		break
	}

}

func (proxy *AddonManagerProxy) sendPropertyChangedNotification(p *Property) {
	data := struct {
		PluginId  string    `json:"pluginId"`
		AdapterId string    `json:"adapterId"`
		DeviceId  string    `json:"deviceId"`
		Property  *Property `json:"property"`
	}{
		PluginId:  proxy.pluginId,
		AdapterId: proxy.pluginId,
		DeviceId:  p.DeviceId,
		Property:  p,
	}
	proxy.send(DevicePropertyChangedNotification, data)
}

func (proxy *AddonManagerProxy) run() {
	proxy.ipcClient.Register()
}

func (proxy *AddonManagerProxy) handleDeviceRemoved(adapterId, devId string) {
	if proxy.verbose {
		fmt.Printf("addon manager handle devices added, deviceId:%v\n", devId)
	}
	message := struct {
		PluginId  string `json:"pluginId"`
		AdapterId string `json:"AdapterId"`
	}{
		PluginId:  proxy.pluginId,
		AdapterId: adapterId,
	}
	proxy.send(AdapterRemoveDeviceResponse, message)
}

func (proxy *AddonManagerProxy) send(messageType int, data interface{}) {

	var message = struct {
		MessageType int         `json:"messageType"`
		Data        interface{} `json:"data"`
	}{MessageType: messageType, Data: data}

	d, er := json.MarshalIndent(message, "", " ")
	if er != nil {
		log.Fatal(er)
		return
	}
	proxy.ipcClient.sendMessage(d)
}

func (proxy *AddonManagerProxy) close() {
	proxy.ipcClient.close()
	proxy.running = false
}
