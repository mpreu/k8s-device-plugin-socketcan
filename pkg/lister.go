package socketcan

import (
	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

const (
	resourceNamespace = "socketcan.mpreu.de"
)

// Lister implements the Lister interface from the
// device plugin manager
type Lister struct{}

// GetResourceNamespace declares the resource namespace in the FQDN format
func (s Lister) GetResourceNamespace() string {
	return resourceNamespace
}

// Discover which device plugins exist in the given resource namespace
func (s Lister) Discover(pluginListCh chan dpm.PluginNameList) {
	plugins := dpm.PluginNameList{"vcan"}

	pluginListCh <- plugins
}

// NewPlugin is called by the device plugin manager to create a device plugin
func (s Lister) NewPlugin(resourceName string) dpm.PluginInterface {
	glog.V(3).Infof("Creating device plugin %s", resourceName)
	return &DevicePlugin{
		allocationCh: make(chan *Allocation),
	}
}
