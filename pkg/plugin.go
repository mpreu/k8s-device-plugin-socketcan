package socketcan

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/kubevirt/kubernetes-device-plugins/pkg/dockerutils"
	"github.com/mpreu/k8s-device-plugin-socketcan/pkg/vcan"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

const (
	fakeDeviceHostPath = "/var/run/device-plugin-socketcan-fakedev"
	nicPoolSize        = 100
	nicCreationRetries = 60
)

// DevicePlugin is the represents the device plugin and implements
// the Kuberentes device plugin interface
type DevicePlugin struct {
	allocationCh chan *Allocation
}

type Allocation struct {
	DeviceID            string
	DeviceContainerPath string
}

// Implementation of the Kubernetes device plugin interface

// GetDevicePluginOptions return options for the device plugin.
// Implementation of the 'DevicePluginServer' interface.
func (DevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return nil, nil
}

// ListAndWatch communicates changes of device states and returns a
// new device list. Implementation of the 'DevicePluginServer' interface.
func (p *DevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	devices := generateVCANDevices()
	s.Send(&pluginapi.ListAndWatchResponse{Devices: devices})

	for {
		time.Sleep(10 * time.Second)
	}
}

// Allocate is resposible to make the device available during the
// container creation process. Implementation of the 'DevicePluginServer' interface.
func (p *DevicePlugin) Allocate(c context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {

	var response pluginapi.AllocateResponse

	for _, req := range r.GetContainerRequests() {
		var devices []*pluginapi.DeviceSpec
		for _, deviceID := range req.GetDevicesIDs() {
			dev := &pluginapi.DeviceSpec{}
			dev.HostPath = fakeDeviceHostPath
			dev.ContainerPath = getDeviceContainerPath(deviceID)
			dev.Permissions = "r"

			devices = append(devices, dev)

			p.allocationCh <- &Allocation{
				DeviceID:            deviceID,
				DeviceContainerPath: dev.ContainerPath,
			}

		}
		response.ContainerResponses = append(response.ContainerResponses, &pluginapi.ContainerAllocateResponse{
			Devices: devices,
		})
	}

	return &response, nil
}

// PreStartContainer is called during registration phase of a container.
// Implementation of the 'DevicePluginServer' interface.
func (DevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return nil, nil
}

// Implement PluginInterfaceStart from kubevirts device-plugin-manager

func (p *DevicePlugin) Start() error {
	err := createFakeDevice()
	if err != nil {
		glog.Exitf("Failed to create fake device: %s", err)
	}

	go p.createVCANNic()

	return nil
}

// Additional functions
func (p *DevicePlugin) createVCANNic() {
	client, err := dockerutils.NewClient()

	if err != nil {
		glog.V(3).Info("Failed to connect to Docker")
		panic(err)
	}

	go func() {
		alloc := <-p.allocationCh
		glog.V(3).Infof("New allocation request: %v", alloc)

		for i := 0; i < nicCreationRetries; i++ {
			if created := func() bool {
				containerID, err := client.GetContainerIDByMountedDevice(alloc.DeviceContainerPath)

				if err != nil {
					glog.V(3).Infof("Container was not found, due to: %s", err.Error())
					return false
				}

				containerPID, err := client.GetPidByContainerID(containerID)
				if err != nil {
					glog.V(3).Infof("Failed to obtain container's pid, due to: %s", err.Error())
					return false
				}

				err = p.createVCANNicInPod(containerID, containerPID)
				if err == nil {
					glog.V(3).Info("Successfully create vcan interface")
					return true
				}

				glog.V(3).Infof("Pod attachment failed with: %s", err.Error())
				return false
			}(); created {
				break
			}
			time.Sleep(time.Duration(i) * time.Second)
		}
	}()
}

func (p *DevicePlugin) createVCANNicInPod(containerID string, containerPID int) error {

	link := &vcan.Vcan{
		LinkAttrs: netlink.LinkAttrs{
			Name: "vcan0",
		},
	}

	// Store current namespace
	originalNS, err := netns.Get()
	if err != nil {
		return err
	}

	// Get namespace of the pod
	ns, err := netns.GetFromPid(containerPID)
	if err != nil {
		return err
	}

	// Set to pod namespace so that interface values could be set
	err = netns.Set(ns)
	if err != nil {
		return err
	}

	// Set back to the original namespace before we leave this function
	defer func() {
		setErr := netns.Set(originalNS)
		if setErr != nil {
			// if we cannot go back the the original namespace
			// the plugin cannot be used anymore and a restart is needed
			panic(setErr)
		}
	}()

	// Add new link
	err = netlink.LinkAdd(link)
	if err != nil {
		return err
	}

	// Set interface up
	err = netlink.LinkSetUp(link)
	if err != nil {
		netlink.LinkDel(link)
		return err
	}

	return nil
}

func getDeviceContainerPath(nic string) string {
	return fmt.Sprintf("/tmp/device-plugin-socketcan/%s", nic)
}

func createFakeDevice() error {
	_, statErr := os.Stat(fakeDeviceHostPath)
	if statErr == nil {
		glog.V(3).Info("Fake block device already exists")
		return nil
	} else if os.IsNotExist(statErr) {
		glog.V(3).Info("Creating fake block device")
		cmd := exec.Command("mknod", fakeDeviceHostPath, "b", "1", "1")
		err := cmd.Run()
		return err
	} else {
		panic(statErr)
	}
}

func generateVCANDevices() []*pluginapi.Device {
	devices := make([]*pluginapi.Device, nicPoolSize)

	for i := range devices {
		devices[i] = &pluginapi.Device{
			ID:     fmt.Sprintf("vcan-%s", getUUID()),
			Health: pluginapi.Healthy,
		}
	}
	return devices
}

func getUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
