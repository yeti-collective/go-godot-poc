package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"server/pkg/packets"
)

func main() {
	packet := &packets.Packet{
		SenderId: 69,
		Msg:      packets.NewChat("Hello, world!"),
	}

	fmt.Println(packet)

	data, err := proto.Marshal(packet)
	if err != nil {
		fmt.Println("Error marshaling packet:", err)
		return
	}

	fmt.Println(data)
}
