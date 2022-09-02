package discordgoself

type Webhook struct {
	ID            string      `json:"id"`
	Type          WebhookType `json:"type"`
	GuildID       string      `json:"guild_id"`
	ChannelID     string      `json:"channel_id"`
	User          *User       `json:"user"`
	Name          string      `json:"name"`
	Avatar        string      `json:"avatar"`
	Token         string      `json:"token"`
	ApplicationID string      `json:"application_id,omitempty"`
}

type WebhookType int

const (
	WebhookTypeIncoming        WebhookType = 1
	WebhookTypeChannelFollower WebhookType = 2
)

type WebhookParams struct {
	Content         string                  `json:"content,omitempty"`
	Username        string                  `json:"username,omitempty"`
	AvatarURL       string                  `json:"avatar_url,omitempty"`
	TTS             bool                    `json:"tts,omitempty"`
	Files           []*File                 `json:"-"`
	Components      []MessageComponent      `json:"components"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
}

type WebhookEdit struct {
	Content         string                  `json:"content,omitempty"`
	Components      []MessageComponent      `json:"components"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	Files           []*File                 `json:"-"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
}
