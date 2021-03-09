package addon

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type onPairingFunc func(ctx context.Context)
type OnCancelPairingFunc func()

type AdapterProxy struct {
	OnPairing       onPairingFunc
	OnCancelPairing OnCancelPairingFunc
	*Adapter
	addonManager *AddonManagerProxy
	locker       *sync.Mutex
	cancelChan   chan struct{}
}

func NewAdapterProxy(adapterId, adapterName, packageName string) *AdapterProxy {
	manager := NewAddonManagerProxy(packageName)
	adp := &AdapterProxy{Adapter: NewAdapter(manager, adapterId, packageName, packageName)}
	adp.addonManager = manager
	adp.addonManager.handleAdapterAdded(adp)
	adp.locker = new(sync.Mutex)
	adp.cancelChan = make(chan struct{})
	return adp
}

func (adapter *AdapterProxy) AddDevice(device *Device) {
	device.AdapterId = adapter.ID
	adapter.HandleDeviceAdded(device)
}

func (adapter *AdapterProxy) Pairing(timeout float64) {

	if adapter.IsPairing {
		fmt.Print("adapter is pairinged")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	adapter.IsPairing = true

	if adapter.OnPairing != nil {
		go adapter.OnPairing(ctx)
	}
	for {
		select {
		case <-ctx.Done():
			cancel()
			adapter.IsPairing = false
			return
		case <-adapter.cancelChan:
			if adapter.OnCancelPairing != nil {
				adapter.OnCancelPairing()
			}
			cancel()
			adapter.IsPairing = false
			return
		}
	}

}

func (adapter *AdapterProxy) cancelPairing() {
	if !adapter.IsPairing {
		fmt.Print("adapter pairing is canceled \t\n")
		return
	}
	adapter.cancelChan <- struct{}{}
}

func (adapter *AdapterProxy) HandleDeviceSaved(devId string, dev interface{}) {
	fmt.Print("on devices saved on the gateway")
}

func (adapter *AdapterProxy) Run() {
}

func (adapter *AdapterProxy) Unload() {
	fmt.Printf("adapter unload, AdapterId:%v", adapter.ID)
}

func (adapter *AdapterProxy) getID() string {
	return adapter.ID
}

func (adapter *AdapterProxy) getPackageName() string {
	return adapter.PackageName
}

func (adapter *AdapterProxy) CloseProxy() {
	fmt.Print("do some thing while adapter close")
}

func (adapter *AdapterProxy) ProxyRunning() bool {
	//return adapter.manager.Running()
	return true
}

func (adapter *AdapterProxy) HandlerDevicePropertyChanged(property *Property) {
	adapter.addonManager.sendPropertyChangedNotification(property)
}
