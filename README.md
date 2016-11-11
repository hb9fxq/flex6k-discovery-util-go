# flex6k-discovery-util-go

Lightweight replacement for: https://github.com/krippendorf/flex6k-discovery-util

Currently tested on PfSense routers.

## Example Usage

### Simple setup (Client/Server)

* 192.168.92.1 is on a VPN site with [n] radios in the subnet
* 192.168.40.1 is a VPN router, connected via a TUN device and routes 192.168.92.0/40

#### Server VPN/Network site
```
./flex6k-discovery-util-go --SERVERIP=192.168.92.1 --SERVERPORT=7777
```

#### Client VPN/Network site
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

