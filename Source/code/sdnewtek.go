package sdnewtek

import (
	"context"

	"github.com/FlowingSPDG/streamdeck"
)

const (
	// AppName Streamdeck plugin app name
	AppName = "dev.flowingspdg.newtek.sdPlugin"
)

// Actions
const (
	// ActionShortcut Shortcut Action
	ActionShortcut = "dev.flowingspdg.newtek.shortcuthttp"

	// ActionShortcutTCP Perform action thru TCP
	ActionShortcutTCP = "dev.flowingspdg.newtek.shortcuttcp"
)

// SDNewTek StreamDeck client
type SDNewTek struct {
	// クライアント情報は保持しない
	// ハンドシェイクの遅延が気になるのであれば、map[string]newtek.ClientV1 でIPベースで保持しても良いかもしれない
	sd *streamdeck.Client
}

// NewSDNewTek Get New StreamDeck plugin instance pointer
func NewSDNewTek(ctx context.Context, params streamdeck.RegistrationParams) *SDNewTek {
	ret := &SDNewTek{
		sd: nil,
	}

	client := streamdeck.NewClient(ctx, params)

	actionShortcut := client.Action(ActionShortcut)
	actionShortcut.RegisterHandler(streamdeck.WillAppear, ret.ShortcutWillAppearHandler)
	actionShortcut.RegisterHandler(streamdeck.KeyDown, ret.ShortcutKeyDownHandler)

	actionShortcutTCP := client.Action(ActionShortcutTCP)
	actionShortcutTCP.RegisterHandler(streamdeck.WillAppear, ret.ShortcutTCPWillAppearHandler)
	actionShortcutTCP.RegisterHandler(streamdeck.KeyDown, ret.ShortcutTCPKeyDownHandler)

	ret.sd = client

	return ret
}

// Run Start client
func (s *SDNewTek) Run(ctx context.Context) error {
	return s.sd.Run()
}
