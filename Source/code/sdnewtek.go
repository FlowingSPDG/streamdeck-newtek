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
	// ActionCut Cut action
	ActionShortcut = "dev.flowingspdg.newtek.shortcut"

	ActionVideoPreview = "dev.flowingspdg.newtek.videopreview"
)

type SDNewTek struct {
	// クライアント情報は保持しない
	// ハンドシェイクの遅延が気になるのであれば、map[string]newtek.ClientV1 でIPベースで保持しても良いかもしれない
	sd *streamdeck.Client

	shortcutContexts    map[string]struct{}
	videoPreviewContext map[string]struct{}
}

func NewSDNewTek(ctx context.Context, params streamdeck.RegistrationParams) *SDNewTek {
	ret := &SDNewTek{
		sd:                  nil,
		shortcutContexts:    map[string]struct{}{},
		videoPreviewContext: map[string]struct{}{},
	}

	client := streamdeck.NewClient(ctx, params)

	actionShortcut := client.Action(ActionShortcut)
	actionShortcut.RegisterHandler(streamdeck.WillAppear, ret.ShortcutWillAppearHandler)
	actionShortcut.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		ret.shortcutContexts[event.Context] = struct{}{}
		return nil
	})
	actionShortcut.RegisterHandler(streamdeck.KeyDown, ret.ShortcutKeyDownHandler)
	actionShortcut.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		delete(ret.shortcutContexts, event.Context)
		return nil
	})

	actionVideoPreview := client.Action(ActionVideoPreview)
	actionVideoPreview.RegisterHandler(streamdeck.WillAppear, ret.VideoPreviewWillAppearHandler)
	actionVideoPreview.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		ret.videoPreviewContext[event.Context] = struct{}{}
		return nil
	})
	actionVideoPreview.RegisterHandler(streamdeck.KeyDown, ret.VideoPreviewKeyDownHandler)
	actionVideoPreview.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		delete(ret.videoPreviewContext, event.Context)
		return nil
	})

	ret.sd = client

	return ret
}

func (s *SDNewTek) Run() error {
	return s.sd.Run()
}
