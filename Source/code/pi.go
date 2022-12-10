package sdnewtek

// ShortcutPI Property Inspector setting for Shortcut action
type ShortcutPI struct {
	Shortcut string `json:"shortcut"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	// 本当はmapにしたいけど、JavaScript側の実装の関係上難しいので一旦１個のみにする
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p *ShortcutPI) IsDefault() bool {
	if p.Shortcut == "" && p.Host == "" && p.User == "" && p.Password == "" && p.Key == "" && p.Value == "" {
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
		Key:      "",
		Value:    "",
	}
}
