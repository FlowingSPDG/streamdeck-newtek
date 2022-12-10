package sdnewtek

// ShortcutPI Property Inspector setting for Shortcut action
type ShortcutPI struct {
	Shortcut string            `json:"shortcut"`
	Host     string            `json:"host"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	KV       map[string]string `json:"kv"`
}
