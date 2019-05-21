package main

import (
	"flag"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	socketcan "github.com/mpreu/k8s-device-plugin-socketcan/pkg"
)

func main() {
	flag.Parse()

	manager := dpm.NewManager(socketcan.Lister{})
	manager.Run()
}
