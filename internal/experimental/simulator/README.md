# Experimental SNMP Simulator

This is an experimental Go-based SNMP simulator that was developed to provide dynamic metrics for testing. However, we are currently using the Python-based `snmpsim` package (in the `/snmpsim` directory) as our primary simulator.

This code is kept for reference and potential future use, but it is not actively used in the project.

## Features

- Dynamic metrics that change over time
- Realistic network switch behaviour
- Written in Go for better integration with our codebase
- Support for system, interface, and resource metrics

## Status

- **NOT IN USE** - We are using the Python-based `snmpsim` package instead
- Code is maintained for reference only
- May be developed further in the future if we need dynamic metrics

## Design Notes

The simulator was designed to provide more realistic network switch behaviour by:
1. Generating dynamic interface statistics
2. Simulating CPU and memory usage patterns
3. Providing temperature sensor data
4. Supporting multiple interfaces with realistic traffic patterns

## Future Considerations

If we need to simulate more complex scenarios, such as:
- Dynamic traffic patterns
- Fault conditions
- SNMP traps
- Custom MIB support

We may revisit this implementation.
