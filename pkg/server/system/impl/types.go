package impl

type GetBIOSSerialNumberRequest struct {
}

type GetBIOSSerialNumberResponse struct {
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

const (
	SERVICE_STATUS_UNKNOWN          = 0
	SERVICE_STATUS_STOPPED          = 1
	SERVICE_STATUS_START_PENDING    = 2
	SERVICE_STATUS_STOP_PENDING     = 3
	SERVICE_STATUS_RUNNING          = 4
	SERVICE_STATUS_CONTINUE_PENDING = 5
	SERVICE_STATUS_PAUSE_PENDING    = 6
	SERVICE_STATUS_PAUSED           = 7
)

type Startype uint32

const (
	START_TYPE_BOOT      = 0
	START_TYPE_SYSTEM    = 1
	START_TYPE_AUTOMATIC = 2
	START_TYPE_MANUAL    = 3
	START_TYPE_DISABLED  = 4
)

type GetServiceRequest struct {
	// Service name (as listed in System\CCS\Services keys)
	Name string
}

type GetServiceResponse struct {
	// Service display name
	DisplayName string

	// Service start type.
	// Used to control whether a service will start on boot, and if so on which
	// boot phase.
	StartType Startype

	// Service status, e.g. stopped, running, paused
	Status ServiceStatus
}
