package discordgoself

import (
	"encoding/json"
)

type ComponentType uint

const (
	ActionsRowComponent ComponentType = 1
	ButtonComponent     ComponentType = 2
)

type MessageComponent interface {
	json.Marshaler
	Type() ComponentType
}

type unmarshalableMessageComponent struct {
	MessageComponent
}

func (umc *unmarshalableMessageComponent) UnmarshalJSON(src []byte) error {
	var v struct {
		Type ComponentType `json:"type"`
	}
	err := json.Unmarshal(src, &v)
	if err != nil {
		return err
	}

	var data MessageComponent
	switch v.Type {
	case ActionsRowComponent:
		v := ActionsRow{}
		err = json.Unmarshal(src, &v)
		data = v
	case ButtonComponent:
		v := Button{}
		err = json.Unmarshal(src, &v)
		data = v
	}
	if err != nil {
		return err
	}
	umc.MessageComponent = data
	return err
}

type ActionsRow struct {
	Components []MessageComponent `json:"components"`
}

func (r ActionsRow) MarshalJSON() ([]byte, error) {
	type actionsRow ActionsRow

	return json.Marshal(struct {
		actionsRow
		Type ComponentType `json:"type"`
	}{
		actionsRow: actionsRow(r),
		Type:       r.Type(),
	})
}

func (r *ActionsRow) UnmarshalJSON(data []byte) error {
	var v struct {
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	r.Components = make([]MessageComponent, len(v.RawComponents))

	for i, v := range v.RawComponents {
		r.Components[i] = v.MessageComponent
	}

	return err
}

func (r ActionsRow) Type() ComponentType {
	return ActionsRowComponent
}

type ButtonStyle uint

const (
	PrimaryButton   ButtonStyle = 1
	SecondaryButton ButtonStyle = 2
	SuccessButton   ButtonStyle = 3
	DangerButton    ButtonStyle = 4
	LinkButton      ButtonStyle = 5
)

type ButtonEmoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}

type Button struct {
	Label    string      `json:"label"`
	Style    ButtonStyle `json:"style"`
	Disabled bool        `json:"disabled"`
	Emoji    ButtonEmoji `json:"emoji"`
	URL      string      `json:"url,omitempty"`
	CustomID string      `json:"custom_id,omitempty"`
}

func (b Button) MarshalJSON() ([]byte, error) {
	type button Button

	if b.Style == 0 {
		b.Style = PrimaryButton
	}

	return json.Marshal(struct {
		button
		Type ComponentType `json:"type"`
	}{
		button: button(b),
		Type:   b.Type(),
	})
}

func (b Button) Type() ComponentType {
	return ButtonComponent
}
