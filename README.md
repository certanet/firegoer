![Firegoer](firegoer.png "Firegoer logo")

**NOTE**: This library is in very early development, with limited functionality

Firegoer provides a way of interacting with Cisco Firepower devices via their REST APIs in Go. Currently FTD devices using FDM (not FMC) are supported.
This aims to be the Firepyer library, but written in Go!

The following versions have been used in development (others should work but YMMV):
- Go 1.17
- FTD 6.6.1-91

## Usage

Import the Fdm type and create an instance, passing in your FTD hostname/IP, password and ignoring SSL verification (if using an untrusted/self-signed cert). Then call any of the available functions, such as getting and printing the hostname:

    package main

    import (
        "fmt"

        "github.com/certanet/firegoer/ftd"
    )

    func main() {
        fdm := ftd.Fdm{
            Host:     "192.168.45.45",
            Password: "Admin123",
            Verify:   false,
        }

        hostname := fdm.GetHostname()
        fmt.Println(hostname)
    }

Then compile and run against your FTD:

    go run main.go
    firepyer2120
