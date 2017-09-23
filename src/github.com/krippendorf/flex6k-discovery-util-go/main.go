/*
 .... work in progress... use on own risk!

 2016, 2017 by Frank Werner-Krippendorf / HB9FXQ, mail@hb9fxq.ch

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
	"reflect"
	"github.com/krippendorf/flex6k-discovery-util-go/flex"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
)

type AppContext struct {
	serverIp      string   // registraton server IP & PORT
	serverPort    int


	localIp       string   // client listener IP & PORT
	localPort     int
	localBroadcast     string

	allLocalIp       string   // client listener IP & PORT

	remotes       []string // remotes to be notified

	registrations map[string]ListenerRegistration

	lastState string
	sync.Mutex

	telegrambot *tgbotapi.BotAPI

	telegramToken string
	telegramChat int64

	lastPackage *flex.DiscoveryPackage
}

type ListenerRegistration struct {
	listenerPort int
	listenerIp   string
	raw          string
	since        int64
}

const NDEF_STRING string = "NDEF"
const FRS_DISCOVEY_ADDR string = "255.255.255.255:4992"
const UDP_NETWORK string = "udp4"

func main() {
	appctx := new(AppContext)

	var remotes string
	flag.StringVar(&remotes, "REMOTES", NDEF_STRING, "List remote server to subscribe to. One or more, format is [SERVER_IP:SERVER_PORT], if more than one, delimit subscriptions by ';'   e.g. --REMOTES=192.168.62.1:7224;192.168.63.1:7228")
	flag.StringVar(&appctx.localIp, "LOCALIFIP", NDEF_STRING, "Client local interface IPinterface, where servers will forward pkgs to")
	flag.IntVar(&appctx.localPort, "LOCALPORT", 0, "Local port")
	flag.StringVar(&appctx.localBroadcast, "LOCALBR", NDEF_STRING, "Local broadcast address address, default 255.255.255.255 - e.g. 192.168.2.255. Required on PfSense!")
	flag.StringVar(&appctx.serverIp, "SERVERIP", NDEF_STRING, "Broadcast server IP address")
	flag.IntVar(&appctx.serverPort, "SERVERPORT", 0, "Broadcast server port")
	flag.Parse()

	go loadTelegramBot(appctx);
	appctx.lastState = "empty"
	appctx.allLocalIp = FetchAllLocalIPs()
	fmt.Println("APP Identified local IPs: " + appctx.allLocalIp)

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    ..:: see https://github.com/krippendorf/flex6k-discovery-util-go for instructions ::..\n")

		flag.PrintDefaults()
	}

	if (remotes != NDEF_STRING && appctx.localIp != NDEF_STRING) {
		appctx.remotes = strings.Split(remotes, ";")
		go NotifyRemotes(appctx)
		go ListenForRelayedPkgs(appctx)
	}

	if (appctx.serverIp != NDEF_STRING && 0 < appctx.serverPort) {
		appctx.registrations = make(map[string]ListenerRegistration)
		fmt.Printf("SRV listening for registrations on: %s:%d \n", appctx.serverIp, appctx.serverPort)
		go BroadcastListener(appctx);
		go ServerListener(appctx);
	}

	for{
		time.Sleep(1*time.Second)
	}
}

func loadTelegramBot(context *AppContext) {

	telegramData := "/flexi/telegramData"


	if _, err := os.Stat(telegramData); os.IsNotExist(err) {
		fmt.Printf("Telegram integration disabled: %s", err)
		return
	}

	dat, err := ioutil.ReadFile(telegramData)
	telegramdata := strings.Split(string(dat), " ")
	context.telegramToken = telegramdata[0]

	chat, err := strconv.ParseInt(telegramdata[1], 10, 64);
	context.telegramChat =  chat

	if(err != nil || chat == 0){
		fmt.Printf("Telegram integration disabled: %s", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(context.telegramToken)

	if(err != nil){
		return
	}

	context.telegrambot = bot;
}

func FetchAllLocalIPs()(allips string) {

	allips = "0.0.0.0 127.0.0.1 "
	ifaces, err := net.Interfaces()
	CheckError("FetchAllLocalIPs", err)

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		CheckError("Fetch if IP", err)
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				allips += v.IP.String() + " "
			case *net.IPAddr:
				allips +=  v.IP.String() + " "
			}
		}
	}

	return allips
}

func ListenForRelayedPkgs(appctx *AppContext) {
	ListenerLocalAddress, err := net.ResolveUDPAddr(UDP_NETWORK, appctx.localIp + ":" + strconv.Itoa(appctx.localPort))
	CheckError("Listener reslolve local", err)
	if(err != nil){
		ListenForRelayedPkgs(appctx)
	}

	ServerConn, err := net.ListenUDP(UDP_NETWORK, ListenerLocalAddress)
	CheckError("Listener listen", err)

	if(err != nil){
		ListenForRelayedPkgs(appctx);
	}

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)


		if(err != nil){
			continue
		}

		if(strings.Contains(appctx.allLocalIp, addr.IP.String())){
			continue // skip, if comes from local server instance, if registered in local loop
		}

		fmt.Println("CLT RECEIVED PKG FROM SRV @", addr.IP.String())

		RelayLocal(appctx, buf[0:n])

		if err != nil {
			fmt.Println("CLT Error: ", err)
		}
	}
}

func notifyTelegramGroup(context *AppContext) {

	if(context.telegrambot == nil || context.lastPackage == nil || len(context.lastPackage.Status) == 0 ){
		return
	}

	filenameStatusImage := "/flexi/" +context.lastPackage.Status+ ".jpg"

	msgImg := tgbotapi.NewPhotoUpload(context.telegramChat, filenameStatusImage)
	context.telegrambot.Send(msgImg)

	msg := tgbotapi.NewMessage(context.telegramChat, "Radio "+ context.lastPackage.Serial + " state changed: '" + context.lastPackage.Status + "' " + context.lastPackage.Inuse_ip + " " + context.lastPackage.Inuse_host)
	context.telegrambot.Send(msg)
}

func RelayLocal(appctx *AppContext, bytes []byte) {
	fmt.Printf("    broadcasting in local subnet\n")

	defAddr := FRS_DISCOVEY_ADDR

	if(NDEF_STRING != appctx.localBroadcast){
		defAddr =  appctx.localBroadcast + ":4992";
	}

	ServerAddr, err := net.ResolveUDPAddr(UDP_NETWORK, defAddr)
	CheckError("broadcasting net.ResolveUDPAddr I", err)
	if(err != nil){
		return
	}

	LocalAddr, err := net.ResolveUDPAddr(UDP_NETWORK, appctx.localIp + ":" + strconv.Itoa(appctx.localPort + 1))
	CheckError("broadcasting net.ResolveUDPAddr II", err)
	if(err != nil){
		return
	}

	Conn, err := net.DialUDP(UDP_NETWORK, LocalAddr, ServerAddr)
	CheckError("broadcasting DialUDP", err)
	if(err != nil){
		return
	}

	defer Conn.Close()

	_, ewrite := Conn.Write(bytes)

	if ewrite != nil {
		fmt.Println("CLT Failed to broadcast", err)
	}
}

func NotifyRemotes(appctx *AppContext) {

	for {
		for _, remote := range appctx.remotes {
			fmt.Printf("	==> Notifying remote [%s]\n", remote)

			ServerAddr, err := net.ResolveUDPAddr(UDP_NETWORK, remote)
			CheckError("net.ResolveUDPAddr I", err)
			if(err != nil){
				continue;
			}

			LocalAddr, err := net.ResolveUDPAddr(UDP_NETWORK, appctx.localIp + ":0")
			CheckError("net.ResolveUDPAddr II", err)

			if(err != nil){
				continue;
			}

			Conn, err := net.DialUDP(UDP_NETWORK, LocalAddr, ServerAddr)
			CheckError("DialUDP", err)

			if(err != nil){
				continue;
			}

			msg := "R;" + appctx.localIp + ";" + strconv.Itoa(appctx.localPort);


			buf := []byte(msg)
			_, ewrite := Conn.Write(buf)

			if ewrite != nil {
				fmt.Println(msg, err)
			}

			Conn.Close();
		}

		time.Sleep(time.Second * 10)
	}
}

func ServerListener(appctx *AppContext) {

	FLexBroadcastAddr, err := net.ResolveUDPAddr(UDP_NETWORK, appctx.serverIp + ":" + strconv.Itoa(appctx.serverPort))
	CheckError("SRV FIND IP", err)
	if(err != nil){
		ServerListener(appctx);
	}

	ServerConn, err := net.ListenUDP(UDP_NETWORK, FLexBroadcastAddr)
	CheckError("SRV LISTEN", err)
	if(err != nil){
		ServerListener(appctx);
	}

	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		content := string(buf[0:n])

		fmt.Println("REGISTRATION ", content, " from ", addr)

		tokens := strings.Split(content, ";")

		tokenPort, err := strconv.Atoi(tokens[2])
		CheckError("PARSE REG CONTENT", err)
		if(err != nil){
			continue
		}

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

func NotifyListener(appctx *AppContext, listener ListenerRegistration, msg []byte) {
	fmt.Printf("	==> Notifying remote [%s]\n", listener.raw)

	ListenerAddr, err := net.ResolveUDPAddr(UDP_NETWORK, listener.listenerIp + ":" + strconv.Itoa(listener.listenerPort))

	if err != nil {
		fmt.Println("SRV ERR, Could not notify listener", err)
		return
	}

	LocalAddr, err := net.ResolveUDPAddr(UDP_NETWORK, appctx.serverIp + ":0")
	if err != nil {
		fmt.Println("SRV ERR, Could not notify listener", err)
		return
	}

	Conn, err := net.DialUDP(UDP_NETWORK, LocalAddr, ListenerAddr)
	if err != nil {
		fmt.Println("SRV ERR, Could not notify listener", err)
		return
	}

	defer Conn.Close()

	_, ewrite := Conn.Write(msg)

	if ewrite != nil {
		fmt.Println(msg, err)
	}
}

func BroadcastListener(appctx *AppContext) {

	LocalAddr := net.UDPAddr{IP: net.IPv4zero, Port:4992}

	ServerConn, err := net.ListenUDP(UDP_NETWORK, &LocalAddr)
	CheckError("BR listen", err)
	defer ServerConn.Close()


	buf := make([]byte, 1024)
	prev := make([]byte, 1024)

	var ackCnt int

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)

		if(!IsFrsFlexDiscoveryPkgInBuffer(buf, n)){
			continue // thats not ourts
		}

		if(reflect.DeepEqual(buf, prev)){
			continue; // skip own pkgs, that where captured on other local network interface
		}

		copy(prev, buf)

		if (  strings.Contains(appctx.allLocalIp, addr.IP.String())) {
			ackCnt++;
			fmt.Println("SRV ACK [" + strconv.Itoa(ackCnt) + "]")
			continue;
		}

		ackCnt = 0;
		fmt.Printf("SRV BROADCAST RECEIVED [%s]\n", addr)
		appctx.Lock()


		parsed := flex.Parse(buf[0:n])
		if(parsed.Status != appctx.lastState){
			appctx.lastState = parsed.Status;
			appctx.lastPackage = &parsed;
			go notifyTelegramGroup(appctx)
		}

		if (0 < len(appctx.registrations)) {
			for _, registration := range appctx.registrations {
				if (registration.since + 30000 < getCurrentUtcLinux()) {
					delete(appctx.registrations, registration.raw)
					fmt.Printf("TTL for registration %s:%d\n", registration.listenerIp, registration.listenerPort)
					continue
				}

				go NotifyListener(appctx, registration, buf[0:n])
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
		fmt.Println("FATAL: (" + where + ") sleeping 5 seconds", err)
		time.Sleep(5*time.Second)
	}
}

func IsFrsFlexDiscoveryPkgInBuffer(buf []byte, length int)(res bool){

	if(900<length){
		res = false
		fmt.Printf("ERROR: INVALID DATA, size: %d", length)
		return
	}

	content := string(buf[0:length])
	res = strings.Contains(content, "serial=") && strings.Contains(content, "version=") && strings.Contains(content, "ip=")

	if(!res){
		fmt.Println("ERROR: INVALID DATA")
	}

	return
}