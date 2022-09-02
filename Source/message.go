package discordgoself

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"
)

type MessageType int

const (
	MessageTypeDefault                               MessageType = 0
	MessageTypeRecipientAdd                          MessageType = 1
	MessageTypeRecipientRemove                       MessageType = 2
	MessageTypeCall                                  MessageType = 3
	MessageTypeChannelNameChange                     MessageType = 4
	MessageTypeChannelIconChange                     MessageType = 5
	MessageTypeChannelPinnedMessage                  MessageType = 6
	MessageTypeGuildMemberJoin                       MessageType = 7
	MessageTypeUserPremiumGuildSubscription          MessageType = 8
	MessageTypeUserPremiumGuildSubscriptionTierOne   MessageType = 9
	MessageTypeUserPremiumGuildSubscriptionTierTwo   MessageType = 10
	MessageTypeUserPremiumGuildSubscriptionTierThree MessageType = 11
	MessageTypeChannelFollowAdd                      MessageType = 12
	MessageTypeGuildDiscoveryDisqualified            MessageType = 14
	MessageTypeGuildDiscoveryRequalified             MessageType = 15
	MessageTypeReply                                 MessageType = 19
	MessageTypeApplicationCommand                    MessageType = 20
)

type Message struct {
	ID               string               `json:"id"`
	ChannelID        string               `json:"channel_id"`
	GuildID          string               `json:"guild_id,omitempty"`
	Content          string               `json:"content"`
	Timestamp        Timestamp            `json:"timestamp"`
	EditedTimestamp  Timestamp            `json:"edited_timestamp"`
	MentionRoles     []string             `json:"mention_roles"`
	TTS              bool                 `json:"tts"`
	MentionEveryone  bool                 `json:"mention_everyone"`
	Author           *User                `json:"author"`
	Attachments      []*MessageAttachment `json:"attachments"`
	Components       []MessageComponent   `json:"-"`
	Embeds           []*MessageEmbed      `json:"embeds"`
	Mentions         []*User              `json:"mentions"`
	Reactions        []*MessageReactions  `json:"reactions"`
	Pinned           bool                 `json:"pinned"`
	Type             MessageType          `json:"type"`
	WebhookID        string               `json:"webhook_id"`
	Member           *Member              `json:"member"`
	MentionChannels  []*Channel           `json:"mention_channels"`
	Activity         *MessageActivity     `json:"activity"`
	Application      *MessageApplication  `json:"application"`
	MessageReference *MessageReference    `json:"message_reference"`
	Flags            MessageFlags         `json:"flags"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message

	var v struct {
		message
		RawComponents []unmarshalableMessageComponent `json:"components"`
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*m = Message(v.message)
	m.Components = make([]MessageComponent, len(v.RawComponents))

	for i, v := range v.RawComponents {
		m.Components[i] = v.MessageComponent
	}

	return err
}

func (m *Message) GetCustomEmojis() []*Emoji {
	var toReturn []*Emoji

	emojis := EmojiRegex.FindAllString(m.Content, -1)

	if len(emojis) < 1 {
		return toReturn
	}

	for _, em := range emojis {
		parts := strings.Split(em, ":")
		toReturn = append(toReturn, &Emoji{
			ID:       parts[2][:len(parts[2])-1],
			Name:     parts[1],
			Animated: strings.HasPrefix(em, "<a:"),
		})
	}

	return toReturn
}

type MessageFlags int

const (
	MessageFlagsCrossPosted          MessageFlags = 1 << 0
	MessageFlagsIsCrossPosted        MessageFlags = 1 << 1
	MessageFlagsSupressEmbeds        MessageFlags = 1 << 2
	MessageFlagsSourceMessageDeleted MessageFlags = 1 << 3
	MessageFlagsUrgent               MessageFlags = 1 << 4
)

type File struct {
	Name        string
	ContentType string
	Reader      io.Reader
}

type MessageSend struct {
	Content         string                  `json:"content,omitempty"`
	Embed           *MessageEmbed           `json:"embed,omitempty"`
	TTS             bool                    `json:"tts"`
	Components      []MessageComponent      `json:"components"`
	Files           []*File                 `json:"-"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Reference       *MessageReference       `json:"message_reference,omitempty"`
	File            *File                   `json:"-"`
}

