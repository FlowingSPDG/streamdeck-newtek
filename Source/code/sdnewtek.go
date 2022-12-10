package sdnewtek

import (
	"context"
	"fmt"
	"time"

	"github.com/FlowingSPDG/newtek-go"
	"github.com/FlowingSPDG/streamdeck"
	sdcontext "github.com/FlowingSPDG/streamdeck/context"
)

const (
	// AppName Streamdeck plugin app name
	AppName = "dev.flowingspdg.newtek.sdPlugin"
)

// Actions
const (
	// ActionShortcut Shortcut Action
	ActionShortcut = "dev.flowingspdg.newtek.shortcuthttp"
	// ActionVideoPreview VideoPreview action
	ActionVideoPreview = "dev.flowingspdg.newtek.videopreview"

	// ActionShortcutTCP Perform action thru TCP
	ActionShortcutTCP = "dev.flowingspdg.newtek.shortcuttcp"
)

type SDNewTek struct {
	// クライアント情報は保持しない
	// ハンドシェイクの遅延が気になるのであれば、map[string]newtek.ClientV1 でIPベースで保持しても良いかもしれない
	sd *streamdeck.Client

	shortcutContexts    map[string]struct{}
	shortcutTCPContexts map[string]struct{}
	videoPreviewContext map[string]VideoPreviewPI
}

func NewSDNewTek(ctx context.Context, params streamdeck.RegistrationParams) *SDNewTek {
	ret := &SDNewTek{
		sd:                  nil,
		shortcutContexts:    map[string]struct{}{},
		shortcutTCPContexts: map[string]struct{}{},
		videoPreviewContext: map[string]VideoPreviewPI{},
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
		ret.videoPreviewContext[event.Context] = VideoPreviewPI{}
		return nil
	})
	actionVideoPreview.RegisterHandler(streamdeck.KeyDown, ret.VideoPreviewKeyDownHandler)
	actionVideoPreview.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		delete(ret.videoPreviewContext, event.Context)
		return nil
	})
	actionVideoPreview.RegisterHandler(streamdeck.DidReceiveSettings, ret.VideoPreviewDidReceiveSettingsHandler)

	actionShortcutTCP := client.Action(ActionShortcutTCP)
	actionShortcutTCP.RegisterHandler(streamdeck.WillAppear, ret.ShortcutTCPWillAppearHandler)
	actionShortcutTCP.RegisterHandler(streamdeck.WillAppear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		ret.shortcutTCPContexts[event.Context] = struct{}{}
		return nil
	})
	actionShortcutTCP.RegisterHandler(streamdeck.KeyDown, ret.ShortcutTCPKeyDownHandler)
	actionShortcutTCP.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		delete(ret.shortcutTCPContexts, event.Context)
		return nil
	})

	ret.sd = client

	return ret
}

func (s *SDNewTek) Run(ctx context.Context) error {
	go s.videoPreviewGoroutine(ctx)
	return s.sd.Run()
}

// videoPreviewGoroutine ctxはStreamDeckのcontextを使う
func (s *SDNewTek) videoPreviewGoroutine(ctx context.Context) error {
	for {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			for ctxStr, pi := range s.videoPreviewContext {
				if !pi.KeepUpdated {
					continue
				}
				ctx = sdcontext.WithContext(ctx, ctxStr)
				c, err := newtek.NewClientV1(pi.Host, pi.User, pi.Password)
				if err != nil {
					msg := fmt.Sprintf("Failed to connect NewTek Client: %s", err)
					s.sd.LogMessage(msg)
					s.sd.ShowAlert(ctx)
					return err
				}

				img, err := c.VideoPreview(pi.Name, 144, 144, 25)
				if err != nil {
					msg := fmt.Sprintf("Failed to send Shortcut to NewTek Client: %s", err)
					s.sd.LogMessage(msg)
					s.sd.ShowAlert(ctx)
					return err
				}

				i, _ := streamdeck.Image(img)
				s.sd.SetImage(ctx, i, streamdeck.HardwareAndSoftware)
			}
		}
	}
}
