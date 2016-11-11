# flex6k-discovery-util-go

Lightweight replacement for: https://github.com/krippendorf/flex6k-discovery-util

Currently tested on PfSense routers.

## Example Usage

### Simple setup (Client/Server)

#### Server VPN/Network site
```
./flex6k-discovery-util-go --SERVERIP=192.168.92.1 --SERVERPORT=7777
```

#### Client VPN/Network site
```
./flex6k-discovery-util-go --REMOTES=192.168.92.1:7777 --LOCALIFIP=192.168.40.1 --LOCALPORT=7788
```

### Multi server (Client/Server)

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
If you host radios on both sided of the tunnel or in different subnets and want to share accross networks you can run CLIENT & SERVER mode at the same time with the same process. 

#### Multisite node / relay loop

```
 ./flex6k-discovery-util-go --REMOTES=192.168.92.1:7777;REMOTES=192.168.87.1:7777 --LOCALIFIP=192.168.40.1 --LOCALPORT=7788 --SERVERIP=192.168.40.1 --SERVERPORT=7777
 ```

