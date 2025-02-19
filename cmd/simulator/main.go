package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sixworks/go-snmp-prometheus-getter/internal/simulator"
)

func main() {
	community := flag.String("community", "public", "SNMP community string")
	port := flag.Uint("port", 161, "SNMP port to listen on")
	flag.Parse()

	sim := simulator.NewSwitchSimulator(*community, uint16(*port))
	if err := sim.Start(); err != nil {
		log.Fatalf("Failed to start simulator: %v", err)
	}

	fmt.Printf("SNMP simulator listening on port %d with community '%s'\n", *port, *community)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down simulator...")
	if err := sim.Stop(); err != nil {
		log.Printf("Error stopping simulator: %v", err)
	}
}
