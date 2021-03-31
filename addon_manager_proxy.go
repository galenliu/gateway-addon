package addon

import (
	"fmt"
	json "github.com/json-iterator/go"
	"log"
	"sync"
	"time"
)

const (
	Aid = "adapterId"
	Pid = "pluginId"
	Did = "deviceId"
)

type IProperty interface {
	UpdateValue(value interface{})
	GetDescription()string
}

type IDevice interface {
	GetProperty(name string) IProperty
	SetCredentials(username string, password string) error
	SetPin(pin interface{}) error
	GetDescription()string
	GetID()string
}

type IAdapter interface {
	//SendPairingPrompt(promt, url string, device *Device)
	StartPairing(timeout float64)
	CancelPairing()
	GetID() string
	GetName() string
	Unload()
	HandleDeviceSaved(device IDevice)
	HandleDeviceRemoved(device IDevice)
	GetDevice(deviceId string) IDevice
}

type AddonManagerProxy struct {
	ipcClient   *IpcClient
	adapters    map[string]IAdapter
	packageName string
	verbose     bool
	running     bool
}

var once sync.Once
var instance *AddonManagerProxy

func NewAddonManagerProxy(packageName string) *AddonManagerProxy {
	once.Do(
		func() {
			instance = &AddonManagerProxy{}
			instance.packageName = packageName
			instance.adapters = make(map[string]IAdapter, 10)
			instance.ipcClient = NewClient(packageName, instance.onMessage)
		},
	)
	return instance
}

func (proxy *AddonManagerProxy) handleAdapterAdded(adapter IAdapter) {
	proxy.adapters[adapter.GetID()] = adapter
	data := make(map[string]interface{})
	data[Aid] = adapter.GetID()
	data["name"] = adapter.GetName
	data["packageName"] = proxy.packageName
	proxy.send(AdapterAddedNotification, data)
}

func (proxy *AddonManagerProxy) HandleDeviceAdded(device *Device) {
	data := make(map[string]interface{})
	data[Aid] = device.AdapterId
	data["device"] = device
	proxy.send(DeviceAddedNotification, data)
}

func (proxy *AddonManagerProxy) HandleDeviceRemoved(device *Device) {

	data := make(map[string]interface{})
	data[Aid] = device.AdapterId
	data["device"] = device
	proxy.send(AdapterRemoveDeviceRequest, data)

}

func (proxy *AddonManagerProxy) getAdapter(adapterId string) (IAdapter, error) {
	adapter, ok := proxy.adapters[adapterId]
	if !ok {
		return nil, fmt.Errorf("adapter id(%s) invaild", adapterId)
	}
	return adapter, nil
}

