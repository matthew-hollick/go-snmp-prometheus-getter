package simulator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
)

// SwitchSimulator simulates a network switch for SNMP testing
type SwitchSimulator struct {
	community string
	port     uint16
	listener *net.UDPConn
	running  bool
	mu       sync.RWMutex
	metrics  *SwitchMetrics
}

// SwitchMetrics represents the simulated switch metrics
type SwitchMetrics struct {
	// System metrics
	Uptime       uint32
	Hostname     string
	Description  string
	Location     string
	Contact      string
	LastBootTime time.Time

	// Interface metrics
	Interfaces []InterfaceMetrics

	// Resource metrics
	CPUUsage    float64
	MemoryUsage float64
	Temperature float64
}

// InterfaceMetrics represents metrics for a single interface
type InterfaceMetrics struct {
	Index       int
	Name        string
	AdminStatus int
	OperStatus  int
	Speed       uint64
	InOctets    uint64
	OutOctets   uint64
	InErrors    uint64
	OutErrors   uint64
}

const (
	OIDSysDescr     = ".1.3.6.1.2.1.1.1.0"
	OIDSysObjectID  = ".1.3.6.1.2.1.1.2.0"
	OIDSysUpTime    = ".1.3.6.1.2.1.1.3.0"
	OIDSysContact   = ".1.3.6.1.2.1.1.4.0"
	OIDSysName      = ".1.3.6.1.2.1.1.5.0"
	OIDSysLocation  = ".1.3.6.1.2.1.1.6.0"
	OIDSysServices  = ".1.3.6.1.2.1.1.7.0"
	OIDIfNumber     = ".1.3.6.1.2.1.2.1.0"
	OIDIfEntry      = ".1.3.6.1.2.1.2.2.1"
	OIDIfIndex      = OIDIfEntry + ".1"
	OIDIfDescr      = OIDIfEntry + ".2"
	OIDIfType       = OIDIfEntry + ".3"
	OIDIfSpeed      = OIDIfEntry + ".5"
	OIDIfAdminStatus = OIDIfEntry + ".7"
	OIDIfOperStatus  = OIDIfEntry + ".8"
	OIDIfInOctets    = OIDIfEntry + ".10"
	OIDIfOutOctets   = OIDIfEntry + ".16"
	OIDIfInErrors    = OIDIfEntry + ".14"
	OIDIfOutErrors   = OIDIfEntry + ".20"
	OIDCPUUsage      = ".1.3.6.1.4.1.9.9.109.1.1.1.1.6.1"
	OIDMemoryUsage   = ".1.3.6.1.4.1.9.9.109.1.1.1.1.7.1"
	OIDTemperature   = ".1.3.6.1.4.1.9.9.109.1.1.1.1.8.1"
)

// NewSwitchSimulator creates a new switch simulator
func NewSwitchSimulator(community string, port uint16) *SwitchSimulator {
	return &SwitchSimulator{
		community: community,
		port:     port,
		metrics:  newDefaultMetrics(),
	}
}

// newDefaultMetrics creates default switch metrics
func newDefaultMetrics() *SwitchMetrics {
	interfaces := make([]InterfaceMetrics, 24) // 24-port switch
	for i := range interfaces {
		interfaces[i] = InterfaceMetrics{
			Index:       i + 1,
			Name:        fmt.Sprintf("GigabitEthernet1/%d", i+1),
			AdminStatus: 1, // up
			OperStatus:  1, // up
			Speed:       1000000000, // 1Gbps
		}
	}

	return &SwitchMetrics{
		Uptime:      0,
		Hostname:    "switch01.network.hedgehog.internal",
		Description: "Network Switch - Hedgehog Analytics",
		Location:    "london",
		Contact:     "hedgehog_admin",
		Interfaces:  interfaces,
		CPUUsage:    0,
		MemoryUsage: 0,
		Temperature: 25.0,
	}
}

// Start starts the SNMP simulator
func (s *SwitchSimulator) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("simulator already running")
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: int(s.port),
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("starting UDP listener: %w", err)
	}

	s.listener = listener
	s.running = true
	s.metrics.LastBootTime = time.Now()

	go s.serve()
	go s.updateMetrics()

	return nil
}

// Stop stops the SNMP simulator
func (s *SwitchSimulator) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("closing listener: %w", err)
	}

	s.running = false
	return nil
}

