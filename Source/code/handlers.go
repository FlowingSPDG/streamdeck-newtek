package sdnewtek

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FlowingSPDG/newtek-go"
	"github.com/FlowingSPDG/streamdeck"
)

func (s *SDNewTek) ShortcutWillAppearHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.WillAppearPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal WillAppear event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	settings := ShortcutPI{}
	if err := json.Unmarshal(payload.Settings, &settings); err != nil {
		// エラー表示
		msg := fmt.Sprintf("Failed to unmarshal WillAppear event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)

		// パースエラーが起きた時点で初期値に直す
		settings.Initialize()
		client.SetSettings(ctx, settings)
	}

	if settings.IsDefault() {
		settings.Initialize()
		client.SetSettings(ctx, settings)
	}

	msg := fmt.Sprintf("WillAppear with settings :%v", settings)
	client.LogMessage(msg)
	return nil
}

// KeyDownHandler keyDown handler
func (s *SDNewTek) ShortcutKeyDownHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.KeyDownPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	pi := ShortcutPI{}
	if err := json.Unmarshal(payload.Settings, &pi); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown settings payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	c, err := newtek.NewClientV1(pi.Host, pi.User, pi.Password)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	kv := make(map[string]string)
	if pi.Key != "" {
		kv[pi.Key] = pi.Value
	}
	msg := fmt.Sprintf("Sending shortcut command %v", pi)
	client.LogMessage(msg)
	if err := c.ShortcutHTTP(pi.Shortcut, kv); err != nil {
		msg := fmt.Sprintf("Failed to send Shortcut to NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	msg = fmt.Sprintf("Sent command %v", pi)
	client.LogMessage(msg)

	return client.ShowOk(ctx)
}

func (s *SDNewTek) VideoPreviewWillAppearHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.WillAppearPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal WillAppear event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	settings := VideoPreviewPI{}
	if err := json.Unmarshal(payload.Settings, &settings); err != nil {
		// エラー表示
		msg := fmt.Sprintf("Failed to unmarshal WillAppear event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)

		// パースエラーが起きた時点で初期値に直す
		settings.Initialize()
		client.SetSettings(ctx, settings)
	}

	if settings.IsDefault() {
		settings.Initialize()
		client.SetSettings(ctx, settings)
	}

	msg := fmt.Sprintf("WillAppear with settings :%v", settings)
	client.LogMessage(msg)
	return nil
}

// KeyDownHandler keyDown handler
func (s *SDNewTek) VideoPreviewKeyDownHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.KeyDownPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	pi := VideoPreviewPI{}
	if err := json.Unmarshal(payload.Settings, &pi); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown settings payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	c, err := newtek.NewClientV1(pi.Host, pi.User, pi.Password)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	msg := fmt.Sprintf("Sending VideoPreview command %v", pi)
	client.LogMessage(msg)
	img, err := c.VideoPreview(pi.Name, 144, 144, 25)
	if err != nil {
		msg := fmt.Sprintf("Failed to send Shortcut to NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}
	i, _ := streamdeck.Image(img)
	client.SetImage(ctx, i, streamdeck.HardwareAndSoftware)

	msg = fmt.Sprintf("Set image %v", pi)
	client.LogMessage(msg)
	return nil
}
