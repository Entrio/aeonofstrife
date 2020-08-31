package game

type PacketType uint16

const (
	MsgNullIota = PacketType(iota)
	MsgPingRequest
	MsgPingResponse
	MsgWelcome
	MsgSpecial1
	MsgSpecial2
	MsgRoomCountRequest
	MsgRoomCountResponse
	MsgRoomUpdateName       PacketType = 1000
	MsgUpdateRoomPayload    PacketType = 1001
	MsgUpdateRoomPayloadAck PacketType = 1002
	MsgRoomUpdateStatus     PacketType = 1003
	MsgRoomMiscUpdate       PacketType = 1004
)
