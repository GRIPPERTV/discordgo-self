package discordgoself

import "strings"

type UserFlags int

const (
	UserFlagDiscordEmployee           UserFlags = 1 << 0
	UserFlagDiscordPartner            UserFlags = 1 << 1
	UserFlagHypeSquadEvents           UserFlags = 1 << 2
	UserFlagBugHunterLevel1           UserFlags = 1 << 3
	UserFlagHouseBravery              UserFlags = 1 << 6
	UserFlagHouseBrilliance           UserFlags = 1 << 7
	UserFlagHouseBalance              UserFlags = 1 << 8
	UserFlagEarlySupporter            UserFlags = 1 << 9
	UserFlagTeamUser                  UserFlags = 1 << 10
	UserFlagSystem                    UserFlags = 1 << 12
	UserFlagBugHunterLevel2           UserFlags = 1 << 14
	UserFlagVerifiedBot               UserFlags = 1 << 16
	UserFlagVerifiedBotDeveloper      UserFlags = 1 << 17
	UserFlagDiscordCertifiedModerator UserFlags = 1 << 18
)

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Avatar        string    `json:"avatar"`
	Bot           bool      `json:"bot"`
	System        bool      `json:"system"`
	Flags         int       `json:"flags"`
	PublicFlags   UserFlags `json:"public_flags"`
	Banner        string    `json:"banner"`
	BannerColor   string    `json:"banner_color"`
	AccentColor   int       `json:"accent_color"`
	Email         string    `json:"email"`
	Token         string    `json:"token"`
}

type Profile struct {
	User struct {
		ID                string            `json:"id"`
		Username          string            `json:"username"`
		Avatar            string            `json:"avatar"`
		Discriminator     string            `json:"discriminator"`
		PublicFlags       UserFlags         `json:"public_flags"`
		Flags             int               `json:"flags"`
		Bio               string            `json:"bio"`
		Bot               bool              `json:"bot"`
		Banner            string            `json:"banner"`
		BannerColor       string            `json:"banner_color"`
		AccentColor       int               `json:"accent_color"`
		PremiumSince      Timestamp         `json:"premium_since"`
		PremiumGuildSince Timestamp         `json:"premium_guild_since"`
		MutualGuilds      []*MutualGuilds   `json:"mutual_guilds"`
		ConnectedAccounts []*UserConnection `json:"connected_accounts"`
	} `json:"user"`
}

func (u *User) Tag() string {
	return u.Username + "#" + u.Discriminator
}

func (u *User) Mention() string {
	return "<@" + u.ID + ">"
}

func (u *User) AvatarURL() string {
	var URL string
	if u.Avatar == "" {
		URL = EndpointDefaultUserAvatar(u.Discriminator)
	} else if strings.HasPrefix(u.Avatar, "a_") {
		URL = EndpointUserAvatarAnimated(u.ID, u.Avatar)
	} else {
		URL = EndpointUserAvatar(u.ID, u.Avatar)
	}

	return URL
}
