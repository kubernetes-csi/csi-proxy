package system

type GetBIOSSerialNumberRequest struct {
	// Intentionally empty
}

type GetBIOSSerialNumberResponse struct {
	// Serial number
	SerialNumber string
}

type StartServiceRequest struct {
	// Service name (as listed in System\CCS\Services keys)
	Name string
}

type StartServiceResponse struct {
	// Intentionally empty
}

type StopServiceRequest struct {
	// Service name (as listed in System\CCS\Services keys)
	Name string

	// Forces stopping of services that has dependent services
	Force bool
}

type StopServiceResponse struct {
	// Intentionally empty
}

type ServiceStatus uint32

// https://docs.microsoft.com/en-us/windows/win32/api/winsvc/ns-winsvc-service_status#members
const (
	SERVICE_STATUS_UNKNOWN ServiceStatus = iota
	SERVICE_STATUS_STOPPED
	SERVICE_STATUS_START_PENDING
	SERVICE_STATUS_STOP_PENDING
	SERVICE_STATUS_RUNNING
	SERVICE_STATUS_CONTINUE_PENDING
	SERVICE_STATUS_PAUSE_PENDING
	SERVICE_STATUS_PAUSED
)

type Startype uint32

// https://docs.microsoft.com/en-us/windows/win32/api/winsvc/nf-winsvc-changeserviceconfiga
const (
	START_TYPE_BOOT Startype = iota
	START_TYPE_SYSTEM
	START_TYPE_AUTOMATIC
	START_TYPE_MANUAL
	START_TYPE_DISABLED
)

type GetServiceRequest struct {
	// Service name (as listed in System\CCS\Services keys)
	Name string
}

type GetServiceResponse struct {
	// Service display name
	DisplayName string

	// Service start type
	// Used to control whether a service will start on boot, and if so on which
	// boot phase
	StartType Startype

	// Service status, e.g. stopped, running, paused
	Status ServiceStatus
}
