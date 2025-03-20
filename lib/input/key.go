package input

// A `Key` represents a single keypress
type Key uint32

// A list of keys required for Doro to work
const (
	KEY_UNKNOWN     Key = 0
	KEY_ARROW_LEFT  Key = (27 << 0) + (91 << 8) + (68 << 16)
	KEY_ARROW_RIGHT Key = (27 << 0) + (91 << 8) + (67 << 16)
	KEY_ARROW_DOWN  Key = (27 << 0) + (91 << 8) + (66 << 16)
	KEY_ARROW_UP    Key = (27 << 0) + (91 << 8) + (65 << 16)
	KEY_S           Key = (115 << 0)
	KEY_SPACE       Key = (32 << 0)
	KEY_ENTER       Key = (13 << 0)
	KEY_CTRL_C      Key = (3 << 0)
)
