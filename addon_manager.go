package addon

type AddonManager struct {
	pluginId string
	verbose  bool
	running  bool
}

func NewAddonManager(pluginId string) *AddonManager {
	manager := &AddonManager{
		pluginId: pluginId,
		verbose:  false,
		running:  false,
	}
	return manager
}

func (manager *AddonManager) GetPluginID() string {
	return manager.pluginId
}

func (manager *AddonManager) Running() bool {
	return manager.running
}

func (manager *AddonManager) Run() {
	manager.running = true
}

func (manager *AddonManager) Stop() {
	manager.running = false
}
func (manager *AddonManager) Start() {
	manager.running = true
}
