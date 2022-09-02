package discordgoself

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	sync.RWMutex
	Token                  string
	MFA                    bool
	LogLevel               int
	ShouldReconnectOnError bool
	Identify               Identify
	Compress               bool
	StateEnabled           bool
	SyncEvents             bool
	DataReady              bool
	MaxRestRetries         int
	VoiceReady             bool
	UDPReady               bool
	VoiceConnections       map[string]*VoiceConnection
	State                  *State
	Client                 *http.Client
	UserAgent              string
	LastHeartbeatAck       time.Time
	LastHeartbeatSent      time.Time
	Ratelimiter            *RateLimiter
	handlersMu             sync.RWMutex
	handlers               map[string][]*eventHandlerInstance
	onceHandlers           map[string][]*eventHandlerInstance
	wsConn                 *websocket.Conn
	listening              chan interface{}
	sequence               *int64
	gateway                string
	sessionID              string
	wsMutex                sync.Mutex
}

type UserConnection struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Revoked      bool           `json:"revoked"`
	Integrations []*Integration `json:"integrations"`
}

type Integration struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Type              string             `json:"type"`
	Enabled           bool               `json:"enabled"`
	Syncing           bool               `json:"syncing"`
	RoleID            string             `json:"role_id"`
	EnableEmoticons   bool               `json:"enable_emoticons"`
	ExpireBehavior    ExpireBehavior     `json:"expire_behavior"`
	ExpireGracePeriod int                `json:"expire_grace_period"`
	User              *User              `json:"user"`
	Account           IntegrationAccount `json:"account"`
	SyncedAt          Timestamp          `json:"synced_at"`
}

type ExpireBehavior int

const (
	ExpireBehaviorRemoveRole ExpireBehavior = 0
	ExpireBehaviorKick       ExpireBehavior = 1
)

type IntegrationAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VoiceRegion struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"sample_hostname"`
	Port     int    `json:"sample_port"`
}

type VoiceICE struct {
	TTL     string       `json:"ttl"`
	Servers []*ICEServer `json:"servers"`
}

type ICEServer struct {
	URL        string `json:"url"`
	Username   string `json:"username"`
	Credential string `json:"credential"`
}

type Note struct {
	UserID     string `json:"user_id"`
	NoteUserID string `json:"note_user_id"`
	Message    string `json:"note"`
}

type Invite struct {
	Guild                    *Guild         `json:"guild"`
	Channel                  *Channel       `json:"channel"`
	Inviter                  *User          `json:"inviter"`
	Code                     string         `json:"code"`
	CreatedAt                Timestamp      `json:"created_at"`
	MaxAge                   int            `json:"max_age"`
	Uses                     int            `json:"uses"`
	MaxUses                  int            `json:"max_uses"`
	Revoked                  bool           `json:"revoked"`
	Temporary                bool           `json:"temporary"`
	Unique                   bool           `json:"unique"`
	TargetUser               *User          `json:"target_user"`
	TargetUserType           TargetUserType `json:"target_user_type"`
	ApproximatePresenceCount int            `json:"approximate_presence_count"`
	ApproximateMemberCount   int            `json:"approximate_member_count"`
}

type TargetUserType int

const (
	TargetUserTypeStream TargetUserType = 1
)

type ChannelType int

const (
	ChannelTypeGuildText     ChannelType = 0
	ChannelTypeDM            ChannelType = 1
	ChannelTypeGuildVoice    ChannelType = 2
	ChannelTypeGroupDM       ChannelType = 3
	ChannelTypeGuildCategory ChannelType = 4
	ChannelTypeGuildNews     ChannelType = 5
	ChannelTypeGuildStore    ChannelType = 6
)

type VideoQualityMode int

const (
	Auto VideoQualityMode = 1
	Full VideoQualityMode = 2
)

type Channel struct {
	ID                   string                 `json:"id"`
	GuildID              string                 `json:"guild_id"`
	Name                 string                 `json:"name"`
	Topic                string                 `json:"topic"`
	Type                 ChannelType            `json:"type"`
	LastMessageID        string                 `json:"last_message_id"`
	LastPinTimestamp     Timestamp              `json:"last_pin_timestamp"`
	NSFW                 bool                   `json:"nsfw"`
	Icon                 string                 `json:"icon"`
	Position             int                    `json:"position"`
	Bitrate              int                    `json:"bitrate"`
	Recipients           []*User                `json:"recipients"`
	Messages             []*Message             `json:"-"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites"`
	UserLimit            int                    `json:"user_limit"`
	ParentID             string                 `json:"parent_id"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user"`
	OwnerID              string                 `json:"owner_id"`
	ApplicationID        string                 `json:"application_id"`
	RTCRegion            string                 `json:"rtc_region"`
	VideoQualityMode     VideoQualityMode       `json:"video_quality_mode"`
}

func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%s>", c.ID)
}

