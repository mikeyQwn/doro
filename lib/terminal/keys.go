package terminal

type Key uint32

const (
	KEY_UNKNOWN     Key = 0
	KEY_ARROW_LEFT      = (27 << 0) + (91 << 8) + (68 << 16)
	KEY_ARROW_RIGHT     = (27 << 0) + (91 << 8) + (67 << 16)
	KEY_ARROW_DOWN      = (27 << 0) + (91 << 8) + (66 << 16)
	KEY_ARROW_UP        = (27 << 0) + (91 << 8) + (65 << 16)
	KEY_SPACE           = (32 << 0)
)