func (proxy *AddonManagerProxy) onMessage(data []byte) {

	var messageType = json.Get(data, "messageType").ToInt()

	switch messageType {
	//卸载plugin
	case PluginUnloadRequest:
		data := make(map[string]interface{})
		proxy.send(PluginUnloadResponse, data)
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
		adapter.StartPairing(timeout)
		return

	case AdapterCancelPairingCommand:
		go adapter.CancelPairing()
		return

		//adapter unload request

	case AdapterUnloadRequest:
		adapter.Unload()
		unloadFunc := func(proxy *AddonManagerProxy, adapter IAdapter) {
			data := make(map[string]interface{})
			data[Aid] = adapter.GetID()
			proxy.send(AdapterUnloadResponse, data)
		}
		go unloadFunc(proxy, adapter)
		delete(proxy.adapters, adapter.GetID())
		break
	}
	var deviceId = json.Get(data, "data", "deviceId").ToString()
	device := adapter.GetDevice(deviceId)
	if device == nil {
		return
	}

	switch messageType {
	case AdapterCancelRemoveDeviceCommand:
		adapter := proxy.adapters[adapterId]
		log.Printf(adapter.GetID())

	case DeviceSavedNotification:

		adapter.HandleDeviceSaved(device)
		return

		//adapter remove devices request

	case AdapterRemoveDeviceRequest:
		adapter.HandleDeviceRemoved(device)

		//devices set properties command

	case DeviceSetPropertyCommand:
		propName := json.Get(data, "data", "propertyName").ToString()
		newValue := json.Get(data, "data", "propertyValue").GetInterface()
		prop := device.GetProperty(propName)
		if err != nil {
			log.Fatal(err.Error())
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
		data := make(map[string]interface{})
		data[Aid] = adapterId
		data[Did] = device.GetID()
		data["property"] =  prop.GetDescription()
		proxy.send(DevicePropertyChangedNotification,data,)


	case DeviceSetPinRequest:
		var pin PIN
		pin.Pattern = json.Get(data, "data", "pin", "pattern").GetInterface()
		pin.Required = json.Get(data, "data", "pin", "required").ToBool()
		messageId := json.Get(data, "data", "message_id").ToInt()
		if messageId == 0 {
			log.Fatal("DeviceSetPinRequest:  non  messageId")
		}

		handleFunc := func() {
			data := make(map[string]interface{})
			data[Aid] = adapterId
			data[Did] = deviceId
			data["devx"] = device
			data["messageId"] = messageId
			err := device.SetPin(pin)
			if err == nil {
				data["success"] = true
				proxy.send(DeviceSetPinResponse, data)

			} else {
				data["success"] = false
				proxy.send(DeviceSetPinResponse, data)
			}
		}
		go handleFunc()

	case DeviceSetCredentialsRequest:
		messageId := json.Get(data, "data", "messageId").ToInt()
		username := json.Get(data, "data", "username").ToString()
		password := json.Get(data, "data", "password").ToString()

		handleFunc := func() {
			err := device.SetCredentials(username, password)
			data := make(map[string]interface{})
			data[Aid] = adapterId
			data[Did] = deviceId
			data["messageId"] = messageId
			if err != nil {
				data["success"] = true
				proxy.send(DeviceSetCredentialsResponse, data)
				fmt.Printf(err.Error())
				return
			}
			data["success"] = false
			proxy.send(DeviceSetCredentialsResponse, data)
			return
		}
		go handleFunc()
		break
	}
}

func (proxy *AddonManagerProxy) sendConnectedStateNotification(device *Device, connected bool) {
	data := make(map[string]interface{})
	data[Aid] = device.AdapterId
	data[Did] = device.ID
	data["connected"] = connected
	proxy.send(DeviceConnectedStateNotification, data)
}



func (proxy *AddonManagerProxy) run() {
	proxy.ipcClient.Register()
}

func (proxy *AddonManagerProxy) handleDeviceRemoved(adapterId, devId string) {
	if proxy.verbose {
		fmt.Printf("addon manager handle devices added, deviceId:%v\n", devId)
	}
	data := make(map[string]interface{})
	data[Aid] = adapterId
	data[Did] = devId

	proxy.send(AdapterRemoveDeviceResponse, data)
}

func (proxy *AddonManagerProxy) send(messageType int, data map[string]interface{}) {
	data[Pid] = proxy.packageName
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

func (proxy *AddonManagerProxy) sendError(adapterID, message string) {

	data := make(map[string]interface{})
	data[Aid] = adapterID
	data["message"] = message
	proxy.send(PluginErrorNotification, data)

}

func (proxy *AddonManagerProxy) sendPairingPrompt(adapter *Adapter, promt, url string, device *Device) {
	data := make(map[string]interface{})
	data[Aid] = adapter.ID
	data["prompt"] = promt
	if device != nil {
		data[Did] = device.ID
	}
	if url != "" {
		data["url"] = url
	}
	proxy.send(AdapterPairingPromptNotification, data)
}

func (proxy *AddonManagerProxy) sendUnPairingPrompt(adapter *Adapter, prompt, url string, device *Device) {
	data := make(map[string]interface{})
	data[Aid] = adapter.ID
	data["prompt"] = prompt
	if device != nil {
		data[Did] = device.ID
	}
	if url != "" {
		data["url"] = url
	}
	proxy.send(AdapterUnpairingPromptNotification, data)
}

func (proxy *AddonManagerProxy) close() {
	proxy.ipcClient.close()
	proxy.running = false
}