type ChannelEdit struct {
	Name                 string                 `json:"name,omitempty"`
	Topic                string                 `json:"topic,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
	Position             int                    `json:"position"`
	Bitrate              int                    `json:"bitrate,omitempty"`
	UserLimit            int                    `json:"user_limit,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
	VideoQualityMode     VideoQualityMode       `json:"video_quality_mode,omitempty"`
}

type ChannelFollow struct {
	ChannelID string `json:"channel_id"`
	WebhookID string `json:"webhook_id"`
}

type PermissionOverwriteType int

const (
	PermissionOverwriteTypeRole   PermissionOverwriteType = 0
	PermissionOverwriteTypeMember PermissionOverwriteType = 1
)

type PermissionOverwrite struct {
	ID    string                  `json:"id"`
	Type  PermissionOverwriteType `json:"type"`
	Deny  int64                   `json:"deny,string"`
	Allow int64                   `json:"allow,string"`
}

type Emoji struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Roles         []string `json:"roles"`
	User          *User    `json:"user"`
	RequireColons bool     `json:"require_colons"`
	Managed       bool     `json:"managed"`
	Animated      bool     `json:"animated"`
	Available     bool     `json:"available"`
}

var (
	EmojiRegex = regexp.MustCompile(`<(a|):[A-z0-9_~]+:[0-9]{18}>`)
)

func (e *Emoji) MessageFormat() string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + e.APIName() + ">"
		}

		return "<:" + e.APIName() + ">"
	}

	return e.APIName()
}

func (e *Emoji) APIName() string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

type WelcomeChannel struct {
	ChannelID   string `json:"channel_id"`
    Description string `json:"description"`
    EmojiID     string `json:"emoji_id"`
    EmojiName   string `json:"emoji_name"`
}

type WelcomeScreen struct {
	Description     string            `json:"description"`
	WelcomeChannels []*WelcomeChannel `json:"name"`
}

type VerificationLevel int

const (
	VerificationLevelNone     VerificationLevel = 0
	VerificationLevelLow      VerificationLevel = 1
	VerificationLevelMedium   VerificationLevel = 2
	VerificationLevelHigh     VerificationLevel = 3
	VerificationLevelVeryHigh VerificationLevel = 4
)

type ExplicitContentFilterLevel int

const (
	ExplicitContentFilterDisabled            ExplicitContentFilterLevel = 0
	ExplicitContentFilterMembersWithoutRoles ExplicitContentFilterLevel = 1
	ExplicitContentFilterAllMembers          ExplicitContentFilterLevel = 2
)

type MFALevel int

const (
	MFALevelNone     MFALevel = 0
	MFALevelElevated MFALevel = 1
)

type PremiumTier int

const (
	PremiumTierNone PremiumTier = 0
	PremiumTier1    PremiumTier = 1
	PremiumTier2    PremiumTier = 2
	PremiumTier3    PremiumTier = 3
)

type NSFWLevel int

const (
	NSFWLevelNone  NSFWLevel = 0
	NSFWLevelTier1 NSFWLevel = 1
	NSFWLevelTier2 NSFWLevel = 2
	NSFWLevelTier3 NSFWLevel = 3
)

type Guild struct {
	ID                          string                     `json:"id"`
	Name                        string                     `json:"name"`
	Icon                        string                     `json:"icon"`
	IconHash                    string                     `json:"icon_hash"`
	Splash                      string                     `json:"splash"`
	DiscoverySplash             string                     `json:"discovery_splash"`
	Owner                       bool                       `json:"owner"`
	OwnerID                     string                     `json:"owner_id"`
	Permissions                 int64                      `json:"permissions,string"`
	Region                      string                     `json:"region"`
	AfkChannelID                string                     `json:"afk_channel_id"`
	AfkTimeout                  int                        `json:"afk_timeout"`
	WidgetEnabled               bool                       `json:"widget_enabled"`
	WidgetChannelID             string                     `json:"widget_channel_id"`
	VerificationLevel           VerificationLevel          `json:"verification_level"`
	DefaultMessageNotifications MessageNotifications       `json:"default_message_notifications"`
	ExplicitContentFilter       ExplicitContentFilterLevel `json:"explicit_content_filter"`
	Roles                       []*Role                    `json:"roles"`
	Emojis                      []*Emoji                   `json:"emojis"`
	Features                    []string                   `json:"features"`
	MFALevel                    MFALevel                   `json:"mfa_level"`
	ApplicationID               string                     `json:"application_id"`
	SystemChannelID             string                     `json:"system_channel_id"`
	SystemChannelFlags          SystemChannelFlag          `json:"system_channel_flags"`
	RulesChannelID              string                     `json:"rules_channel_id"`
	JoinedAt                    Timestamp                  `json:"joined_at"`
	Large                       bool                       `json:"large"`
	Unavailable                 bool                       `json:"unavailable"`
	MemberCount                 int                        `json:"member_count"`
	VoiceStates                 []*VoiceState              `json:"voice_states"`
	Members                     []*Member                  `json:"members"`
	Channels                    []*Channel                 `json:"channels"`
	Presences                   []*Presence                `json:"presences"`
	MaxPresences                int                        `json:"max_presences"`
	MaxMembers                  int                        `json:"max_members"`
	VanityURLCode               string                     `json:"vanity_url_code"`
	Description                 string                     `json:"description"`
	Banner                      string                     `json:"banner"`
	PremiumTier                 PremiumTier                `json:"premium_tier"`
	PremiumSubscriptionCount    int                        `json:"premium_subscription_count"`
	PreferredLocale             string                     `json:"preferred_locale"`
	PublicUpdatesChannelID      string                     `json:"public_updates_channel_id"`
	MaxVideoChannelUsers        int                        `json:"max_video_channel_users"`
	ApproximateMemberCount      int                        `json:"approximate_member_count"`
	ApproximatePresenceCount    int                        `json:"approximate_presence_count"`
	WelcomeScreen               WelcomeScreen              `json:"welcome_screen"`
	NSFWLevel                   NSFWLevel                  `json:"nsfw_level"`
	Lazy                        bool                       `json:"lazy"`
}

