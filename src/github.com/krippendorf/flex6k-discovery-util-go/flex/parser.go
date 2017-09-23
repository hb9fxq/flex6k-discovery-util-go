package flex

import (
	"strings"
)

type DiscoveryPackage struct {
	Discovery_protocol_version string
	Model string
	Serial string
	Version string
	Nickname string
	Callsign string
	Ip string
	Port string
	Status string
	Inuse_ip string
	Inuse_host string
	Max_licensed_version string
	Radio_license_id string
}

func  Parse(msg []byte) (DiscoveryPackage) {

	s := string(msg[:])
	s = s[strings.Index(s, "discovery_protocol_version"):]
	tokens := strings.Split(s, " ")
	var discoveryPackage DiscoveryPackage

	for i := 0; i < len(tokens); i++ {


		values := strings.Split(tokens[i], "=")

		switch values[0] {
		case "discovery_protocol_version":
			discoveryPackage.Discovery_protocol_version = values[1]
			continue;
		case "model":
			discoveryPackage.Model = values[1]
			continue;
		case "serial":
			discoveryPackage.Serial = values[1]
			continue;
		case "version":
			discoveryPackage.Version = values[1]
			continue;
		case "nickname":
			discoveryPackage.Nickname = values[1]
			continue;
		case "callsign":
			discoveryPackage.Callsign = values[1]
			continue;
		case "ip":
			discoveryPackage.Ip = values[1]
			continue;
		case "port":
			discoveryPackage.Port = values[1]
			continue;
		case "status":
			discoveryPackage.Status = values[1]
			continue;
		case "inuse_ip":
			discoveryPackage.Inuse_ip = values[1]
			continue;
		case "inuse_host":
			discoveryPackage.Inuse_host = values[1]
			continue;
		case "max_licensed_version":
			discoveryPackage.Max_licensed_version = values[1]
			continue;
		case "radio_license_id":
			discoveryPackage.Radio_license_id = values[1]
			continue;
		}
	}

	return discoveryPackage
}