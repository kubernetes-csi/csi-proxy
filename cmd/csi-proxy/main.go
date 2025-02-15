package main

import (
	"flag"

	diskapi "github.com/kubernetes-csi/csi-proxy/pkg/os/disk"
	filesystemapi "github.com/kubernetes-csi/csi-proxy/pkg/os/filesystem"
	iscsiapi "github.com/kubernetes-csi/csi-proxy/pkg/os/iscsi"
	smbapi "github.com/kubernetes-csi/csi-proxy/pkg/os/smb"
	sysapi "github.com/kubernetes-csi/csi-proxy/pkg/os/system"
	volumeapi "github.com/kubernetes-csi/csi-proxy/pkg/os/volume"
	"github.com/kubernetes-csi/csi-proxy/pkg/server"
	disksrv "github.com/kubernetes-csi/csi-proxy/pkg/server/disk"
	filesystemsrv "github.com/kubernetes-csi/csi-proxy/pkg/server/filesystem"
	iscsisrv "github.com/kubernetes-csi/csi-proxy/pkg/server/iscsi"
	"github.com/kubernetes-csi/csi-proxy/pkg/server/metrics"
	smbsrv "github.com/kubernetes-csi/csi-proxy/pkg/server/smb"
	syssrv "github.com/kubernetes-csi/csi-proxy/pkg/server/system"
	srvtypes "github.com/kubernetes-csi/csi-proxy/pkg/server/types"
	volumesrv "github.com/kubernetes-csi/csi-proxy/pkg/server/volume"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"k8s.io/klog/v2"
)

type workingDirFlags []string

func (i *workingDirFlags) String() string {
	return "Not implemented"
}

func (i *workingDirFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	kubeletPath     = flag.String("kubelet-path", `C:\var\lib\kubelet`, "Prefix path of the kubelet directory in the host file system")
	windowsSvc      = flag.Bool("windows-service", false, "Configure as a Windows Service")
	requirePrivacy  = flag.Bool("require-privacy", true, "If true, New-SmbGlobalMapping will be called with -RequirePrivacy $true")
	metricsBindAddr = flag.String("metrics-bind-address", "", "The address the metric endpoint binds to. Defaults to empty in which case metrics are disabled")
	service         *handler
	workingDirs     workingDirFlags
)

type handler struct {
	tosvc   chan bool
	fromsvc chan error
}

func init() {
	flag.Var(&workingDirs, "working-dir", "Prefix path of the csi-proxy working directory in the host file system")
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

	klog.Info("Starting CSI-Proxy Server ...")
	klog.Infof("Version: %s", version)
	apiGroups, err := apiGroups()
	if err != nil {
		panic(err)
	}

	enableMetrics := *metricsBindAddr != ""
	if enableMetrics {
		err := metrics.SetupMetricsServer(*metricsBindAddr)
		if err != nil {
			panic(err)
		}
	}
	s := server.NewServer(enableMetrics, apiGroups...)

	if err := s.Start(nil); err != nil {
		panic(err)
	}
}

// apiGroups returns the list of enabled API groups.
func apiGroups() ([]srvtypes.APIGroup, error) {
	workingDirs = append(workingDirs, *kubeletPath)
	fssrv, err := filesystemsrv.NewServer(workingDirs, filesystemapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}
	klog.Infof("Working directories: %v", fssrv.GetWorkingDirs())
	klog.Infof("Require privacy: %t", *requirePrivacy)

	volumesrv, err := volumesrv.NewServer(volumeapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	disksrv, err := disksrv.NewServer(diskapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	smbsrv, err := smbsrv.NewServer(smbapi.New(*requirePrivacy), fssrv)
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	syssrv, err := syssrv.NewServer(sysapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	iscsisrv, err := iscsisrv.NewServer(iscsiapi.New())
	if err != nil {
		return []srvtypes.APIGroup{}, err
	}

	return []srvtypes.APIGroup{
		fssrv,
		disksrv,
		volumesrv,
		smbsrv,
		syssrv,
		iscsisrv,
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
