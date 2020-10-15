package main

import (
	"flag"

	diskapi "github.com/kubernetes-csi/csi-proxy/internal/os/disk"
	filesystemapi "github.com/kubernetes-csi/csi-proxy/internal/os/filesystem"
	smbapi "github.com/kubernetes-csi/csi-proxy/internal/os/smb"
	sysapi "github.com/kubernetes-csi/csi-proxy/internal/os/system"
	volumeapi "github.com/kubernetes-csi/csi-proxy/internal/os/volume"
	"github.com/kubernetes-csi/csi-proxy/internal/server"
	disksrv "github.com/kubernetes-csi/csi-proxy/internal/server/disk"
	filesystemsrv "github.com/kubernetes-csi/csi-proxy/internal/server/filesystem"
	smbsrv "github.com/kubernetes-csi/csi-proxy/internal/server/smb"
	syssrv "github.com/kubernetes-csi/csi-proxy/internal/server/system"
	srvtypes "github.com/kubernetes-csi/csi-proxy/internal/server/types"
	volumesrv "github.com/kubernetes-csi/csi-proxy/internal/server/volume"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"k8s.io/klog"
)

var (
	kubeletCSIPluginsPath = flag.String("kubelet-csi-plugins-path", `C:\var\lib\kubelet`, "Prefix path of the Kubelet plugin directory in the host file system")
	kubeletPodPath        = flag.String("kubelet-pod-path", `C:\var\lib\kubelet`, "Prefix path of the kubelet pod directory in the host file system")
	windowsSvc            = flag.Bool("windows-service", false, "Configure as a Windows Service")
	service               *handler
)

type handler struct {
	tosvc   chan bool
	fromsvc chan error
}

func main() {
	defer klog.Flush()
	klog.InitFlags(nil)

	flag.Parse()

	if *windowsSvc {
		if err := initService(); err != nil {
			panic(err)
		}
	}

	apiGroups, err := apiGroups()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(apiGroups...)

	klog.Info("Starting CSI-Proxy Server ...")
	klog.Infof("Version: %s", version)
	if err := s.Start(nil); err != nil {
		panic(err)
	}
}

// apiGroups returns the list of enabled API groups.
func apiGroups() ([]srvtypes.APIGroup, error) {
	fssrv, err := filesystemsrv.NewServer(*kubeletCSIPluginsPath, *kubeletPodPath, filesystemapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	volumesrv, err := volumesrv.NewServer(volumeapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	disksrv, err := disksrv.NewServer(diskapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	smbsrv, err := smbsrv.NewServer(smbapi.New(), fssrv)
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	syssrv, err := syssrv.NewServer(sysapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	return []srvtypes.APIGroup{
		fssrv,
		disksrv,
		volumesrv,
		smbsrv,
		syssrv,
	}, nil
}

// configure as a Windows service managed by Windows SCM
// code borrowed from
// https://github.com/kubernetes/kubernetes/blob/323f34858de18b862d43c40b2cced65ad8e24052/pkg/windows/service/service.go
func initService() error {
	h := &handler{
		tosvc:   make(chan bool),
		fromsvc: make(chan error),
	}

	service = h
	var err error
	go func() {
		err = svc.Run("csiproxy", h)
		h.fromsvc <- err
	}()

	// Wait for the first signal from the service handler.
	err = <-h.fromsvc
	if err != nil {
		return err
	}
	klog.Infof("Running as a Windows service.")
	return nil
}

func (h *handler) Execute(_ []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	s <- svc.Status{State: svc.StartPending, Accepts: 0}
	// Unblock initService()
	h.fromsvc <- nil

	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown | svc.Accepted(windows.SERVICE_ACCEPT_PARAMCHANGE)}
	klog.Infof("Windows Service initialized through SCM")
Loop:
	for {
		select {
		case <-h.tosvc:
			break Loop
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				s <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// todo: need to add a ctx to servers
				// from main and cancel it from here
				s <- svc.Status{State: svc.StopPending}
				break Loop
			}
		}
	}

	return false, 0
}
