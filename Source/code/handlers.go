package sdnewtek

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FlowingSPDG/newtek-go"
	"github.com/FlowingSPDG/streamdeck"
)

// ShortcutWillAppearHandler WillAppear handler for ShortcutHTTP
func (s *SDNewTek) ShortcutWillAppearHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.WillAppearPayload[ShortcutPI]{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal WillAppear event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	if payload.Settings.IsDefault() {
		payload.Settings.Initialize()
		client.SetSettings(ctx, payload.Settings)
	}

	msg := fmt.Sprintf("Context %s WillAppear with settings :%v", event.Context, payload.Settings)
	client.LogMessage(msg)
	return nil
}

// ShortcutKeyDownHandler keyDown handler
func (s *SDNewTek) ShortcutKeyDownHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.KeyDownPayload[ShortcutPI]{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	c, err := newtek.NewClientV1(payload.Settings.Host, payload.Settings.User, payload.Settings.Password)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	kv := make(map[string]string)
	if payload.Settings.Key != "" {
		kv[payload.Settings.Key] = payload.Settings.Value
	}
	msg := fmt.Sprintf("Sending shortcut command %v", payload.Settings)
	client.LogMessage(msg)
	if err := c.ShortcutHTTP(payload.Settings.Shortcut, kv); err != nil {
		msg := fmt.Sprintf("Failed to send Shortcut(HTTP) to NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	msg = fmt.Sprintf("Sent command %v", payload.Settings)
	client.LogMessage(msg)

	return client.ShowOk(ctx)
}

// ShortcutTCPWillAppearHandler WillAppear handler for ShortcutTCP
func (s *SDNewTek) ShortcutTCPWillAppearHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.WillAppearPayload[ShortcutTCPPI]{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal WillAppear event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	if payload.Settings.IsDefault() {
		payload.Settings.Initialize()
		client.SetSettings(ctx, payload.Settings)
	}

	msg := fmt.Sprintf("Context %s WillAppear with settings :%v", event.Context, payload.Settings)
	client.LogMessage(msg)
	return nil
}

// ShortcutTCPKeyDownHandler keyDown handler
func (s *SDNewTek) ShortcutTCPKeyDownHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.KeyDownPayload[ShortcutTCPPI]{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown event payload: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	c, err := newtek.NewClientV1TCP(payload.Settings.Host)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}
	defer c.Close()

	msg := fmt.Sprintf("Sending shortcut command %v", payload.Settings)
	client.LogMessage(msg)
	if err := c.Shortcut(payload.Settings.ToShortcuts()); err != nil {
		msg := fmt.Sprintf("Failed to send Shortcut(TCP) to NewTek Client: %s", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	msg = fmt.Sprintf("Sent command %v", payload.Settings)
	client.LogMessage(msg)

	return client.ShowOk(ctx)
}
