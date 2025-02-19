package simulator

// Common SNMP OIDs for network switches
const (
	// System OIDs
	OIDSysDescr      = ".1.3.6.1.2.1.1.1.0"
	OIDSysObjectID   = ".1.3.6.1.2.1.1.2.0"
	OIDSysUpTime     = ".1.3.6.1.2.1.1.3.0"
	OIDSysContact    = ".1.3.6.1.2.1.1.4.0"
	OIDSysName       = ".1.3.6.1.2.1.1.5.0"
	OIDSysLocation   = ".1.3.6.1.2.1.1.6.0"
	OIDSysServices   = ".1.3.6.1.2.1.1.7.0"

	// Interface OIDs
	OIDIfNumber      = ".1.3.6.1.2.1.2.1.0"
	OIDIfTable       = ".1.3.6.1.2.1.2.2"
	OIDIfEntry       = ".1.3.6.1.2.1.2.2.1"
	OIDIfIndex       = ".1.3.6.1.2.1.2.2.1.1"
	OIDIfDescr       = ".1.3.6.1.2.1.2.2.1.2"
	OIDIfType        = ".1.3.6.1.2.1.2.2.1.3"
	OIDIfMtu         = ".1.3.6.1.2.1.2.2.1.4"
	OIDIfSpeed       = ".1.3.6.1.2.1.2.2.1.5"
	OIDIfPhysAddress = ".1.3.6.1.2.1.2.2.1.6"
	OIDIfAdminStatus = ".1.3.6.1.2.1.2.2.1.7"
	OIDIfOperStatus  = ".1.3.6.1.2.1.2.2.1.8"
	OIDIfInOctets    = ".1.3.6.1.2.1.2.2.1.10"
	OIDIfInErrors    = ".1.3.6.1.2.1.2.2.1.14"
	OIDIfOutOctets   = ".1.3.6.1.2.1.2.2.1.16"
	OIDIfOutErrors   = ".1.3.6.1.2.1.2.2.1.20"

	// Resource OIDs (enterprise-specific)
	OIDCPUUsage      = ".1.3.6.1.4.1.9.9.109.1.1.1.1.3.1" // Cisco-style CPU usage
	OIDMemoryUsage   = ".1.3.6.1.4.1.9.9.48.1.1.1.5.1"    // Cisco-style memory usage
	OIDTemperature   = ".1.3.6.1.4.1.9.9.13.1.3.1.3.1"    // Cisco-style temperature
)

// OIDMap maps OIDs to their descriptions
var OIDMap = map[string]string{
	OIDSysDescr:      "System Description",
	OIDSysObjectID:   "System Object ID",
	OIDSysUpTime:     "System Uptime",
	OIDSysContact:    "System Contact",
	OIDSysName:       "System Name",
	OIDSysLocation:   "System Location",
	OIDSysServices:   "System Services",
	OIDIfNumber:      "Number of Interfaces",
	OIDIfIndex:       "Interface Index",
	OIDIfDescr:       "Interface Description",
	OIDIfType:        "Interface Type",
	OIDIfMtu:         "Interface MTU",
	OIDIfSpeed:       "Interface Speed",
	OIDIfPhysAddress: "Interface Physical Address",
	OIDIfAdminStatus: "Interface Admin Status",
	OIDIfOperStatus:  "Interface Operational Status",
	OIDIfInOctets:    "Interface Input Octets",
	OIDIfInErrors:    "Interface Input Errors",
	OIDIfOutOctets:   "Interface Output Octets",
	OIDIfOutErrors:   "Interface Output Errors",
	OIDCPUUsage:      "CPU Usage",
	OIDMemoryUsage:   "Memory Usage",
	OIDTemperature:   "Temperature",
}
