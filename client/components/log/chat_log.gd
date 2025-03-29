extends RichTextLabel
class_name Log

# This method (underscore means local / private by convention, even though private methods aren't really a thing in Godot)
# will append the given text. BBCode is used to set the color for the text snippet to the given argument.
func _message(message: String, color: Color = Color.WHITE) -> void:
    append_text("[color=#%s]%s[/color]\n" % [color.to_html(false), str(message)])


# These are just convenience functions that add the given text with a preset color depending on the role.
func info(message: String) -> void:
    _message(message, Color.WHITE)

func warning(message: String) -> void:
    _message(message, Color.YELLOW)

func error(message: String) -> void:
    _message(message, Color.ORANGE_RED)

func success(message: String) -> void:
    _message(message, Color.LAWN_GREEN)
    
# This adds a chat message, which is really just a blue username followed by a default color (white, as per the arg list
# in _message, above) message in italics.
func chat(sender_name: String, message: String) -> void:
    _message("[color=#%s]%s:[/color] [i]%s[/i]" % [Color.CORNFLOWER_BLUE.to_html(false), sender_name, message])
