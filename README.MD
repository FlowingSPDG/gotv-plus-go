# GOTV+ for Go  

## Author
Shugo Kawamura  
Github : [**FlowingSPDG**](http://github.com/FlowingSPDG)  
Twitter : [**@FlowingSPDG**](http://twitter.com/FlowingSPDG) / [**@FlowingSPDG_EN**](http://twitter.com/FlowingSPDG_EN)

## About
This is [GOTV+](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Broadcast) broadcast server written in Go and Gin.  
Basically this is same with [kunalpowar/csgo-broadacst](https://github.com/kunalpowar/csgo-broadcast)(Go with mux) and [S5Projects/csgo-broadcast](https://github.com/S5Projects/csgo-broadcast)(NodeJS/TypeScript with Express).  

## Features
- Multi matches Support
- 

## TODO
- Match management(Give specify ID for each matches for support better path like : `/match/this_is_good_match_path/sync` )
- Dynamic auth system
- RtDelay/RcVage Support
- CDN Support
- Front end support
- Hidden args support (C/B/F/A). F=fragment specification w/o sync, C=play from first? A = Play from [1]frag? B=?
- Internal gRPC API