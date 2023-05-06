package sdnewtek

import (
	"github.com/FlowingSPDG/newtek-go"
)

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
	p.User = "admin"
	p.Password = "admin"
}

// ShortcutTCPPI Property Inspector setting for Shortcut action
type ShortcutTCPPI struct {
	Shortcut string `json:"shortcut"`
	Host     string `json:"host"`
	Value    string `json:"value"`
}

func (p *ShortcutTCPPI) IsDefault() bool {
	if p.Shortcut == "" && p.Host == "" && p.Value == "" {
		return true
	}
	return false
}

func (p *ShortcutTCPPI) Initialize() {
	p.Host = "localhost"
	p.Shortcut = ""
	p.Value = ""
}

func (p *ShortcutTCPPI) ToShortcuts() newtek.Shortcuts {
	return newtek.Shortcuts{
		Shortcut: []newtek.Shortcut{
			{
				Name:  p.Shortcut,
				Value: p.Value,
			},
		},
	}
}

// DialShortcutTCPPI Property Inspector setting for Dial Shortcut action
type DialShortcutTCPPI struct {
	Host string `json:"host"`

	// Push
	PushShortcut string `json:"push_shortcut"`
	PushValue    string `json:"push_value"`

	// Rotate
	RotateShortcut string `json:"rotate_shortcut"`
	RotateValue    string `json:"rotate_value"`
	RotateUseTicks bool   `json:"rotate_use_ticks"`
	XOfTicks       string `json:"x_of_ticks"`

	// Touch
	TouchShortcut string `json:"touch_shortcut"`
	TouchValue    string `json:"touch_value"`
}

func (p *DialShortcutTCPPI) IsDefault() bool {
	return p.Host == ""
}

func (p *DialShortcutTCPPI) Initialize() {
	p.Host = "localhost"
}

func (p *DialShortcutTCPPI) PushToShortcuts() newtek.Shortcuts {
	return newtek.Shortcuts{
		Shortcut: []newtek.Shortcut{
			{
				Name:  p.PushShortcut,
				Value: p.PushValue,
			},
		},
	}
}

func (p *DialShortcutTCPPI) RotateToShortcuts(val string) newtek.Shortcuts {
	return newtek.Shortcuts{
		Shortcut: []newtek.Shortcut{
			{
				Name:  p.RotateShortcut,
				Value: val,
			},
		},
	}
}

func (p *DialShortcutTCPPI) TouchToShortcuts() newtek.Shortcuts {
	return newtek.Shortcuts{
		Shortcut: []newtek.Shortcut{
			{
				Name:  p.TouchShortcut,
				Value: p.TouchValue,
			},
		},
	}
}
