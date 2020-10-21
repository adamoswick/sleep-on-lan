package main
 
import (
    "net"
    "strconv"
    "strings"
    "os/exec"
    "os"
    "flag"
    "log"
)

// Set flag variables
var port string // Listen port
var logPath string // Set log path
var testMode bool // Set test mode (I.E. Don't actually poweroff just log)


func isWoLPacket(buffer []byte) bool {
	  // Checks first 6 bytes for the "FF FF FF FF FF FF" payload
		var firstSixBytes string
		// Loop through bytes 0 - 5 and append to string
	  for i := 0; i < 6; i++ {
		    char := strconv.FormatInt(int64(buffer[i]), 16)
		  	firstSixBytes = firstSixBytes + char
		}
		// If string contains the payload of F's, return True
		// Else return False
		if firstSixBytes == "ffffffffffff" {
			  return true
		} else {
			  return false
		}
}

func getMacAddressFromPacket(buffer []byte) string {
	  // Create MacAddress variable
	  var macAddress string
	  // Loop through bytes 6 to 12
	  for i := 6; i < 12; i++ {
	  	// Convert to hex
	  	char := strconv.FormatInt(int64(buffer[i]), 16)
	  	// Apend to MacAddress variable and add : after each section
	  	macAddress = macAddress + char + ":"
	  }
    // Return the MacAddress in the packet (removing the trailing :)
	  return strings.ToUpper(macAddress[:len(macAddress)-1])
}

func getInterfaces() []net.Interface {
	// Get interfaces on local device
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error. Unable to get interfaces:", err)
	}
	// Return them for use in checkIfInterfaceExists()
	return interfaces
}

func checkIfInterfaceExists(mac string) bool {
	// Get interfaces
	interfaces := getInterfaces()
	// Loop through interfaces
	for _, iface := range interfaces {
		// If interface matches function input, return true
		if strings.ToUpper(mac) == strings.ToUpper(iface.HardwareAddr.String()) {
			return true
		}
	}
	// If not, return false
	return false
}

func init() {
    flag.StringVar(&port, "port", "9", "Set the WoL listen port (default is UDP/9)")
    flag.StringVar(&logPath, "log-path", "", "File to log to (default is stdout)")
    flag.BoolVar(&testMode, "test-mode", false, "Don't poweroff, just log attempts")
    flag.Parse()

    if logPath == "" {
    	  log.Println("Logging to stdout")
    } else {
    	  logFile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    	  if err != nil {
    	  	log.Fatal("ERROR. Could not open the log file:", err)
    	  } else {
    	  	log.SetOutput(logFile)
    	  }
    	  log.Println("Logging to file", logPath)
    }
    if port == "9" {
        log.Println("Listening on port 9/UDP (default)")
    } else {
        log.Println("Listening on custom port,", port + "/UDP")
    }
    if testMode {
    	  log.Println("Test mode enabled")
    }
}
 
func main() {
    // Open UDP address
    port = ":" + port
    ServerAddr, err := net.ResolveUDPAddr("udp", port)
		if err != nil {
			  log.Fatal("ERROR. Could not prepare an address: ",err)
		} 
 
    // Listen on UDP port and defer close
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
		if err != nil {
			  log.Fatal("ERROR. Could not bind to port: ",err)
		}
    defer ServerConn.Close()
 
    // Create buffer for incoming packets
    packetBuffer := make([]byte, 1024)

    // Begin to for loop through packets in buffer
    log.Println("Listening for WoL packets")
    for {
        n, addr, err := ServerConn.ReadFromUDP(packetBuffer)
        if err != nil {
            log.Fatal("ERROR. Could not read from packet buffer: ",err)
        } 

        // Check length is correct, then check if packet is WoL
        if n == 102 && isWoLPacket(packetBuffer) {
        		log.Println("WoL Packet Sent From:", addr)

        		// Get MAC address from packet
        		packetMacAddress := getMacAddressFromPacket(packetBuffer)
        		log.Println("MAC Address in WoL Packet Is:", packetMacAddress)

        		// See if matching MAC address exists
        		if checkIfInterfaceExists(packetMacAddress) {
        			log.Println("This MAC Address Exists On The Device.")

        			// Is test mode enabled? 
        		  if testMode {
        		  	log.Println("Not shutting down due to test mode being enabled")
        		  } else {
        		  	log.Println("Shutting down now!")

        		  	// Try to power off
        		  	cmd := exec.Command("poweroff")
        		  	err := cmd.Run()
        		  	if err != nil {
        		  		log.Println("ERROR. Failed to shutdown:", err)
        		  	}
        		  }

            // Write to log if matching MAC address does not exist
        		} else {
        			log.Println("This MAC Address Is Not Associated With This Device")
        		}
        }

    }
}