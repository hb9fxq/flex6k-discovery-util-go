# flex6k-discovery-util-go

Lightweight replacement for: https://github.com/krippendorf/flex6k-discovery-util
util to use FRS 6XXX(R) signature series radios across subnets / routed VPNs

Currently tested on PfSense routers.


## Download precompiled for your system

Most recent binary release is here (precompiled)
https://github.com/krippendorf/flex6k-discovery-util-go/files/586780/flex6k-discovery-util-go-0.1-REL.zip

...currently for 386/amd64 Linux (ubuntu, etc) FreeBSD (PfSense) and ARM 5 linux (RaspberryPi)
If your platform isn't listed, send me a pull request for the pretty simple build.sh file. 

## Example Usage

### Simple setup (Client/Server)

* 192.168.92.1 is on a VPN site with [n] radios in the subnet
* 192.168.40.1 is a VPN router, connected via a TUN device and routes 192.168.92.0/40

#### Server VPN/Network site
```
./flex6k-discovery-util-go --SERVERIP=192.168.92.1 --SERVERPORT=7777
```

#### Client VPN/Network site

if you need to redirect the traffic on clientside to anything other than 255.255.255.255 (default) you can apply the ``` LOCALBR``` argument e.g. ``` --LOCALBR=192.168.40.255```

```
./flex6k-discovery-util-go --REMOTES=192.168.92.1:7777 --LOCALIFIP=192.168.40.1 --LOCALPORT=7788
```

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
```
./flex6k-discovery-util-go --REMOTES=192.168.92.1:7777;REMOTES=192.168.87.1:7777 --LOCALIFIP=192.168.40.1 --LOCALPORT=7788
```


### Multisite
If you host one or more radios on both sides of the tunnel or in different subnets and want to share accross networks you can run CLIENT & SERVER mode at the same time with the same process. 

#### Multisite node / relay loop

```
 ./flex6k-discovery-util-go --REMOTES=192.168.92.1:7777;REMOTES=192.168.87.1:7777 --LOCALIFIP=192.168.40.1 --LOCALPORT=7788 --SERVERIP=192.168.40.1 --SERVERPORT=7777
 ```

