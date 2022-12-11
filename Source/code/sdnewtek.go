package sdnewtek

import (
	"context"
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/FlowingSPDG/newtek-go"
	"github.com/FlowingSPDG/streamdeck"
	sdcontext "github.com/FlowingSPDG/streamdeck/context"
	"github.com/oliamb/cutter"
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
	videoPreviewContext sync.Map
}

func NewSDNewTek(ctx context.Context, params streamdeck.RegistrationParams) *SDNewTek {
	ret := &SDNewTek{
		sd:                  nil,
		shortcutContexts:    map[string]struct{}{},
		shortcutTCPContexts: map[string]struct{}{},
		videoPreviewContext: sync.Map{},
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
	actionVideoPreview.RegisterHandler(streamdeck.WillDisappear, func(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
		ret.videoPreviewContext.Delete(event.Context)
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

type cacher struct {
	m sync.Map
}

type cacherKey struct {
	host string
	name string
}

func (c *cacher) getCache(key cacherKey) (image.Image, bool) {
	v, ok := c.m.Load(key)
	if !ok {
		return nil, ok
	}
	i, o := v.(image.Image)
	if !o {
		return nil, o
	}
	return i, ok
}

func (c *cacher) storeCache(key cacherKey, img image.Image) {
	c.m.Store(key, img)
}

func (s *SDNewTek) videoPreviewGoroutine(ctx context.Context) error {
	// それぞれ独立のgoroutineでループしないと処理が止まる
	for {
		time.Sleep(time.Second / 4)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// map-based Image cacher
			cs := cacher{}
			s.videoPreviewContext.Range(func(key, value interface{}) bool {
				ctxStr := key.(string)
				pi := value.(VideoPreviewPI)
				sctx := sdcontext.WithContext(ctx, ctxStr)

				// Hostなどが指定されていない場合止める
				if pi.Host == "" || pi.Name == "" {
					return true
				}

				doCrop := true
				anchor := image.Point{}
				var img image.Image
				const buttonSize = 144

				cache, ok := cs.getCache(cacherKey{
					host: pi.Host,
					name: pi.Name,
				})
				if ok {
					// Use cache
					img = cache
				} else {
					// Fetch new image
					c, err := newtek.NewClientV1(pi.Host, pi.User, pi.Password)
					if err != nil {
						msg := fmt.Sprintf("Failed to connect NewTek Client: %s", err)
						s.sd.LogMessage(msg)
						s.sd.ShowAlert(sctx)
						return true
					}
					img, err = c.VideoPreview(pi.Name, buttonSize*5, buttonSize*3, 25)
					if err != nil {
						msg := fmt.Sprintf("Failed to send Shortcut to NewTek Client: %s", err)
						s.sd.LogMessage(msg)
						s.sd.ShowAlert(sctx)
						return true
					}

					cs.storeCache(cacherKey{
						host: pi.Host,
						name: pi.Name,
					}, img)
				}

				switch pi.Segment {
				case 1:
					anchor = image.Point{X: buttonSize * 0, Y: buttonSize * 0}
				case 2:
					anchor = image.Point{X: buttonSize * 1, Y: buttonSize * 0}
				case 3:
					anchor = image.Point{X: buttonSize * 2, Y: buttonSize * 0}
				case 4:
					anchor = image.Point{X: buttonSize * 3, Y: buttonSize * 0}
				case 5:
					anchor = image.Point{X: buttonSize * 4, Y: buttonSize * 0}
				case 6:
					anchor = image.Point{X: buttonSize * 0, Y: buttonSize * 1}
				case 7:
					anchor = image.Point{X: buttonSize * 1, Y: buttonSize * 1}
				case 8:
					anchor = image.Point{X: buttonSize * 2, Y: buttonSize * 1}
				case 9:
					anchor = image.Point{X: buttonSize * 3, Y: buttonSize * 1}
				case 10:
					anchor = image.Point{X: buttonSize * 4, Y: buttonSize * 1}
				case 11:
					anchor = image.Point{X: buttonSize * 0, Y: buttonSize * 2}
				case 12:
					anchor = image.Point{X: buttonSize * 1, Y: buttonSize * 2}
				case 13:
					anchor = image.Point{X: buttonSize * 2, Y: buttonSize * 2}
				case 14:
					anchor = image.Point{X: buttonSize * 3, Y: buttonSize * 2}
				case 15:
					anchor = image.Point{X: buttonSize * 4, Y: buttonSize * 2}
				default:
					doCrop = false
				}

				// cropする場合新規取得
				if doCrop {
					var err error
					img, err = cutter.Crop(img, cutter.Config{
						Width:  buttonSize,
						Height: buttonSize,
						Anchor: anchor,
						Mode:   cutter.TopLeft,
					})
					if err != nil {
						msg := fmt.Sprintf("Failed to send Shortcut to NewTek Client: %s", err)
						s.sd.LogMessage(msg)
						s.sd.ShowAlert(sctx)
						return true
					}
				}

				i, err := streamdeck.Image(img)
				if err != nil {
					msg := fmt.Sprintf("Failed to convert preview image for NewTek Client: %s", err)
					s.sd.LogMessage(msg)
					s.sd.ShowAlert(sctx)
					return true
				}
				go s.sd.SetImage(sctx, i, streamdeck.HardwareAndSoftware)
				return true
			})
		}
	}
}
