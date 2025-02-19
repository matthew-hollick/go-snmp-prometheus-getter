# SNMP Simulator

This directory contains the SNMP simulator configuration for development and testing.

## Sample Device

The simulator is configured with a sample network switch device with the following characteristics:

### System Information
- Name: switch01.hedgehog.internal
- Location: Server Room
- Contact: hedgehog_admin
- System Description: Sample Network Switch

### Interfaces
1. Management Interface
   - Type: Ethernet
   - Speed: 1 Gbps
   - Status: Up
   - Counters: Sample in/out bytes

2. User Network Interface
   - Type: Ethernet
   - Speed: 10 Gbps
   - Status: Up
   - Counters: Sample in/out bytes

### System Resources
- Physical Memory
  - Total: 8192 units
  - Used: 6144 units
  - Free: 2048 units
- CPU Usage: 45%

## Testing the Simulator

You can test the SNMP simulator using snmpget or snmpwalk:

```bash
# Get system description
snmpget -v2c -c public localhost:11161 1.3.6.1.2.1.1.1.0

# Walk all interfaces
snmpwalk -v2c -c public localhost:11161 1.3.6.1.2.1.2.2.1

# Get CPU usage
snmpget -v2c -c public localhost:11161 1.3.6.1.2.1.25.3.3.1.2.1
```

## Adding New Devices

To add a new device:

1. Create a new .snmprec file in the `data` directory
2. Use the following format for each line:
   ```
   OID|type-tag|value
   ```
   Where type-tag is:
   - 2: Integer
   - 4: Octet string
   - 6: Object identifier
   - 65: Counter32
   - 66: Gauge32
   - 67: TimeTicks
   - 70: Counter64

3. Update the docker-compose.yml if you need to expose additional ports

## MIB References

The sample device implements parts of the following MIBs:
- SNMPv2-MIB
- IF-MIB
- HOST-RESOURCES-MIB
