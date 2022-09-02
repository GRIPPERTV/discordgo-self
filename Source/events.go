package discordgoself

import (
	"encoding/json"
)

type Connect struct{}
type Disconnect struct{}

type RateLimit struct {
	*TooManyRequests
	URL string
}

type Event struct {
	Operation int             `json:"op"`
	Sequence  int64           `json:"s"`
	Type      string          `json:"t"`
	RawData   json.RawMessage `json:"d"`
	Struct    interface{}     `json:"-"`
}

type Ready struct {
	Version           int                  `json:"v"`
	SessionID         string               `json:"session_id"`
	User              *User                `json:"user"`
	ReadState         []*ReadState         `json:"read_state"`
	PrivateChannels   []*Channel           `json:"private_channels"`
	Guilds            []*Guild             `json:"guilds"`
	Settings          *Settings            `json:"user_settings"`
	UserGuildSettings []*UserGuildSettings `json:"user_guild_settings"`
	Relationships     []*Relationship      `json:"relationships"`
	Presences         []*Presence          `json:"presences"`
	Notes             map[string]string    `json:"notes"`
}

type ChannelCreate struct {
	*Channel
}

type ChannelUpdate struct {
	*Channel
}

type ChannelDelete struct {
	*Channel
}

type ChannelPinsUpdate struct {
	LastPinTimestamp string `json:"last_pin_timestamp"`
	ChannelID        string `json:"channel_id"`
	GuildID          string `json:"guild_id,omitempty"`
}

type GuildCreate struct {
	*Guild
}

type GuildUpdate struct {
	*Guild
}

type GuildDelete struct {
	*Guild
}

type GuildBanAdd struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

type GuildBanRemove struct {
	User    *User  `json:"user"`
	GuildID string `json:"guild_id"`
}

type GuildMemberAdd struct {
	*Member
}

type GuildMemberUpdate struct {
	*Member
}

type GuildMemberRemove struct {
	*Member
}

type GuildRoleCreate struct {
	*GuildRole
}

type GuildRoleUpdate struct {
	*GuildRole
}

type GuildRoleDelete struct {
	RoleID  string `json:"role_id"`
	GuildID string `json:"guild_id"`
}

type GuildEmojisUpdate struct {
	GuildID string   `json:"guild_id"`
	Emojis  []*Emoji `json:"emojis"`
}

type GuildMembersChunk struct {
	GuildID    string      `json:"guild_id"`
	Members    []*Member   `json:"members"`
	ChunkIndex int         `json:"chunk_index"`
	ChunkCount int         `json:"chunk_count"`
	Presences  []*Presence `json:"presences,omitempty"`
}

type GuildIntegrationsUpdate struct {
	GuildID string `json:"guild_id"`
}

type MessageAck struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
}

type MessageCreate struct {
	*Message
}

func (m *MessageCreate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

type MessageUpdate struct {
	*Message
	BeforeUpdate *Message `json:"-"`
}

func (m *MessageUpdate) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

type MessageDelete struct {
	*Message
	BeforeDelete *Message `json:"-"`
}

func (m *MessageDelete) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.Message)
}

type MessageReactionAdd struct {
	*MessageReaction
}

type MessageReactionRemove struct {
	*MessageReaction
}

type MessageReactionRemoveAll struct {
	*MessageReaction
}

type PresencesReplace []*Presence

type PresenceUpdate struct {
	Presence
	GuildID string `json:"guild_id"`
}

type Resumed struct {
	Trace []string `json:"_trace"`
}

type RelationshipAdd struct {
	*Relationship
}

type RelationshipRemove struct {
	*Relationship
}

type TypingStart struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Timestamp int    `json:"timestamp"`
}

type UserUpdate struct {
	*User
}

type UserSettingsUpdate map[string]interface{}

type UserGuildSettingsUpdate struct {
	*UserGuildSettings
}

type UserNoteUpdate struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}

type VoiceServerUpdate struct {
	Token    string `json:"token"`
	GuildID  string `json:"guild_id"`
	Endpoint string `json:"endpoint"`
}

type VoiceStateUpdate struct {
	*VoiceState
	BeforeUpdate *VoiceState `json:"-"`
}

type MessageDeleteBulk struct {
	Messages  []string `json:"ids"`
	ChannelID string   `json:"channel_id"`
	GuildID   string   `json:"guild_id"`
}

type WebhooksUpdate struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
}
