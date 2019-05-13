module github.com/mpreu/k8s-device-plugin-socketcan

require (
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/uuid v1.1.1
	github.com/kubevirt/device-plugin-manager v1.9.3-0.20180705123155-a2dafa739e03
	github.com/kubevirt/kubernetes-device-plugins v0.0.1
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20180720170159-13995c7128cc
	golang.org/x/net v0.0.0-20190509222800-a4d6f7feada5 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/grpc v1.20.1 // indirect
	k8s.io/kubernetes v1.14.1
)

replace golang.org/x/net => github.com/golang/net v0.0.0-20190509222800-a4d6f7feada5
