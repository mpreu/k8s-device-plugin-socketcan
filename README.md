# SocketCAN Device Plugin for Kubernetes

## Intro
`SocketCAN` is a set of CAN drivers in the Linux kernel contributed by Volkswagen. Via different kernel modules support for real hardware CAN devices as well as virtual loopback CAN devices can be enabled. This Kubernetes device plugin is responsible to provide the corresponding devices inside of docker containers.

## Supported vcan modules
Right now, the device plugin supports the usage of the virtual can devices provided by the `vcan` kernel module.

## Build
To build the device plugin in a ready to use docker plugin run:

```
make build-docker
```

## Deployment
To deploy the device plugin as a DaemonSet in the Kubernetes cluster an example [yaml configuration](./deployments/socketcan-ds.yml) is provided.

To deploy the device plugin via a DaemonSet run:

```
make deploy-ds
```

To remove the device plugin run:

```
make remove-ds
```

## Usage of the device plugin
The device plugin is available through the namespace `socketcan.mpreu.de/vcan`. An example deployment is given under [example/consumer-vcan](./example/consumer-vcan/dc.yml).

## Dependencies on the Host
To be able to consume the `SocketCAN` devices the corresponding kernel modules must be available on the Kubernetes compute nodes. The kernel modules should be already provided by the distributions, one known exception at the time of writing are RHEL 7.0-7.2 (https://access.redhat.com/solutions/2259931)