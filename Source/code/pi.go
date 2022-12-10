package sdnewtek

// ShortcutPI Property Inspector setting for Shortcut action
type ShortcutPI struct {
	Shortcut string            `json:"shortcut"`
	Host     string            `json:"host"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	KV       map[string]string `json:"kv"`
}

func (p *ShortcutPI) IsDefault() bool {
	if p.Shortcut == "" && p.Host == "" && p.User == "" && p.Password == "" && len(p.KV) == 0 {
		return true
	}
	return false
}

func (p *ShortcutPI) Initialize() {
	p = &ShortcutPI{
		Shortcut: "",
		Host:     "",
		User:     "admin",
		Password: "admin",
		KV: map[string]string{
			"": "",
		},
	}
}
