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
		msg := fmt.Sprintf("Failed to unmarshal WillAppear event payload: %s\n", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}
	// PIが完全に空だった場合初期情報を入れて返したい
	return nil
}

// KeyDownHandler keyDown handler
func (s *SDNewTek) ShortcutKeyDownHandler(ctx context.Context, client *streamdeck.Client, event streamdeck.Event) error {
	payload := streamdeck.KeyDownPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown event payload: %s\n", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	pi := ShortcutPI{}
	if err := json.Unmarshal(payload.Settings, &pi); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal KeyDown settings payload: %s\n", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	c, err := newtek.NewClientV1(pi.Host, pi.User, pi.Password)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect NewTek Client: %s\n", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	if err := c.ShortcutHTTP(pi.Shortcut, pi.KV); err != nil {
		msg := fmt.Sprintf("Failed to send Shortcut to NewTek Client: %s\n", err)
		client.LogMessage(msg)
		client.ShowAlert(ctx)
		return err
	}

	return client.ShowOk(ctx)
}