type GuildPreview struct {
	ID                       string   `json:"id"`
	Name                     string   `json:"name"`
	Icon                     string   `json:"icon"`
	Splash                   string   `json:"splash"`
	DiscoverySplash          string   `json:"discovery_splash"`
	Emojis                   []*Emoji `json:"emojis"`
	Features                 []string `json:"features"`
	ApproximateMemberCount   int      `json:"approximate_member_count"`
	ApproximatePresenceCount int      `json:"approximate_presence_count"`
	Description              string   `json:"description"`
}

type MessageNotifications int

const (
	MessageNotificationsAllMessages  MessageNotifications = 0
	MessageNotificationsOnlyMentions MessageNotifications = 1
)

type SystemChannelFlag int

const (
	SystemChannelFlagsSuppressJoin    SystemChannelFlag = 1 << 0
	SystemChannelFlagsSuppressPremium SystemChannelFlag = 1 << 1
)

func (g *Guild) IconURL() string {
	if g.Icon == "" {
		return ""
	}

	if strings.HasPrefix(g.Icon, "a_") {
		return EndpointGuildIconAnimated(g.ID, g.Icon)
	}

	return EndpointGuildIcon(g.ID, g.Icon)
}

type UserGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int64  `json:"permissions,string"`
}

type GuildParams struct {
	Name                        string                     `json:"name,omitempty"`
	Icon                        string                     `json:"icon,omitempty"`
	Splash                      string                     `json:"splash,omitempty"`
	Region                      string                     `json:"region,omitempty"`
	AfkChannelID                string                     `json:"afk_channel_id,omitempty"`
	AfkTimeout                  int                        `json:"afk_timeout,omitempty"`
	WidgetEnabled               bool                       `json:"widget_enabled,omitempty"`
	WidgetChannelID             string                     `json:"widget_channel_id,omitempty"`
	VerificationLevel           *VerificationLevel         `json:"verification_level,omitempty"`
	DefaultMessageNotifications int                        `json:"default_message_notifications,omitempty"`
	ExplicitContentFilter       int                        `json:"explicit_content_filter,omitempty"`
	MFALevel                    int                        `json:"mfa_level,omitempty"`
	SystemChannelID             string                     `json:"system_channel_id,omitempty"`
	RulesChannelID              string                     `json:"rules_channel_id,omitempty"`
	VanityURLCode               string                     `json:"vanity_url_code,omitempty"`
	Description                 string                     `json:"description,omitempty"`
	Banner                      string                     `json:"banner,omitempty"`
	PreferredLocale             string                     `json:"preferred_locale,omitempty"`
	PublicUpdatesChannelID      string                     `json:"public_updates_channel_id,omitempty"`
	WelcomeScreen               *WelcomeScreen             `json:"welcome_screen,omitempty"`
	NSFWLevel                   int                        `json:"nsfw_level,omitempty"`
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`
	Hoist       bool   `json:"hoist"`
	Color       int    `json:"color"`
	Position    int    `json:"position"`
	Permissions int64  `json:"permissions,string"`
}

func (r *Role) Mention() string {
	return fmt.Sprintf("<@&%s>", r.ID)
}

type Roles []*Role

func (r Roles) Len() int {
	return len(r)
}

func (r Roles) Less(i, j int) bool {
	return r[i].Position > r[j].Position
}

func (r Roles) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type VoiceState struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Suppress  bool   `json:"suppress"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	Mute      bool   `json:"mute"`
	Deaf      bool   `json:"deaf"`
}

type Presence struct {
	User       *User       `json:"user"`
	Status     Status      `json:"status"`
	Activities []*Activity `json:"activities"`
	Since      *int        `json:"since"`
}

