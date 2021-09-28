package game

import "fmt"

func handleData(data []byte, receivedData *Packet) bool {
	packetLength := uint32(0)
	receivedData.SetBytes(data)

	if receivedData.UnreadLength() >= 4 {
		packetLength = receivedData.ReadUInt32()
		if packetLength <= 0 {
			return true
		}
	}

	for packetLength > 0 && packetLength <= receivedData.UnreadLength() {
		packetBytes := receivedData.ReadBytes(packetLength)

		//TODO: Maybe look into goroutine
		newPaket := NewUnknownPacket(packetBytes)
		newPaket.Connection = receivedData.Connection
		fmt.Println(fmt.Sprintf("Total handlers: %d", len(ServerInstance.packetHandler)))

		t := newPaket.GetMessageType()
		_, ok := ServerInstance.packetHandler[t]

		if !ok {
			fmt.Println(fmt.Sprintf("There is no handler registered for packet type %d", t))
		}

		ServerInstance.packetHandler[t].handle(newPaket)

		packetLength = 0
		if receivedData.UnreadLength() >= 4 {
			packetLength = receivedData.UnreadLength()
			if packetLength <= 0 {
				return true
			}
		}
	}

	if packetLength <= 1 {
		return true
	}

	return false
}
