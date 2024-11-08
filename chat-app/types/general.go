package types

type WsResponseDto struct {
	OnlineCount int32      `json:"onlineCount"`
	Message     MessageDto `json:"message"`
}