type TimeStamps struct {
	EndTimestamp   int64 `json:"end,omitempty"`
	StartTimestamp int64 `json:"start,omitempty"`
}

func (t *TimeStamps) UnmarshalJSON(b []byte) error {
	temp := struct {
		End   float64 `json:"end,omitempty"`
		Start float64 `json:"start,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	t.EndTimestamp = int64(temp.End)
	t.StartTimestamp = int64(temp.Start)
	return nil
}

type Assets struct {
	LargeImageID string `json:"large_image,omitempty"`
	SmallImageID string `json:"small_image,omitempty"`
	LargeText    string `json:"large_text,omitempty"`
	SmallText    string `json:"small_text,omitempty"`
}

type Member struct {
	GuildID      string    `json:"guild_id"`
	JoinedAt     Timestamp `json:"joined_at"`
	Nick         string    `json:"nick"`
	Deaf         bool      `json:"deaf"`
	Mute         bool      `json:"mute"`
	User         *User     `json:"user"`
	Roles        []string  `json:"roles"`
	PremiumSince Timestamp `json:"premium_since"`
	Pending      bool      `json:"pending"`
	Permissions  int64     `json:"permissions,string"`
}

func (m *Member) Mention() string {
	return "<@!" + m.User.ID + ">"
}

type Settings struct {
	AfkTimeout                    int               `json:"afk_timeout"`
	AllowAccessibilityDetection   bool              `json:"allow_accessibility_detection"`
	AnimateEmoji                  bool              `json:"animate_emoji"`
	AnimateStickers               int               `json:"animate_stickers"`
	ContactSyncEnabled            bool              `json:"contact_sync_enabled"`
	ConvertEmoticons              bool              `json:"convert_emoticons"`
	CustomStatus                  interface{}       `json:"custom_status"`
	DefaultGuildsRestricted       bool              `json:"default_guilds_restricted"`
	DetectPlatformAccounts        bool              `json:"detect_platform_accounts"`
	DeveloperMode                 bool              `json:"developer_mode"`
	DisableGamesTab               bool              `json:"disable_games_tab"`
	EnableTTSCommand              bool              `json:"enable_tts_command"`
	ExplicitContentFilter         int               `json:"explicit_content_filter"`
	FriendDiscoveryFlags          int               `json:"friend_discovery_flags"`
	FriendSourceFlags             FriendSourceFlags `json:"friend_source_flags"`
	GifAutoPlay                   bool              `json:"gif_auto_play"`
	GuildFolders                  []GuildFolder     `json:"guild_folders"`
	GuildPositions                []string          `json:"guild_positions"`
	InlineAttachmentMedia         bool              `json:"inline_attachment_media"`
	InlineEmbedMedia              bool              `json:"inline_embed_media"`
	Locale                        string            `json:"locale"`
	MessageDisplayCompact         bool              `json:"message_display_compact"`
	NativePhoneIntegrationEnabled bool              `json:"native_phone_integration_enabled"`
	RenderEmbeds                  bool              `json:"render_embeds"`
	RenderReactions               bool              `json:"render_reactions"`
	RestrictedGuilds              []string          `json:"restricted_guilds"`
	ShowCurrentGame               bool              `json:"show_current_game"`
	Status                        string            `json:"status"`
	StreamNotificationsEnabled    bool              `json:"stream_notifications_enabled"`
	Theme                         string            `json:"theme"`
	TimezoneOffset                int               `json:"timezone_offset"`
	ViewNsfwGuilds                bool              `json:"view_nsfw_guilds"`
}

type Status string

const (
	StatusOnline       Status = "online"
	StatusIdle         Status = "idle"
	StatusDoNotDisturb Status = "dnd"
	StatusInvisible    Status = "invisible"
	StatusOffline      Status = "offline"
)

type MutualGuilds struct {
	ID   string `json:"id"`
	Nick string `json:"nick"`
}

type GuildFolder struct {
	ID       string   `json:"name"`
	GuildIDs []string `json:"guild_ids"`
	Name     string   `json:"Name"`
	Color    int      `json:"color"`
}

type FriendSourceFlags struct {
	All           bool `json:"all"`
	MutualGuilds  bool `json:"mutual_guilds"`
	MutualFriends bool `json:"mutual_friends"`
}

type Relationship struct {
	User *User  `json:"user"`
	Type int    `json:"type"`
	ID   string `json:"id"`
}

type TooManyRequests struct {
	Bucket     string        `json:"bucket"`
	Message    string        `json:"message"`
	RetryAfter time.Duration `json:"retry_after"`
}

func (t *TooManyRequests) UnmarshalJSON(b []byte) error {
	u := struct {
		Bucket     string  `json:"bucket"`
		Message    string  `json:"message"`
		RetryAfter float64 `json:"retry_after"`
	}{}
	err := json.Unmarshal(b, &u)
	if err != nil {
		return err
	}

	t.Bucket = u.Bucket
	t.Message = u.Message
	whole, frac := math.Modf(u.RetryAfter)
	t.RetryAfter = time.Duration(whole)*time.Second + time.Duration(frac*1000)*time.Millisecond
	return nil
}

type ReadState struct {
	MentionCount  int    `json:"mention_count"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

type Ack struct {
	Token string `json:"token"`
}

type GuildRole struct {
	Role    *Role  `json:"role"`
	GuildID string `json:"guild_id"`
}

type GuildBan struct {
	Reason string `json:"reason"`
	User   *User  `json:"user"`
}

type GuildEmbed struct {
	Enabled   bool   `json:"enabled"`
	ChannelID string `json:"channel_id"`
}

type GuildAuditLog struct {
	Webhooks        []*Webhook       `json:"webhooks,omitempty"`
	Users           []*User          `json:"users,omitempty"`
	AuditLogEntries []*AuditLogEntry `json:"audit_log_entries"`
	Integrations    []*Integration   `json:"integrations"`
}

type AuditLogEntry struct {
	TargetID   string            `json:"target_id"`
	Changes    []*AuditLogChange `json:"changes"`
	UserID     string            `json:"user_id"`
	ID         string            `json:"id"`
	ActionType *AuditLogAction   `json:"action_type"`
	Options    *AuditLogOptions  `json:"options"`
	Reason     string            `json:"reason"`
}

type AuditLogChange struct {
	NewValue interface{}        `json:"new_value"`
	OldValue interface{}        `json:"old_value"`
	Key      *AuditLogChangeKey `json:"key"`
}

type AuditLogChangeKey string

const (
	AuditLogChangeKeyName                       AuditLogChangeKey = "name"
	AuditLogChangeKeyIconHash                   AuditLogChangeKey = "icon_hash"
	AuditLogChangeKeySplashHash                 AuditLogChangeKey = "splash_hash"
	AuditLogChangeKeyOwnerID                    AuditLogChangeKey = "owner_id"
	AuditLogChangeKeyRegion                     AuditLogChangeKey = "region"
	AuditLogChangeKeyAfkChannelID               AuditLogChangeKey = "afk_channel_id"
	AuditLogChangeKeyAfkTimeout                 AuditLogChangeKey = "afk_timeout"
	AuditLogChangeKeyMFALevel                   AuditLogChangeKey = "mfa_level"
	AuditLogChangeKeyVerificationLevel          AuditLogChangeKey = "verification_level"
	AuditLogChangeKeyExplicitContentFilter      AuditLogChangeKey = "explicit_content_filter"
	AuditLogChangeKeyDefaultMessageNotification AuditLogChangeKey = "default_message_notifications"
	AuditLogChangeKeyVanityURLCode              AuditLogChangeKey = "vanity_url_code"
	AuditLogChangeKeyRoleAdd                    AuditLogChangeKey = "$add"
	AuditLogChangeKeyRoleRemove                 AuditLogChangeKey = "$remove"
	AuditLogChangeKeyPruneDeleteDays            AuditLogChangeKey = "prune_delete_days"
	AuditLogChangeKeyWidgetEnabled              AuditLogChangeKey = "widget_enabled"
	AuditLogChangeKeyWidgetChannelID            AuditLogChangeKey = "widget_channel_id"
	AuditLogChangeKeySystemChannelID            AuditLogChangeKey = "system_channel_id"
	AuditLogChangeKeyPosition                   AuditLogChangeKey = "position"
	AuditLogChangeKeyTopic                      AuditLogChangeKey = "topic"
	AuditLogChangeKeyBitrate                    AuditLogChangeKey = "bitrate"
	AuditLogChangeKeyPermissionOverwrite        AuditLogChangeKey = "permission_overwrites"
	AuditLogChangeKeyNSFW                       AuditLogChangeKey = "nsfw"
	AuditLogChangeKeyApplicationID              AuditLogChangeKey = "application_id"
	AuditLogChangeKeyRateLimitPerUser           AuditLogChangeKey = "rate_limit_per_user"
	AuditLogChangeKeyPermissions                AuditLogChangeKey = "permissions"
	AuditLogChangeKeyColor                      AuditLogChangeKey = "color"
	AuditLogChangeKeyHoist                      AuditLogChangeKey = "hoist"
	AuditLogChangeKeyMentionable                AuditLogChangeKey = "mentionable"
	AuditLogChangeKeyAllow                      AuditLogChangeKey = "allow"
	AuditLogChangeKeyDeny                       AuditLogChangeKey = "deny"
	AuditLogChangeKeyCode                       AuditLogChangeKey = "code"
	AuditLogChangeKeyChannelID                  AuditLogChangeKey = "channel_id"
	AuditLogChangeKeyInviterID                  AuditLogChangeKey = "inviter_id"
	AuditLogChangeKeyMaxUses                    AuditLogChangeKey = "max_uses"
	AuditLogChangeKeyUses                       AuditLogChangeKey = "uses"
	AuditLogChangeKeyMaxAge                     AuditLogChangeKey = "max_age"
	AuditLogChangeKeyTempoary                   AuditLogChangeKey = "temporary"
	AuditLogChangeKeyDeaf                       AuditLogChangeKey = "deaf"
	AuditLogChangeKeyMute                       AuditLogChangeKey = "mute"
	AuditLogChangeKeyNick                       AuditLogChangeKey = "nick"
	AuditLogChangeKeyAvatarHash                 AuditLogChangeKey = "avatar_hash"
	AuditLogChangeKeyID                         AuditLogChangeKey = "id"
	AuditLogChangeKeyType                       AuditLogChangeKey = "type"
	AuditLogChangeKeyEnableEmoticons            AuditLogChangeKey = "enable_emoticons"
	AuditLogChangeKeyExpireBehavior             AuditLogChangeKey = "expire_behavior"
	AuditLogChangeKeyExpireGracePeriod          AuditLogChangeKey = "expire_grace_period"
)

type AuditLogOptions struct {
	DeleteMemberDays string               `json:"delete_member_days"`
	MembersRemoved   string               `json:"members_removed"`
	ChannelID        string               `json:"channel_id"`
	MessageID        string               `json:"message_id"`
	Count            string               `json:"count"`
	ID               string               `json:"id"`
	Type             *AuditLogOptionsType `json:"type"`
	RoleName         string               `json:"role_name"`
}

type AuditLogOptionsType string

const (
	AuditLogOptionsTypeMember AuditLogOptionsType = "member"
	AuditLogOptionsTypeRole   AuditLogOptionsType = "role"
)

type AuditLogAction int

const (
	AuditLogActionGuildUpdate            AuditLogAction = 1
	AuditLogActionChannelCreate          AuditLogAction = 10
	AuditLogActionChannelUpdate          AuditLogAction = 11
	AuditLogActionChannelDelete          AuditLogAction = 12
	AuditLogActionChannelOverwriteCreate AuditLogAction = 13
	AuditLogActionChannelOverwriteUpdate AuditLogAction = 14
	AuditLogActionChannelOverwriteDelete AuditLogAction = 15
	AuditLogActionMemberKick             AuditLogAction = 20
	AuditLogActionMemberPrune            AuditLogAction = 21
	AuditLogActionMemberBanAdd           AuditLogAction = 22
	AuditLogActionMemberBanRemove        AuditLogAction = 23
	AuditLogActionMemberUpdate           AuditLogAction = 24
	AuditLogActionMemberRoleUpdate       AuditLogAction = 25
	AuditLogActionRoleCreate             AuditLogAction = 30
	AuditLogActionRoleUpdate             AuditLogAction = 31
	AuditLogActionRoleDelete             AuditLogAction = 32
	AuditLogActionInviteCreate           AuditLogAction = 40
	AuditLogActionInviteUpdate           AuditLogAction = 41
	AuditLogActionInviteDelete           AuditLogAction = 42
	AuditLogActionWebhookCreate          AuditLogAction = 50
	AuditLogActionWebhookUpdate          AuditLogAction = 51
	AuditLogActionWebhookDelete          AuditLogAction = 52
	AuditLogActionEmojiCreate            AuditLogAction = 60
	AuditLogActionEmojiUpdate            AuditLogAction = 61
	AuditLogActionEmojiDelete            AuditLogAction = 62
	AuditLogActionMessageDelete          AuditLogAction = 72
	AuditLogActionMessageBulkDelete      AuditLogAction = 73
	AuditLogActionMessagePin             AuditLogAction = 74
	AuditLogActionMessageUnpin           AuditLogAction = 75
	AuditLogActionIntegrationCreate      AuditLogAction = 80
	AuditLogActionIntegrationUpdate      AuditLogAction = 81
	AuditLogActionIntegrationDelete      AuditLogAction = 82
)

type UserGuildSettingsChannelOverride struct {
	Muted                bool            `json:"muted"`
	MuteConfig           *UserMuteConfig `json:"mute_config"`
	MessageNotifications int             `json:"message_notifications"`
	ChannelID            string          `json:"channel_id"`
}

type UserMuteConfig struct {
	SelectedTimeWindow int       `json:"selected_time_window"`
	EndTime            time.Time `json:"end_time"`
}

type UserGuildSettings struct {
	SupressEveryone      bool                                `json:"suppress_everyone"`
	SupressRoles         bool                                `json:"suppress_roles"`
	Muted                bool                                `json:"muted"`
	MuteConfig           *UserMuteConfig                     `json:"mute_config"`
	MobilePush           bool                                `json:"mobile_push"`
	MessageNotifications int                                 `json:"message_notifications"`
	GuildID              string                              `json:"guild_id"`
	ChannelOverrides     []*UserGuildSettingsChannelOverride `json:"channel_overrides"`
}

type UserGuildSettingsEdit struct {
	SupressEveryone      bool                                         `json:"suppress_everyone"`
	SupressRoles         bool                                         `json:"suppress_roles"`
	Muted                bool                                         `json:"muted"`
	MuteConfig           *UserMuteConfig                              `json:"mute_config"`
	MobilePush           bool                                         `json:"mobile_push"`
	MessageNotifications int                                          `json:"message_notifications"`
	ChannelOverrides     map[string]*UserGuildSettingsChannelOverride `json:"channel_overrides"`
}

type APIErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MessageReaction struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     Emoji  `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

type GatewayStatusUpdate struct {
	Since  int      `json:"since"`
	Game   Activity `json:"game"`
	Status string   `json:"status"`
	AFK    bool     `json:"afk"`
}

type Activity struct {
	Name          string       `json:"name"`
	Type          ActivityType `json:"type"`
	URL           string       `json:"url,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	ApplicationID string       `json:"application_id,omitempty"`
	State         string       `json:"state,omitempty"`
	Details       string       `json:"details,omitempty"`
	Timestamps    TimeStamps   `json:"timestamps,omitempty"`
	Emoji         Emoji        `json:"emoji,omitempty"`
	Party         Party        `json:"party,omitempty"`
	Assets        Assets       `json:"assets,omitempty"`
	Secrets       Secrets      `json:"secrets,omitempty"`
	Instance      bool         `json:"instance,omitempty"`
	Flags         int          `json:"flags,omitempty"`
}

func (activity *Activity) UnmarshalJSON(b []byte) error {
	temp := struct {
		Name          string       `json:"name"`
		Type          ActivityType `json:"type"`
		URL           string       `json:"url,omitempty"`
		CreatedAt     int64        `json:"created_at"`
		ApplicationID string       `json:"application_id,omitempty"`
		State         string       `json:"state,omitempty"`
		Details       string       `json:"details,omitempty"`
		Timestamps    TimeStamps   `json:"timestamps,omitempty"`
		Emoji         Emoji        `json:"emoji,omitempty"`
		Party         Party        `json:"party,omitempty"`
		Assets        Assets       `json:"assets,omitempty"`
		Secrets       Secrets      `json:"secrets,omitempty"`
		Instance      bool         `json:"instance,omitempty"`
		Flags         int          `json:"flags,omitempty"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	activity.CreatedAt = time.Unix(0, temp.CreatedAt*1000000)
	activity.ApplicationID = temp.ApplicationID
	activity.Assets = temp.Assets
	activity.Details = temp.Details
	activity.Emoji = temp.Emoji
	activity.Flags = temp.Flags
	activity.Instance = temp.Instance
	activity.Name = temp.Name
	activity.Party = temp.Party
	activity.Secrets = temp.Secrets
	activity.State = temp.State
	activity.Timestamps = temp.Timestamps
	activity.Type = temp.Type
	activity.URL = temp.URL
	return nil
}

type Party struct {
	ID   string `json:"id,omitempty"`
	Size []int  `json:"size,omitempty"`
}

type Secrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

type ActivityType int

const (
	ActivityTypeGame      ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeCustom    ActivityType = 4
)

type Identify struct {
	Token              string              `json:"token"`
	Properties         IdentifyProperties  `json:"properties"`
	Compress           bool                `json:"compress"`
	LargeThreshold     int                 `json:"large_threshold"`
	Presence           GatewayStatusUpdate `json:"presence,omitempty"`
}

type IdentifyProperties struct {
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referrer        string `json:"$referer"`
	ReferringDomain string `json:"$referring_domain"`
}

const (
	PermissionReadMessages       = 0x0000000000000400
	PermissionSendMessages       = 0x0000000000000800
	PermissionSendTTSMessages    = 0x0000000000001000
	PermissionManageMessages     = 0x0000000000002000
	PermissionEmbedLinks         = 0x0000000000004000
	PermissionAttachFiles        = 0x0000000000008000
	PermissionReadMessageHistory = 0x0000000000010000
	PermissionMentionEveryone    = 0x0000000000020000
	PermissionUseExternalEmojis  = 0x0000000000040000
	PermissionUseSlashCommands   = 0x0000000080000000
)

const (
	PermissionVoicePrioritySpeaker = 0x0000000000000100
	PermissionVoiceStreamVideo     = 0x0000000000000200
	PermissionVoiceConnect         = 0x0000000000100000
	PermissionVoiceSpeak           = 0x0000000000200000
	PermissionVoiceMuteMembers     = 0x0000000000400000
	PermissionVoiceDeafenMembers   = 0x0000000000800000
	PermissionVoiceMoveMembers     = 0x0000000001000000
	PermissionVoiceUseVAD          = 0x0000000002000000
	PermissionVoiceRequestToSpeak  = 0x0000000100000000
)

const (
	PermissionChangeNickname  = 0x0000000004000000
	PermissionManageNicknames = 0x0000000008000000
	PermissionManageRoles     = 0x0000000010000000
	PermissionManageWebhooks  = 0x0000000020000000
	PermissionManageEmojis    = 0x0000000040000000
)

const (
	PermissionCreateInstantInvite = 0x0000000000000001
	PermissionKickMembers         = 0x0000000000000002
	PermissionBanMembers          = 0x0000000000000004
	PermissionAdministrator       = 0x0000000000000008
	PermissionManageChannels      = 0x0000000000000010
	PermissionManageServer        = 0x0000000000000020
	PermissionAddReactions        = 0x0000000000000040
	PermissionViewAuditLogs       = 0x0000000000000080
	PermissionViewChannel         = 0x0000000000000400
	PermissionViewGuildInsights   = 0x0000000000080000

	PermissionAllText = PermissionViewChannel |
		PermissionSendMessages |
		PermissionSendTTSMessages |
		PermissionManageMessages |
		PermissionEmbedLinks |
		PermissionAttachFiles |
		PermissionReadMessageHistory |
		PermissionMentionEveryone
	PermissionAllVoice = PermissionViewChannel |
		PermissionVoiceConnect |
		PermissionVoiceSpeak |
		PermissionVoiceMuteMembers |
		PermissionVoiceDeafenMembers |
		PermissionVoiceMoveMembers |
		PermissionVoiceUseVAD |
		PermissionVoicePrioritySpeaker
	PermissionAllChannel = PermissionAllText |
		PermissionAllVoice |
		PermissionCreateInstantInvite |
		PermissionManageRoles |
		PermissionManageChannels |
		PermissionAddReactions |
		PermissionViewAuditLogs
	PermissionAll = PermissionAllChannel |
		PermissionKickMembers |
		PermissionBanMembers |
		PermissionManageServer |
		PermissionAdministrator |
		PermissionManageWebhooks |
		PermissionManageEmojis
)

const (
	ErrCodeUnknownAccount                            = 10001
	ErrCodeUnknownApplication                        = 10002
	ErrCodeUnknownChannel                            = 10003
	ErrCodeUnknownGuild                              = 10004
	ErrCodeUnknownIntegration                        = 10005
	ErrCodeUnknownInvite                             = 10006
	ErrCodeUnknownMember                             = 10007
	ErrCodeUnknownMessage                            = 10008
	ErrCodeUnknownOverwrite                          = 10009
	ErrCodeUnknownProvider                           = 10010
	ErrCodeUnknownRole                               = 10011
	ErrCodeUnknownToken                              = 10012
	ErrCodeUnknownUser                               = 10013
	ErrCodeUnknownEmoji                              = 10014
	ErrCodeUnknownWebhook                            = 10015
	ErrCodeUnknownBan                                = 10026
	ErrCodeBotsCannotUseEndpoint                     = 20001
	ErrCodeOnlyBotsCanUseEndpoint                    = 20002
	ErrCodeMaximumGuildsReached                      = 30001
	ErrCodeMaximumFriendsReached                     = 30002
	ErrCodeMaximumPinsReached                        = 30003
	ErrCodeMaximumGuildRolesReached                  = 30005
	ErrCodeTooManyReactions                          = 30010
	ErrCodeUnauthorized                              = 40001
	ErrCodeMissingAccess                             = 50001
	ErrCodeInvalidAccountType                        = 50002
	ErrCodeCannotExecuteActionOnDMChannel            = 50003
	ErrCodeEmbedDisabled                             = 50004
	ErrCodeCannotEditFromAnotherUser                 = 50005
	ErrCodeCannotSendEmptyMessage                    = 50006
	ErrCodeCannotSendMessagesToThisUser              = 50007
	ErrCodeCannotSendMessagesInVoiceChannel          = 50008
	ErrCodeChannelVerificationLevelTooHigh           = 50009
	ErrCodeOAuth2ApplicationDoesNotHaveBot           = 50010
	ErrCodeOAuth2ApplicationLimitReached             = 50011
	ErrCodeInvalidOAuthState                         = 50012
	ErrCodeMissingPermissions                        = 50013
	ErrCodeInvalidAuthenticationToken                = 50014
	ErrCodeNoteTooLong                               = 50015
	ErrCodeTooFewOrTooManyMessagesToDelete           = 50016
	ErrCodeCanOnlyPinMessageToOriginatingChannel     = 50019
	ErrCodeCannotExecuteActionOnSystemMessage        = 50021
	ErrCodeMessageProvidedTooOldForBulkDelete        = 50034
	ErrCodeInvalidFormBody                           = 50035
	ErrCodeInviteAcceptedToGuildApplicationsBotNotIn = 50036
	ErrCodeReactionBlocked                           = 90001
)