type MessageEdit struct {
	Content         *string                 `json:"content,omitempty"`
	Components      []MessageComponent      `json:"components"`
	Embed           *MessageEmbed           `json:"embed,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	ID              string
	Channel         string
}

func NewMessageEdit(channelID string, messageID string) *MessageEdit {
	return &MessageEdit{
		Channel: channelID,
		ID:      messageID,
	}
}

func (m *MessageEdit) SetContent(str string) *MessageEdit {
	m.Content = &str
	return m
}

func (m *MessageEdit) SetEmbed(embed *MessageEmbed) *MessageEdit {
	m.Embed = embed
	return m
}

type AllowedMentionType string

const (
	AllowedMentionTypeRoles    AllowedMentionType = "roles"
	AllowedMentionTypeUsers    AllowedMentionType = "users"
	AllowedMentionTypeEveryone AllowedMentionType = "everyone"
)

type MessageAllowedMentions struct {
	Parse []AllowedMentionType `json:"parse"`
	Roles []string             `json:"roles,omitempty"`
	Users []string             `json:"users,omitempty"`
}

type MessageAttachment struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
}

type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        EmbedType              `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int64                  `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

type EmbedType string

const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGifv    EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

type MessageReactions struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

type MessageActivity struct {
	Type    MessageActivityType `json:"type"`
	PartyID string              `json:"party_id"`
}

type MessageActivityType int

const (
	MessageActivityTypeJoin        MessageActivityType = 1
	MessageActivityTypeSpectate    MessageActivityType = 2
	MessageActivityTypeListen      MessageActivityType = 3
	MessageActivityTypeJoinRequest MessageActivityType = 5
)

type MessageApplication struct {
	ID          string `json:"id"`
	CoverImage  string `json:"cover_image"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
}

type MessageReference struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

func (m *Message) Reference() *MessageReference {
	return &MessageReference{
		GuildID:   m.GuildID,
		ChannelID: m.ChannelID,
		MessageID: m.ID,
	}
}

func (m *Message) ContentWithMentionsReplaced() (content string) {
	content = m.Content

	for _, user := range m.Mentions {
		content = strings.NewReplacer(
			"<@"+user.ID+">", "@"+user.Username,
			"<@!"+user.ID+">", "@"+user.Username,
		).Replace(content)
	}
	return
}

var patternChannels = regexp.MustCompile("<#[^>]*>")

func (m *Message) ContentWithMoreMentionsReplaced(s *Session) (content string, err error) {
	content = m.Content

	if !s.StateEnabled {
		content = m.ContentWithMentionsReplaced()
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		content = m.ContentWithMentionsReplaced()
		return
	}

	for _, user := range m.Mentions {
		nick := user.Username

		member, err := s.State.Member(channel.GuildID, user.ID)
		if err == nil && member.Nick != "" {
			nick = member.Nick
		}

		content = strings.NewReplacer(
			"<@"+user.ID+">", "@"+user.Username,
			"<@!"+user.ID+">", "@"+nick,
		).Replace(content)
	}

	for _, roleID := range m.MentionRoles {
		role, err := s.State.Role(channel.GuildID, roleID)
		if err != nil || !role.Mentionable {
			continue
		}

		content = strings.Replace(content, "<@&"+role.ID+">", "@"+role.Name, -1)
	}

	content = patternChannels.ReplaceAllStringFunc(content, func(mention string) string {
		channel, err := s.State.Channel(mention[2 : len(mention)-1])
		if err != nil || channel.Type == ChannelTypeGuildVoice {
			return mention
		}

		return "#" + channel.Name
	})

	return
}
