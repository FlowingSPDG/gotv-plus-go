# GOTV+ for Go  

## Author
Shugo Kawamura  
Github : [**FlowingSPDG**](http://github.com/FlowingSPDG)  
Twitter : [**@rjzefl**](http://twitter.com/rjzefl) / [**@FlowingSPDG**](http://twitter.com/FlowingSPDG) / [**@HiiragiNil**](http://twitter.com/HiiragiNil)

## About
This is [GOTV+](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Broadcast) broadcast server written in Go and Gin.  
`tv_broadcast` sends GOTV fragment data API server,and API server relays fragment data for CS:GO Clients.  
CS:GO clients will get those fragments through HTTP,so viewer can play GOTV like a demo file.(e.g. Shift+F2 demoui, 1/4 speed play, HLAE and other external effects.)  
Also you can add CDN for GOTV+ for multi region stable connections such as Major tournament.  
  
Basically this is same with [kunalpowar/csgo-broadacst](https://github.com/kunalpowar/csgo-broadcast)(Go with mux) and [S5Projects/csgo-broadcast](https://github.com/S5Projects/csgo-broadcast)(NodeJS/TypeScript with Express).  

## Usage
- Launch GOTV+ Server(`go run ./main.go -addr YOURGLOBALIP:8080 -auth gopher`) or binary file : (`./gotv_plus -addr YOURGLOBALIP:8080 -auth gopher`)
- Enable GOTV broadcast in CS:GO Server
(`tv_enable 1;tv_broadcast_url "http://YOURGLOBALIP:8080/id/YOUR_MATCH_ID"; tv_broadcast_origin_auth "gopher"; tv_broadcast 1`)
- Play broadcast by CS:GO Client(`playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID"`)
### Tips
- You can get list of matches by http://localhost:8080/ .
- You should configure `tv_broadcast_...` cvars **before** enable broadcast(`tv_broadcast 1`). And you may need to restart broadcast if you changed cvars.(`tv_broadcast 0;tv_broadcast 1`)
- `playcast` and `tv_broadcast_url` URL needs to be wrapped by double quotes (e.g. `playcast "http://example.com/match/MATCH_ID"` .) because of `http:` colon(`:`) is recognized as command seprator for CS:GO.
- `playcast` and `tv_broadcast_url` does not need slash(`/`) in last character of URL. for example `playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID` works properly,however, `playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID/` **WILL NOT** work properly. because requested URL will be `http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID//sync` (Double-slash).

## Hidden Options
There are several hidden options that is not on document. I haven't covored all of these options,but some of them are supported.  
- "F" (e.g. `playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID" f500`) will play match from `500` fragment. CS:GO client will send request to http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID/sync?fragment=500 .
- "A" (e.g. `playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID" a`) will play match from `1` fragment. CS:GO client will send request to http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID/sync?fragment=1 .
- "C" (e.g. `playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID" c`) will play match from `1` fragment. CS:GO client **WILL NOT** send request to http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID/sync . CS:GO client will get 1st fragment directly.
- "B" Option / Skips dem_stop control frame(?)

## Example with Public CDN
In this example we run the application with user `gotv`, and then use NGINX to proxy TCP/80 (HTTP) and TCP/443 (HTTPS) traffic to the application. 
We advice that you limit who can send POST requests (CS:GO servers external/internal IP address) directly to the service with local firewall (iptables, nftables etc), this is the reason why we limit in NGINX all requests to only GET. 

We use two page rules on Cloudflare. We bypass Cache for all requests to /sync as we want that to be served directly from the application, and then cache everything on rest of the URL's. 

```
gotv.domain.com/*/sync
Cache Level: Bypass
gotv.domain.com/*
Cache Level: Cache Everything
```

nginx config
```
server {
	listen 80 default_server;
	listen [::]:80 default_server;
	server_name _;
	location / {
		limit_except GET {
			deny all;
		}
		proxy_pass http://ip.ip.ip.ip:port; # Address to the gotv-plus-go application.
	}
}
server {
	listen 443 ssl;
	listen [::]:443 ssl;
	ssl_certificate /etc/ssl/cloudflare/cloudflare-cert.pem;
	ssl_certificate_key /etc/ssl/cloudflare/cloudflare-key.key;
	server_name gotv.domain.com;
	location / {
		limit_except GET {
			deny all;
		}
		proxy_pass http://ip.ip.ip.ip:port; # Address to the gotv-plus-go application.
	}
}
```

simple systemd service
```
[Unit]
Description=GOTV+ service
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=gotv
WorkingDirectory=/path/to/gotv # Required as the application needs template
ExecStart=/path/to/gotv_plus -addr ip.ip.ip.ip:port -auth gopher
[Install]
WantedBy=multi-user.target
```

## Features
- Multi matches Support
- Match management(Give specify ID for each matches for support better path like : `/match/this_is_good_match_path/sync` )
- Match Delete
- RtDelay/RcVage Support
- Front end support
- Internal gRPC API
- CDN Support(Tested with Google Cloud CDN)
- Hidden args support (F/A). F=fragment specification w/o sync request,  A = Play from [1]frag with sync request

## TODO
- Dynamic auth system
- Rich Front end support
- Hidden args support (C/B). C=play from first? B=?
