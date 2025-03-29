extends Node2D

const packets := preload("res://packets.gd")

@onready var chat_log: ChatLog = $GUI/ChatLog

func _ready() -> void:
	WebsocketClient.connected_to_server.connect(_on_ws_connected_to_server)
	WebsocketClient.connection_closed.connect(_on_ws_connection_closed)
	WebsocketClient.packet_received.connect(_on_ws_packet_received)
	
	chat_log.info("Connecting to server...")
	WebsocketClient.connect_to_url("ws://127.0.0.1:8080/ws")

func _on_ws_connected_to_server() -> void:
	var packet := packets.Packet.new()
	var chat_msg := packet.new_chat()
	chat_msg.set_msg("Hello, Golang!")
	
	var err := WebsocketClient.send(packet)
	if err:
		chat_log.error("Error sending packet")
	else:
		chat_log.success("Sent packet")
	
func _on_ws_connection_closed() -> void:
	chat_log.info("Connection closed")
	
func _on_ws_packet_received(packet: packets.Packet) -> void:
	chat_log.info("Received packet from the server: %s" % packet)
