/*
 .... work in progress... use on own risk!

 2016 by Frank Werner-Krippendorf / HB9FXQ, 2016 mail@hb9fxq.ch

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:
 The above copyright notice and this permission notice shall be included in
 all copies or substantial portions of the Software.
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 THE SOFTWARE.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"net"
	"strconv"
	"strings"
	"time"
	"sync"
)

type AppContext struct {
	serverIp      string   // registraton server IP & PORT
	serverPort    int

	localIp       string   // client listener IP & PORT
	localPort     int

	remotes       []string // remotes to be notified

	registrations map[string]ListenerRegistration

	sync.Mutex
}

type ListenerRegistration struct {
	listenerPort int
	listenerIp   string
	raw          string
	since        int64
}

const NDEF_STRING string = "NDEF"
const FRS_DISCOVEY_ADDR string = "255.255.255.255:4992"

func main() {
	appctx := new(AppContext)

	var remotes string

	flag.StringVar(&remotes, "REMOTES", NDEF_STRING, "List of remotes to be notified. Delimited by ;   e.g. 192.168.62.1:7224;192.168.63.1:7228")
	flag.StringVar(&appctx.localIp, "LOCALIFIP", NDEF_STRING, "Local interface IP, interface, on that a client listens for relayd pks from a server")
	flag.IntVar(&appctx.localPort, "LOCALPORT", 0, "Port, that the client listens on for server pkgs")
	flag.StringVar(&appctx.serverIp, "SERVERIP", NDEF_STRING, "Broadcast server IP address")
	flag.IntVar(&appctx.serverPort, "SERVERPORT", 0, "Proadcast server port")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    TODO ...\n")

		flag.PrintDefaults()
	}

	if (remotes != NDEF_STRING && appctx.localIp != NDEF_STRING) {
		appctx.remotes = strings.Split(remotes, ";")
		go NotifyRemotes(appctx)
		go ListenForRelayedPkgs(appctx)
	}

	if (appctx.serverIp != NDEF_STRING && 0 < appctx.serverPort) {
		appctx.registrations = make(map[string]ListenerRegistration)
		fmt.Printf("Starting opmode SERVER on %s:%d \n", appctx.serverIp, appctx.serverPort)
		go BroadcastListener(appctx);
		go ServerListener(appctx);
	}

	fmt.Scanln()
}
func ListenForRelayedPkgs(appctx *AppContext) {
	ListenerLocalAddress, err := net.ResolveUDPAddr("udp4", appctx.localIp + ":" + strconv.Itoa(appctx.localPort))
	CheckError("Listener reslolve local", err)

	ServerConn, err := net.ListenUDP("udp4", ListenerLocalAddress)
	CheckError("Listener listen", err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Got relayed flex discovery: ", string(buf[0:n]), " from ", addr)

		relayLocal(appctx, buf[0:n])

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
func relayLocal(appctx *AppContext, bytes []byte) {
	fmt.Printf("...broadcasting in local subnet\n")

	ServerAddr, err := net.ResolveUDPAddr("udp4", FRS_DISCOVEY_ADDR)
	CheckError("broadcasting net.ResolveUDPAddr I", err)

	LocalAddr, err := net.ResolveUDPAddr("udp4", appctx.localIp + ":0")
	CheckError("broadcasting net.ResolveUDPAddr II", err)

	Conn, err := net.DialUDP("udp4", LocalAddr, ServerAddr)
	CheckError("broadcasting DialUDP", err)

	defer Conn.Close()

	_, ewrite := Conn.Write(bytes)

	if ewrite != nil {
		fmt.Println("Failed to broadcast", err)
	}
}

func NotifyRemotes(appctx *AppContext) {

	for {
		for _, remote := range appctx.remotes {
			fmt.Printf("Notifying remote %s\n", remote)

			ServerAddr, err := net.ResolveUDPAddr("udp4", remote)
			CheckError("net.ResolveUDPAddr I", err)

			LocalAddr, err := net.ResolveUDPAddr("udp4", appctx.localIp + ":0")
			CheckError("net.ResolveUDPAddr II", err)

			Conn, err := net.DialUDP("udp4", LocalAddr, ServerAddr)
			CheckError("DialUDP", err)

			defer Conn.Close()

			msg := "R;" + appctx.localIp + ";" + strconv.Itoa(appctx.localPort);


			buf := []byte(msg)
			_, ewrite := Conn.Write(buf)

			if ewrite != nil {
				fmt.Println(msg, err)
			}

		}

		time.Sleep(time.Second * 10)
	}

}

func ServerListener(appctx *AppContext) {

	FLexBroadcastAddr, err := net.ResolveUDPAddr("udp4", appctx.serverIp + ":" + strconv.Itoa(appctx.serverPort))
	CheckError("SRV FIND IP", err)

	ServerConn, err := net.ListenUDP("udp4", FLexBroadcastAddr)
	CheckError("SRV LISTEN", err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		content := string(buf[0:n])

		fmt.Println("REGISTRATION ", content, " from ", addr)

		tokens := strings.Split(content, ";")

		tokenPort, err := strconv.Atoi(tokens[2])
		CheckError("PARSE REG CONTENT", err)

		appctx.Lock();
		appctx.registrations[content] = ListenerRegistration{listenerIp: tokens[1], listenerPort: tokenPort, since:getCurrentUtcLinux(), raw:content}
		appctx.Unlock();

		fmt.Printf("SRV: Number of regs: %d\n", len(appctx.registrations))

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}

func getCurrentUtcLinux() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}

func notifyListener(appctx *AppContext, listener ListenerRegistration, msg []byte) {
	fmt.Printf("Notifying listener %s\n", listener.raw)

	ListenerAddr, err := net.ResolveUDPAddr("udp4", listener.listenerIp + ":" + strconv.Itoa(listener.listenerPort))

	if err != nil {
		fmt.Println("Could not notify listener", err)
		return
	}

	LocalAddr, err := net.ResolveUDPAddr("udp4", appctx.serverIp + ":0")
	if err != nil {
		fmt.Println("Could not notify listener", err)
		return
	}

	Conn, err := net.DialUDP("udp4", LocalAddr, ListenerAddr)
	if err != nil {
		fmt.Println("Could not notify listener", err)
		return
	}

	defer Conn.Close()

	_, ewrite := Conn.Write(msg)

	if ewrite != nil {
		fmt.Println(msg, err)
	}
}

func BroadcastListener(appctx *AppContext) {

	LocalAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:4992")
	CheckError("BR resolve broadcast A", err)

	ServerConn, err := net.ListenUDP("udp4", LocalAddr)
	CheckError("BR listen", err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Flex discovery: ", string(buf[0:n]), " from ", addr)

		appctx.Lock()

		if (0 < len(appctx.registrations)) {
			for _, registration := range appctx.registrations {
				if (registration.since + 30000 < getCurrentUtcLinux()) {
					delete(appctx.registrations, registration.raw)
					fmt.Printf("TTL for registration %s:%d\n", registration.listenerIp, registration.listenerPort)
					continue
				}

				if (addr.IP.String() != registration.listenerIp) {
					go notifyListener(appctx, registration, buf[0:n])
				}

			}
		}

		appctx.Unlock()

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}

func CheckError(where string, err error) {
	if err != nil {
		fmt.Println("FATAL: (" + where + ") ", err)
		os.Exit(0)
	}
}