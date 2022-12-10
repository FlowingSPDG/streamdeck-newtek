package sdnewtek

import (
	"time"

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

// VideoPreviewPI Property Inspector setting for Video Thumbnail Preview
type VideoPreviewPI struct {
	Host           string        `json:"host"`
	User           string        `json:"user"`
	Password       string        `json:"password"`
	Name           string        `json:"name"`
	RefreshSeconds time.Duration `json:"refresh_seconds,string"` // in Mili seconds
}

func (p *VideoPreviewPI) IsDefault() bool {
	if p.Name == "" && p.Host == "" && p.User == "" && p.Password == "" && p.RefreshSeconds <= 0 {
		return true
	}
	return false
}

func (p *VideoPreviewPI) Initialize() {
	p.User = "admin"
	p.Password = "admin"
	p.Name = "output1"
	p.RefreshSeconds = time.Second
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
	//
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
