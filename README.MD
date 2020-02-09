# GOTV+ for Go  

## Author
Shugo Kawamura  
Github : [**FlowingSPDG**](http://github.com/FlowingSPDG)  
Twitter : [**@FlowingSPDG**](http://twitter.com/FlowingSPDG) / [**@FlowingSPDG_EN**](http://twitter.com/FlowingSPDG_EN)

## About
This is [GOTV+](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Broadcast) broadcast server written in Go and Gin.  
tv_broadcast sends GOTV fragment data API server,and API server relays fragment data for CS:GO Clients.  
CS:GO clients will get those fragments through HTTP,so viewer can play GOTV like a demo file.(e.g. Shift+F2 demoui, 1/4 speed play, HLAE and other external effects.)  
Also you can (probably) add CDN for GOTV+ for multi region stable connections.  
  
Basically this is same with [kunalpowar/csgo-broadacst](https://github.com/kunalpowar/csgo-broadcast)(Go with mux) and [S5Projects/csgo-broadcast](https://github.com/S5Projects/csgo-broadcast)(NodeJS/TypeScript with Express).  

## Usage
- Launch GOTV+ Server(`go run ./main.go -addr YOURGLOBALIP:8080 -auth gopher`) or binary file : (`./gotv_plus -addr YOURGLOBALIP:8080 -auth gopher`)
- Enable GOTV broadcast in CS:GO Server
(`tv_enable 1;tv_broadcast_url "http://YOURGLOBALIP:8080/id/YOUR_MATCH_ID"; tv_broadcast_origin_auth "gopher"; tv_broadcast 1`)
- Play broadcast by CS:GO Client(`playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID`)
### Tips
- You can get list of matches by http://localhost:8080/ .
- You should configure `tv_broadcast_...` cvars **before** enable broadcast(`tv_broadcast 1`). And you may need to restart broadcast if you changed cvars.(`tv_broadcast 0;tv_broadcast 1`)
- `playcast` and `tv_broadcast_url` URL needs to be wrapped by double quotes (e.g. `playcast "http://example.com/match/MATCH_ID"` .) because of `http:` colon(`:`) is recognized as command seprator for CS:GO.
- `playcast` and `tv_broadcast_url` does not need slash(`/`) in last character of URL. for example `playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID` works properly,however, `playcast "http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID/` **WILL NOT** work properly. because requested URL will be `http://YOURGLOBALIP:8080/match/id/YOUR_MATCH_ID//sync` (Double-slash).

## Features
- Multi matches Support
- Match management(Give specify ID for each matches for support better path like : `/match/this_is_good_match_path/sync` )
- RtDelay/RcVage Support(?)
- Front end support
- Internal gRPC API

## TODO
- Dynamic auth system
- Match Delete(to save memory and Earth!)
- CDN Support
- Rich Front end support
- Hidden args support (C/B/F/A). F=fragment specification w/o sync, C=play from first? A = Play from [1]frag? B=?