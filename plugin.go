package main

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
	"log"
	"net"
	"os"
	"path"
)

const (
	resourceName     = "android/device"
	serverSocketName = pluginapi.DevicePluginPath + "android.sock"
)

type AndroidDevicesPlugin struct {
	devicesListener DevicesListener

	server *grpc.Server

	socketPath    string
	kubeletSocket string
	resourceName  string
}

func NewAndroidDevicesPlugin(devicesListener DevicesListener) *AndroidDevicesPlugin {
	return &AndroidDevicesPlugin{
		devicesListener: devicesListener,

		socketPath:    serverSocketName,
		kubeletSocket: pluginapi.KubeletSocket,

		resourceName: resourceName,
	}
}

func (p *AndroidDevicesPlugin) GetDevicePluginOptions(
	ctx context.Context,
	in *pluginapi.Empty,
) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

func (p *AndroidDevicesPlugin) ListAndWatch(
	in *pluginapi.Empty,
	server pluginapi.DevicePlugin_ListAndWatchServer,
) error {
	panic("implement me")
}

func (p *AndroidDevicesPlugin) Allocate(
	ctx context.Context,
	in *pluginapi.AllocateRequest,
) (*pluginapi.AllocateResponse, error) {
	panic("implement me")
}

func (p *AndroidDevicesPlugin) PreStartContainer(
	ctx context.Context,
	in *pluginapi.PreStartContainerRequest,
) (*pluginapi.PreStartContainerResponse, error) {
	panic("implement me")
}

func (p *AndroidDevicesPlugin) Start() {
	err := p.cleanupPluginServerSocket()
	if err != nil {
		log.Fatal("Failed to cleanup plugin server socket: ", err)
	}

	log.Println("Starting serving device plugin on ", p.socketPath)
	err = p.startPluginGrpcServer()
	if err != nil {
		log.Fatal("Something went wrong during starting plugin grpc server: ", err)
	}
	log.Println("Device plugin's grpc server started")

	log.Println("Register plugin...")
	err = p.registerPlugin()
	if err != nil {
		_ = p.Stop()
		log.Fatal("Something went wrong during register plugin: ", err)
	}
	log.Println("Plugin registered")

	log.Println("Starting devices listener")
	p.devicesListener.Listen()
}

func (p *AndroidDevicesPlugin) Stop() error {
	if p.server == nil {
		return nil
	}

	p.server.Stop()
	p.server = nil

	return p.cleanupPluginServerSocket()
}

func (p *AndroidDevicesPlugin) startPluginGrpcServer() error {
	err := p.cleanupPluginServerSocket()
	if err != nil {
		return err
	}

	socket, err := net.Listen("unix", p.socketPath)
	if err != nil {
		return err
	}

	p.server = grpc.NewServer()
	pluginapi.RegisterDevicePluginServer(p.server, p)

	//noinspection GoUnhandledErrorResult
	go func() {
		err = p.server.Serve(socket)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = p.waitForPluginServer()
	if err != nil {
		return err
	}

	return nil
}

func (p *AndroidDevicesPlugin) registerPlugin() error {
	connectionEstablishContext, connectionEstablishContextCancel := grpcContext()
	defer connectionEstablishContextCancel()

	connection, err := dial(p.kubeletSocket, connectionEstablishContext)
	if err != nil {
		return err
	}

	//noinspection GoUnhandledErrorResult
	defer connection.Close()

	registrationClient := pluginapi.NewRegistrationClient(connection)
	registrationRequest := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(p.socketPath),
		ResourceName: p.resourceName,
	}

	registerContext, registerContextCancel := grpcContext()
	defer registerContextCancel()

	_, err = registrationClient.Register(registerContext, registrationRequest)
	if err != nil {
		return err
	}

	return nil
}

func (p *AndroidDevicesPlugin) waitForPluginServer() error {
	if p.server == nil {
		return errors.New("failed to wait grpc (not initialized)")
	}

	ctx, cancel := grpcContext()
	defer cancel()

	connection, err := dial(p.socketPath, ctx)
	if err != nil {
		return err
	}

	_ = connection.Close()

	return nil
}

func (p *AndroidDevicesPlugin) cleanupPluginServerSocket() error {
	if err := os.Remove(p.socketPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