// serve handles incoming SNMP requests
func (s *SwitchSimulator) serve() {
	buffer := make([]byte, 2048)
	for {
		s.mu.RLock()
		if !s.running {
			s.mu.RUnlock()
			return
		}
		s.mu.RUnlock()

		n, remoteAddr, err := s.listener.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		msg := gosnmp.SnmpPacket{}
		if err := msg.UnmarshalBinary(buffer[:n]); err != nil {
			continue
		}

		if msg.Community != s.community {
			continue
		}

		response := s.handleMessage(&msg)
		if response == nil {
			continue
		}

		responseBytes, err := response.MarshalBinary()
		if err != nil {
			continue
		}

		s.listener.WriteToUDP(responseBytes, remoteAddr)
	}
}

// updateMetrics periodically updates simulated metrics
func (s *SwitchSimulator) updateMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		s.mu.RLock()
		if !s.running {
			s.mu.RUnlock()
			return
		}
		s.mu.RUnlock()

		s.mu.Lock()
		// Update uptime
		s.metrics.Uptime = uint32(time.Since(s.metrics.LastBootTime).Seconds() * 100)

		// Update interface metrics
		for i := range s.metrics.Interfaces {
			if s.metrics.Interfaces[i].AdminStatus == 1 && s.metrics.Interfaces[i].OperStatus == 1 {
				s.metrics.Interfaces[i].InOctets += uint64(rand.Int63n(1000000))
				s.metrics.Interfaces[i].OutOctets += uint64(rand.Int63n(1000000))
				if rand.Float64() < 0.01 { // 1% chance of errors
					s.metrics.Interfaces[i].InErrors++
					s.metrics.Interfaces[i].OutErrors++
				}
			}
		}

		// Update resource metrics
		s.metrics.CPUUsage = 20 + 10*rand.Float64()    // 20-30%
		s.metrics.MemoryUsage = 40 + 20*rand.Float64() // 40-60%
		s.metrics.Temperature = 25 + 5*rand.Float64()   // 25-30°C
		s.mu.Unlock()

		<-ticker.C
	}
}

// handleMessage processes an SNMP message and returns a response
func (s *SwitchSimulator) handleMessage(msg *gosnmp.SnmpPacket) *gosnmp.SnmpPacket {
	s.mu.RLock()
	defer s.mu.RUnlock()

	response := &gosnmp.SnmpPacket{
		Version:    msg.Version,
		Community:  s.community,
		RequestID:  msg.RequestID,
		Error:      0,
		ErrorIndex: 0,
	}

	switch msg.Type {
	case gosnmp.GetRequest:
		response.Type = gosnmp.GetResponse
		response.Variables = s.handleGetRequest(msg.Variables)
	case gosnmp.GetNextRequest:
		response.Type = gosnmp.GetResponse
		response.Variables = s.handleGetNextRequest(msg.Variables)
	default:
		return nil
	}

	return response
}

// handleGetRequest handles SNMP GET requests
func (s *SwitchSimulator) handleGetRequest(vars []gosnmp.SnmpPDU) []gosnmp.SnmpPDU {
	result := make([]gosnmp.SnmpPDU, len(vars))
	for i, v := range vars {
		result[i] = s.getMetricValue(v.Name)
	}
	return result
}

// handleGetNextRequest handles SNMP GETNEXT requests
func (s *SwitchSimulator) handleGetNextRequest(vars []gosnmp.SnmpPDU) []gosnmp.SnmpPDU {
	result := make([]gosnmp.SnmpPDU, len(vars))
	for i, v := range vars {
		result[i] = s.getNextMetricValue(v.Name)
	}
	return result
}

// getMetricValue returns the value for a specific OID
func (s *SwitchSimulator) getMetricValue(oid string) gosnmp.SnmpPDU {
	switch oid {
	// System metrics
	case OIDSysDescr:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.OctetString,
			Value: []byte(s.metrics.Description),
		}
	case OIDSysObjectID:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.ObjectIdentifier,
			Value: ".1.3.6.1.4.1.9.1.1", // Simulated Cisco Catalyst
		}
	case OIDSysUpTime:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.TimeTicks,
			Value: s.metrics.Uptime,
		}
	case OIDSysContact:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.OctetString,
			Value: []byte(s.metrics.Contact),
		}
	case OIDSysName:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.OctetString,
			Value: []byte(s.metrics.Hostname),
		}
	case OIDSysLocation:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.OctetString,
			Value: []byte(s.metrics.Location),
		}
	case OIDSysServices:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: int(72), // Layer 2 and 3 services
		}
	case OIDIfNumber:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: len(s.metrics.Interfaces),
		}
	case OIDCPUUsage:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: int(s.metrics.CPUUsage),
		}
	case OIDMemoryUsage:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: int(s.metrics.MemoryUsage),
		}
	case OIDTemperature:
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: int(s.metrics.Temperature),
		}
	}

	// Handle interface table entries
	if ifIndex := s.parseInterfaceIndex(oid); ifIndex >= 0 {
		return s.getInterfaceMetric(oid, ifIndex)
	}

	return gosnmp.SnmpPDU{
		Name:  oid,
		Type:  gosnmp.NoSuchObject,
		Value: nil,
	}
}

