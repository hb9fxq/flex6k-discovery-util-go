# flex6k-discovery-util-go

Lightweight replacement for: https://github.com/krippendorf/flex6k-discovery-util
util to use FRS 6XXX(R) signature series radios across subnets / routed VPNs


Currently tested on RasberryPi, X64 Linux, and Windows


## Download precompiled for your system

Most recent binary release is here (precompiled)
https://github.com/krippendorf/flex6k-discovery-util-go/releases/download/v0.1-br/flex6k-discovery-util-go-0.1-BR-REL_ret.zip

...currently for 386/amd64 Linux (ubuntu, etc and Windows) FreeBSD (PfSense) and ARM 5 linux (RaspberryPi)
If your platform isn't listed, send me a pull request for the pretty simple build.sh file. 

## Example Usage

How it works: You need a client and a server which establish a link. One of them has to be installed in the subnet where your radio is connected and the other where SmartSDR runs. 
It does not matter where you install the client and where the server.

### Simple setup (Client/Server):

* Server: 192.168.1.4 is on a VPN site with [n] radios in the subnet
* Client: 10.147.20.144 is a VPN router, connected via a TUN device and routes ?????  192.168.92.0/40
* 192.168.1.7 Radio in this example

#### Server VPN/Network site (server is installed close to the radio)
```
Linux command: ./flexi --SERVERIP=192.168.1.4 --SERVERPORT=7777
Windows command: flexi -SERVERIP=192.168.1.4 -SERVERPORT=7777
```

#### Client VPN/Network site (Client in the subnet of SmartSDR)

```
Linux command ./flexi --REMOTES=10.147.20.144:7777 --LOCALIFIP=192.168.1.4 --LOCALPORT=7788
Windows command: flexi -REMOTES=10.147.20.144:7777 -LOCALIFIP=192.168.1.4 -LOCALPORT=7788
```
The ports can be changed. SERVERPORT and REMOTES port has to be the same (here 7777)

If you execute these two commands the following output should result:
Client side:
CLT RECEIVED PKG FROM SRV @ 192.168.1.4
    broadcasting in local subnet
CLT RECEIVED PKG FROM SRV @ 192.168.1.4
    broadcasting in local subnet
CLT RECEIVED PKG FROM SRV @ 192.168.1.4
    broadcasting in local subnet

Server side:
REGISTRATION  R;10.147.20.144;7788  from  10.147.20.144:55973
SRV: Number of regs: 1
SRV BROADCAST RECEIVED [192.168.1.7:4992]
        ==> Notifying remote [R;10.147.20.144;7788]
SRV BROADCAST RECEIVED [192.168.1.7:4992]
        ==> Notifying remote [R;10.147.20.144;7788]

Hint: ZeroTier works also with LTE routers and non-public IP adress


-------------- END OF standard installation -----------


If you need to redirect the traffic on clientside to anything other than 255.255.255.255 (default) you can apply the ``` LOCALBR``` argument e.g. ``` --LOCALBR=192.168.40.255``` I've discoverd that especially PfSense drops UDP packages directly to 255.255.255.255 - probably due to the fact it does not decide on which interface to send out the traffic

### Multi server (Client/Server/Server)

#### Server 1 VPN/Network site
```
./flex6k-discovery-util-go --SERVERIP=192.168.92.1 --SERVERPORT=7777
```

#### Server 2 VPN/Network site
```
./flex6k-discovery-util-go --SERVERIP=192.168.87.1 --SERVERPORT=7777
```

#### Client VPN/Network site
Simple add all server to the REMOTES argument.

```
./flex6k-discovery-util-go --REMOTES=192.168.92.1:7777;192.168.87.1:7777 --LOCALIFIP=192.168.40.1 --LOCALPORT=7788
```


### Multi site
If you host one or more radios on both sides of the tunnel or in different subnets and want to share accross networks you can run CLIENT & SERVER mode at the same time with the same process. 

#### Multi site node / relay loop

```
 ./flex6k-discovery-util-go --REMOTES=192.168.92.1:7777;REMOTES=192.168.87.1:7777 --LOCALIFIP=192.168.40.1 --LOCALPORT=7788 --SERVERIP=192.168.40.1 --SERVERPORT=7777
 ```



