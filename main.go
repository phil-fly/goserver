package main

import (
	"fmt"
	"github.com/phil-fly/goserver/rtsp/rtspserver"
	"github.com/phil-fly/goserver/rtsp/log"
	"github.com/phil-fly/goserver/rtsp/auth"
)

func main() {

	mode := "console"

	config := `{"level":0,"filename":"test.log"}`
	log.NewLogger(0, mode, config)
	//auth.AuthorizationHeader{}
	authinfo := auth.NewAuthDatabase("pot_rtspserver")
	authinfo.InsertUserRecord("admin","123456")
	server := rtspserver.New(authinfo)

	portNum := 8554
	err := server.Listen("0.0.0.0",portNum)
	if err != nil {
		fmt.Printf("Failed to bind port: %d\n", portNum)
		return
	}

	if !server.SetupTunnelingOverHTTP(80) ||
		!server.SetupTunnelingOverHTTP(8000) ||
		!server.SetupTunnelingOverHTTP(8080) {
		fmt.Printf("We use port %d for optional RTSP-over-HTTP tunneling, "+
			"or for HTTP live streaming (for indexed Transport Stream files only).\n", server.HTTPServerPortNum())
	} else {
		fmt.Println("(RTSP-over-HTTP tunneling is not available.)")
	}

	urlPrefix := server.RtspURLPrefix()
	fmt.Println("This server's URL: " + urlPrefix + "<filename>.")

	server.Start()

	select {}
}