// parseInterfaceIndex extracts the interface index from an OID
func (s *SwitchSimulator) parseInterfaceIndex(oid string) int {
	var ifIndex int
	if n, err := fmt.Sscanf(oid, OIDIfEntry+".1.%d", &ifIndex); err == nil && n == 1 {
		if ifIndex > 0 && ifIndex <= len(s.metrics.Interfaces) {
			return ifIndex - 1
		}
	}
	return -1
}

// getInterfaceMetric returns interface-specific metrics
func (s *SwitchSimulator) getInterfaceMetric(oid string, ifIndex int) gosnmp.SnmpPDU {
	iface := s.metrics.Interfaces[ifIndex]

	switch {
	case strings.HasPrefix(oid, OIDIfIndex):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: iface.Index,
		}
	case strings.HasPrefix(oid, OIDIfDescr):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.OctetString,
			Value: []byte(iface.Name),
		}
	case strings.HasPrefix(oid, OIDIfType):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: 6, // ethernetCsmacd
		}
	case strings.HasPrefix(oid, OIDIfSpeed):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Gauge32,
			Value: iface.Speed,
		}
	case strings.HasPrefix(oid, OIDIfAdminStatus):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: iface.AdminStatus,
		}
	case strings.HasPrefix(oid, OIDIfOperStatus):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Integer,
			Value: iface.OperStatus,
		}
	case strings.HasPrefix(oid, OIDIfInOctets):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Counter64,
			Value: iface.InOctets,
		}
	case strings.HasPrefix(oid, OIDIfOutOctets):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Counter64,
			Value: iface.OutOctets,
		}
	case strings.HasPrefix(oid, OIDIfInErrors):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Counter32,
			Value: iface.InErrors,
		}
	case strings.HasPrefix(oid, OIDIfOutErrors):
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.Counter32,
			Value: iface.OutErrors,
		}
	}

	return gosnmp.SnmpPDU{
		Name:  oid,
		Type:  gosnmp.NoSuchObject,
		Value: nil,
	}
}

// getNextMetricValue returns the next value in the OID tree
func (s *SwitchSimulator) getNextMetricValue(oid string) gosnmp.SnmpPDU {
	// Find the next OID in our tree
	var nextOID string
	var found bool

	// System OIDs
	systemOIDs := []string{
		OIDSysDescr,
		OIDSysObjectID,
		OIDSysUpTime,
		OIDSysContact,
		OIDSysName,
		OIDSysLocation,
		OIDSysServices,
		OIDIfNumber,
	}

	// Find next system OID
	for _, candidate := range systemOIDs {
		if oid < candidate {
			nextOID = candidate
			found = true
			break
		}
	}

	// If not found in system OIDs, check interface table
	if !found {
		if oid < OIDIfEntry+".1.1" {
			nextOID = OIDIfEntry + ".1.1"
			found = true
		} else {
			// Parse current interface index
			var currentIndex int
			fmt.Sscanf(oid, OIDIfEntry+".1.%d", &currentIndex)

			if currentIndex < len(s.metrics.Interfaces) {
				nextOID = fmt.Sprintf("%s.1.%d", OIDIfEntry, currentIndex+1)
				found = true
			}
		}
	}

	// If not found in interface table, check resource OIDs
	if !found {
		resourceOIDs := []string{
			OIDCPUUsage,
			OIDMemoryUsage,
			OIDTemperature,
		}

		for _, candidate := range resourceOIDs {
			if oid < candidate {
				nextOID = candidate
				found = true
				break
			}
		}
	}

	if !found {
		return gosnmp.SnmpPDU{
			Name:  oid,
			Type:  gosnmp.EndOfMibView,
			Value: nil,
		}
	}

	return s.getMetricValue(nextOID)
}

// DumpMetrics returns the current metrics as JSON
func (s *SwitchSimulator) DumpMetrics() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.MarshalIndent(s.metrics, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling metrics: %w", err)
	}

	return string(data), nil
}
