extends Node2D

const packets := preload("res://packets.gd")

@onready var chat_log: ChatLog = $GUI/ChatLog
@onready var line_input: LineEdit = $GUI/LineEdit

func _ready() -> void:
	# Connect the websocket client to our packet handling methods
	WebsocketClient.connected_to_server.connect(_on_ws_connected_to_server)
	WebsocketClient.connection_closed.connect(_on_ws_connection_closed)
	WebsocketClient.packet_received.connect(_on_ws_packet_received)

	# Connect the line input
	line_input.text_submitted.connect(_on_line_edit_text_submitted)
	
	chat_log.info("Connecting to server...")
	WebsocketClient.connect_to_url("ws://127.0.0.1:8080/ws")

func _on_ws_connected_to_server() -> void:
	chat_log.success("Connected to server!")
	
func _on_ws_connection_closed() -> void:
	chat_log.info("Connection closed")
	
func _on_ws_packet_received(packet: packets.Packet) -> void:
	var sender_id := packet.get_sender_id()
	if packet.has_id():
		_handle_id_msg(sender_id, packet.get_id())
	elif packet.has_chat():
		_handle_chat_msg(sender_id, packet.get_chat())	

# This is the method that gets called on registration with the server. It means we
# get to know our own client id.
func _handle_id_msg(_sender_id: int, id_msg: packets.IdMessage) -> void:
	var client_id := id_msg.get_id()
	chat_log.success("You are ID #: %d" % client_id)

# This one is the method that gets called when we receive a chat message.
func _handle_chat_msg(sender_id: int, chat_msg: packets.ChatMessage) -> void:
	chat_log.chat("Client %d" % sender_id, chat_msg.get_msg())

# This gets called when the line edit is submitted
func _on_line_edit_text_submitted(text: String) -> void:

	# Compose a packet out of the chat message so that we can send it to the server
	var packet := packets.Packet.new()
	var chat_msg := packet.new_chat()
	chat_msg.set_msg(text)

	# Send the packet, and add the message to the log so that we can track our own
	# contributions to the chat
	var err := WebsocketClient.send(packet)
	if err:
		chat_log.error("Error sending chat message")
	else:
		chat_log.chat("You", text)
	
	# Clear the input field
	line_input.text = ""
