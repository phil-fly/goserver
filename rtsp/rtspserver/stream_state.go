package rtspserver

import "github.com/phil-fly/goserver/rtsp/livemedia"

type StreamServerState struct {
	subsession  livemedia.IServerMediaSubsession
	streamToken *livemedia.StreamState
}
