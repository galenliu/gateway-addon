package addon

import (
	"fmt"
	"sync"
)

type onPairingFunc func(timeout float64)
type OnCancelPairingFunc func()
type OnDeviceSavedFunc func(deivceId string, device *Device)
type OnSetCredentialsFunc func(deivceId, username, password string)
type OnSetPinFunc func(deivceId string, pin PIN) error

type AdapterProxy struct {
	*Adapter
	devicesProxy     map[string]*DeviceProxy
	OnPairing        onPairingFunc
	OnCancelPairing  OnCancelPairingFunc
	OnDeviceSaved    OnDeviceSavedFunc
	OnSetCredentials OnSetCredentialsFunc
	OnSetPin         OnSetPinFunc
	managerProxy     *AddonManagerProxy
	locker           *sync.Mutex
	cancelChan       chan struct{}
}

func NewAdapterProxy(adapterId, adapterName, packageName string) *AdapterProxy {
	manager := NewAddonManagerProxy(packageName)
	adp := &AdapterProxy{Adapter: NewAdapter(manager, adapterId, packageName, packageName)}
	adp.managerProxy = manager
	adp.managerProxy.handleAdapterAdded(adp)
	adp.locker = new(sync.Mutex)
	adp.devicesProxy = make(map[string]*DeviceProxy)
	adp.cancelChan = make(chan struct{})
	return adp
}

func (proxy *AdapterProxy) HandleDeviceAdded(dx *DeviceProxy) {
	proxy.devicesProxy[dx.ID] = dx
	proxy.Adapter.HandleDeviceAdded(dx.Device)
}

func (proxy *AdapterProxy) HandleDeviceRemoved(dx *DeviceProxy) {
	delete(proxy.devicesProxy, dx.ID)
	proxy.Adapter.HandleDeviceRemoved(dx.Device)
}

func (proxy *AdapterProxy) SendError(messsage string) {
	proxy.managerProxy.sendError(proxy.ID, messsage)
}

func (proxy *AdapterProxy) ConnectedNotify(device *Device, connected bool) {
	proxy.managerProxy.sendConnectedStateNotification(device, connected)
}

//向前端UI发送提示
func (proxy *AdapterProxy) SendPairingPrompt(promt, url string, device *Device) {
	proxy.managerProxy.sendPairingPrompt(proxy.Adapter, promt, url, device)
}

func (proxy *AdapterProxy) SendUnpairingPrompt(promt, url string, device *Device) {
	proxy.managerProxy.sendPairingPrompt(proxy.Adapter, promt, url, device)
}

func (proxy *AdapterProxy) handleDeviceSaved(devId string, dev *Device) {
	fmt.Print("on devices saved on the gateway")
	if proxy.OnDeviceSaved != nil {
		proxy.OnDeviceSaved(devId, dev)
	}
}

func (proxy *AdapterProxy) startPairing(timeout float64) {

	if proxy.IsPairing {
		fmt.Print("proxy is pairinged")
		return
	}
	if proxy.OnPairing != nil {
		proxy.OnPairing(timeout)
	}

}

func (proxy *AdapterProxy) setCredentials(deivceId, username, password string) {
	if proxy.OnSetCredentials != nil {
		go proxy.OnSetCredentials(deivceId, username, password)
	}
}

func (proxy *AdapterProxy) setPin(deivceId string, pin PIN) error {
	if proxy.OnSetPin != nil {
		return proxy.OnSetPin(deivceId, pin)
	}
	return nil
}

func (proxy *AdapterProxy) cancelPairing() {
	if !proxy.IsPairing {
		return
	}
	proxy.IsPairing = false
	if proxy.OnCancelPairing != nil {
		go proxy.OnCancelPairing()
	}
}

func (proxy *AdapterProxy) getDeviceProxy(id string) *DeviceProxy {
	return proxy.devicesProxy[id]
}

func (proxy *AdapterProxy) Unload() {
	fmt.Printf("proxy unload, AdapterId:%v", proxy.ID)
}

func (proxy *AdapterProxy) CloseProxy() {
	fmt.Print("do some thing while proxy close")
	proxy.managerProxy.close()
}

func (proxy *AdapterProxy) ProxyRunning() bool {
	return proxy.managerProxy.running
}
