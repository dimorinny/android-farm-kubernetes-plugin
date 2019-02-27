package main

import (
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

type AndroidDevicesPlugin struct {
	devices []*pluginapi.Device
	server  *grpc.Server
}
