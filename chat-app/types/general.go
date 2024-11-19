package types

import "time"

type WsResponseDto struct {
	OnlineCount int32      `json:"onlineCount"`
	Message     MessageDto `json:"message"`
}

type Invite struct {
	ChannelId string        `json:"channelId"`
	Exp      time.Duration `json:"exp"`
}
