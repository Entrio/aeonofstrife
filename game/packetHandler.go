package game

import "fmt"

type PacketHandler interface {
	handle(packet *Packet)
}

type RoomCountHandler struct{}

/*
We are asking to send us a list of rooms
*/
func (r RoomCountHandler) handle(packet *Packet) {

	for _, v := range ServerInstance.roomList {
		fmt.Println(fmt.Sprintf("Sending room %s upstream", v.ID))
		msg := NewPacket(MsgRoomCountResponse)
		msg.WriteRoomData(v)
		sendMessageToConnection(packet.Connection, *msg)
	}
}
