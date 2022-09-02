package discordgoself

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrJSONUnmarshal           = errors.New("json unmarshal")
	ErrVerificationLevelBounds = errors.New("VerificationLevel out of bounds, should be between 0 and 3")
	ErrPruneDaysBounds         = errors.New("the number of days should be more than or equal to 1")
	ErrGuildNoIcon             = errors.New("guild does not have an icon set")
	ErrGuildNoSplash           = errors.New("guild does not have a splash set")
	ErrUnauthorized            = errors.New("HTTP request was unauthorized.")
)

func (s *Session) Request(method, urlStr string, data interface{}) (response []byte, err error) {
	return s.RequestWithBucketID(method, urlStr, data, strings.SplitN(urlStr, "?", 2)[0])
}

func (s *Session) RequestWithBucketID(method, urlStr string, data interface{}, bucketID string) (response []byte, err error) {
	var body []byte
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return
		}
	}

	return s.request(method, urlStr, "application/json", body, bucketID, 0)
}

func (s *Session) request(method, urlStr, contentType string, b []byte, bucketID string, sequence int) (response []byte, err error) {
	if bucketID == "" {
		bucketID = strings.SplitN(urlStr, "?", 2)[0]
	}
	return s.RequestWithLockedBucket(method, urlStr, contentType, b, s.Ratelimiter.LockBucket(bucketID), sequence)
}

func (s *Session) RequestWithLockedBucket(method, urlStr, contentType string, b []byte, bucket *Bucket, sequence int) (response []byte, err error) {

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(b))
	if err != nil {
		bucket.Release(nil)
		return
	}

	if s.Identify.Token != "" {
		req.Header.Set("authorization", s.Identify.Token)
	}

	if b != nil {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.Client.Do(req)
	if err != nil {
		bucket.Release(nil)
		return
	}
	defer func() {
		err2 := resp.Body.Close()
		if err2 != nil {
			log.Println("error closing resp body")
		}
	}()

	err = bucket.Release(resp.Header)
	if err != nil {
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusNoContent:
	case http.StatusBadGateway:
		if sequence < s.MaxRestRetries {

			s.log(LogInformational, "%s Failed (%s), Retrying...", urlStr, resp.Status)
			response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.Ratelimiter.LockBucketObject(bucket), sequence+1)
		} else {
			err = fmt.Errorf("Exceeded Max retries HTTP %s, %s", resp.Status, response)
		}
	case 429:
		rl := TooManyRequests{}
		err = json.Unmarshal(response, &rl)
		if err != nil {
			if e, ok := err.(*json.SyntaxError); ok {
				s.log(LogError, "syntax error at byte offset %d", e.Offset)
			} else {
				s.log(LogError, "rate limit unmarshal error, %s", err)
			}
			s.log(LogInformational, "response: %s", response)
			return
		}
		s.log(LogInformational, "Rate Limiting %s, retry in %v", urlStr, rl.RetryAfter)
		s.handleEvent(rateLimitEventType, &RateLimit{TooManyRequests: &rl, URL: urlStr})

		time.Sleep(rl.RetryAfter)

		response, err = s.RequestWithLockedBucket(method, urlStr, contentType, b, s.Ratelimiter.LockBucketObject(bucket), sequence)
	case http.StatusUnauthorized:
		s.log(LogInformational, ErrUnauthorized.Error())
		err = ErrUnauthorized
		fallthrough
	default:
		err = newRestError(req, resp, response)
	}

	return
}

func unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrJSONUnmarshal, err)
	}

	return nil
}

func (s *Session) Login(email, password string) (err error) {

	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{email, password}

	response, err := s.RequestWithBucketID("POST", EndpointLogin, data, EndpointLogin)
	if err != nil {
		return
	}

	temp := struct {
		Token string `json:"token"`
		MFA   bool   `json:"mfa"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	s.Identify.Token = temp.Token
	s.MFA = temp.MFA
	return
}

func (s *Session) Register(username string) (token string, err error) {

	data := struct {
		Username string `json:"username"`
	}{username}

	response, err := s.RequestWithBucketID("POST", EndpointRegister, data, EndpointRegister)
	if err != nil {
		return
	}

	temp := struct {
		Token string `json:"token"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	token = temp.Token
	return
}

func (s *Session) Logout() (err error) {

	if s.Identify.Token == "" {
		return
	}

	data := struct {
		Token string `json:"token"`
	}{s.Identify.Token}

	_, err = s.RequestWithBucketID("POST", EndpointLogout, data, EndpointLogout)
	return
}

func (s *Session) User(userID string) (st *User, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointUser(userID), nil, EndpointUsers)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserProfile(userID string) (st *Profile, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointUserProfile(userID), nil, EndpointUserProfile(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserNote(userID string) (st *Note, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointUserNotes(userID), nil, EndpointUserNotes(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserNoteSet(userID string, message string) (err error) {
	data := struct {
		Note string `json:"note"`
	}{message}

	_, err = s.RequestWithBucketID("PUT", EndpointUserNotes(userID), data, EndpointUserNotes(""))
	return
}

func (s *Session) UserAvatarDecode(u *User) (img image.Image, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointUserAvatar(u.ID, u.Avatar), nil, EndpointUserAvatar("", ""))
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

func (s *Session) UserSetNewPassword(password, newPassword string) (st *User, err error) {
	data := struct {
		Password    string `json:"password,omitempty"`
		NewPassword string `json:"new_password,omitempty"`
	}{password, newPassword}

	body, err := s.RequestWithBucketID("PATCH", EndpointUser("@me"), data, EndpointUsers)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserSettings() (st *Settings, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointUserSettings("@me"), nil, EndpointUserSettings(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserUpdateStatus(status Status) (st *Settings, err error) {
	data := struct {
		Status Status `json:"status"`
	}{status}

	body, err := s.RequestWithBucketID("PATCH", EndpointUserSettings("@me"), data, EndpointUserSettings(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserChannels() (st []*Channel, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointUserChannels("@me"), nil, EndpointUserChannels(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserChannelCreate(recipientID string) (st *Channel, err error) {

	data := struct {
		RecipientID string `json:"recipient_id"`
	}{recipientID}

	body, err := s.RequestWithBucketID("POST", EndpointUserChannels("@me"), data, EndpointUserChannels(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserGuilds(limit int, beforeID, afterID string) (st []*UserGuild, err error) {

	v := url.Values{}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}

	uri := EndpointUserGuilds("@me")

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointUserGuilds(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) UserGuildSettingsEdit(guildID string, settings *UserGuildSettingsEdit) (st *UserGuildSettings, err error) {

	body, err := s.RequestWithBucketID("PATCH", EndpointUserGuildSettings("@me", guildID), settings, EndpointUserGuildSettings("", guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func memberPermissions(guild *Guild, channel *Channel, userID string, roles []string) (apermissions int64) {
	if userID == guild.OwnerID {
		apermissions = PermissionAll
		return
	}

	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			apermissions |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range roles {
			if role.ID == roleID {
				apermissions |= role.Permissions
				break
			}
		}
	}

	if apermissions&PermissionAdministrator == PermissionAdministrator {
		apermissions |= PermissionAll
	}

	for _, overwrite := range channel.PermissionOverwrites {
		if guild.ID == overwrite.ID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	var denies, allows int64
	for _, overwrite := range channel.PermissionOverwrites {
		for _, roleID := range roles {
			if overwrite.Type == PermissionOverwriteTypeRole && roleID == overwrite.ID {
				denies |= overwrite.Deny
				allows |= overwrite.Allow
				break
			}
		}
	}

	apermissions &= ^denies
	apermissions |= allows

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == PermissionOverwriteTypeMember && overwrite.ID == userID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	if apermissions&PermissionAdministrator == PermissionAdministrator {
		apermissions |= PermissionAllChannel
	}

	return apermissions
}


func (s *Session) Guild(guildID string) (st *Guild, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointGuild(guildID), nil, EndpointGuild(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildPreview(guildID string) (st *GuildPreview, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointGuildPreview(guildID), nil, EndpointGuildPreview(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildCreate(name string) (st *Guild, err error) {

	data := struct {
		Name string `json:"name"`
	}{name}

	body, err := s.RequestWithBucketID("POST", EndpointGuildCreate, data, EndpointGuildCreate)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildEdit(guildID string, g GuildParams) (st *Guild, err error) {

	if g.VerificationLevel != nil {
		val := *g.VerificationLevel
		if val < 0 || val > 4 {
			err = ErrVerificationLevelBounds
			return
		}
	}

	if g.Region != "" {
		isValid := false
		regions, _ := s.VoiceRegions()
		for _, r := range regions {
			if g.Region == r.ID {
				isValid = true
			}
		}
		if !isValid {
			var valid []string
			for _, r := range regions {
				valid = append(valid, r.ID)
			}
			err = fmt.Errorf("Region not a valid region (%q)", valid)
			return
		}
	}

	body, err := s.RequestWithBucketID("PATCH", EndpointGuild(guildID), g, EndpointGuild(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildDelete(guildID string) (st *Guild, err error) {

	body, err := s.RequestWithBucketID("DELETE", EndpointGuild(guildID), nil, EndpointGuild(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildLeave(guildID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointUserGuild("@me", guildID), nil, EndpointUserGuild("", guildID))
	return
}

func (s *Session) GuildBans(guildID string) (st []*GuildBan, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildBans(guildID), nil, EndpointGuildBans(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildBanCreate(guildID, userID string, days int) (err error) {
	return s.GuildBanCreateWithReason(guildID, userID, "", days)
}

func (s *Session) GuildBan(guildID, userID string) (st *GuildBan, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildBan(guildID, userID), nil, EndpointGuildBan(guildID, userID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildBanCreateWithReason(guildID, userID, reason string, days int) (err error) {

	uri := EndpointGuildBan(guildID, userID)

	queryParams := url.Values{}
	if days > 0 {
		queryParams.Set("delete_message_days", strconv.Itoa(days))
	}
	if reason != "" {
		queryParams.Set("reason", reason)
	}

	if len(queryParams) > 0 {
		uri += "?" + queryParams.Encode()
	}

	_, err = s.RequestWithBucketID("PUT", uri, nil, EndpointGuildBan(guildID, ""))
	return
}

func (s *Session) GuildBanDelete(guildID, userID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointGuildBan(guildID, userID), nil, EndpointGuildBan(guildID, ""))
	return
}

func (s *Session) GuildMembers(guildID string, after string, limit int) (st []*Member, err error) {

	uri := EndpointGuildMembers(guildID)

	v := url.Values{}

	if after != "" {
		v.Set("after", after)
	}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointGuildMembers(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildMember(guildID, userID string) (st *Member, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildMember(guildID, userID), nil, EndpointGuildMember(guildID, ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildMemberAdd(accessToken, guildID, userID, nick string, roles []string, mute, deaf bool) (err error) {

	data := struct {
		AccessToken string   `json:"access_token"`
		Nick        string   `json:"nick,omitempty"`
		Roles       []string `json:"roles,omitempty"`
		Mute        bool     `json:"mute,omitempty"`
		Deaf        bool     `json:"deaf,omitempty"`
	}{accessToken, nick, roles, mute, deaf}

	_, err = s.RequestWithBucketID("PUT", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	if err != nil {
		return err
	}

	return err
}

func (s *Session) GuildMemberDelete(guildID, userID string) (err error) {

	return s.GuildMemberDeleteWithReason(guildID, userID, "")
}

func (s *Session) GuildMemberDeleteWithReason(guildID, userID, reason string) (err error) {

	uri := EndpointGuildMember(guildID, userID)
	if reason != "" {
		uri += "?reason=" + url.QueryEscape(reason)
	}

	_, err = s.RequestWithBucketID("DELETE", uri, nil, EndpointGuildMember(guildID, ""))
	return
}

func (s *Session) GuildMemberEdit(guildID, userID string, roles []string) (err error) {

	data := struct {
		Roles []string `json:"roles"`
	}{roles}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	return
}

func (s *Session) GuildMemberMove(guildID string, userID string, channelID *string) (err error) {
	data := struct {
		ChannelID *string `json:"channel_id"`
	}{channelID}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	return
}

func (s *Session) GuildMemberNickname(guildID, userID, nickname string) (err error) {

	data := struct {
		Nick string `json:"nick"`
	}{nickname}

	if userID == "@me" {
		userID += "/nick"
	}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	return
}

func (s *Session) GuildMemberMute(guildID string, userID string, mute bool) (err error) {
	data := struct {
		Mute bool `json:"mute"`
	}{mute}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	return
}

func (s *Session) GuildMemberDeafen(guildID string, userID string, deaf bool) (err error) {
	data := struct {
		Deaf bool `json:"deaf"`
	}{deaf}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildMember(guildID, userID), data, EndpointGuildMember(guildID, ""))
	return
}

func (s *Session) GuildMemberRoleAdd(guildID, userID, roleID string) (err error) {

	_, err = s.RequestWithBucketID("PUT", EndpointGuildMemberRole(guildID, userID, roleID), nil, EndpointGuildMemberRole(guildID, "", ""))

	return
}

func (s *Session) GuildMemberRoleRemove(guildID, userID, roleID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointGuildMemberRole(guildID, userID, roleID), nil, EndpointGuildMemberRole(guildID, "", ""))

	return
}

func (s *Session) GuildChannels(guildID string) (st []*Channel, err error) {

	body, err := s.request("GET", EndpointGuildChannels(guildID), "", nil, EndpointGuildChannels(guildID), 0)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

type GuildChannelCreateData struct {
	Name                 string                 `json:"name"`
	Type                 ChannelType            `json:"type"`
	Topic                string                 `json:"topic,omitempty"`
	Bitrate              int                    `json:"bitrate,omitempty"`
	UserLimit            int                    `json:"user_limit,omitempty"`
	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
	Position             int                    `json:"position,omitempty"`
	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
	NSFW                 bool                   `json:"nsfw,omitempty"`
}

func (s *Session) GuildChannelCreateComplex(guildID string, data GuildChannelCreateData) (st *Channel, err error) {
	body, err := s.RequestWithBucketID("POST", EndpointGuildChannels(guildID), data, EndpointGuildChannels(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildChannelCreate(guildID, name string, ctype ChannelType) (st *Channel, err error) {
	return s.GuildChannelCreateComplex(guildID, GuildChannelCreateData{
		Name: name,
		Type: ctype,
	})
}

func (s *Session) GuildChannelsReorder(guildID string, channels []*Channel) (err error) {

	data := make([]struct {
		ID       string `json:"id"`
		Position int    `json:"position"`
	}, len(channels))

	for i, c := range channels {
		data[i].ID = c.ID
		data[i].Position = c.Position
	}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildChannels(guildID), data, EndpointGuildChannels(guildID))
	return
}

func (s *Session) GuildInvites(guildID string) (st []*Invite, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointGuildInvites(guildID), nil, EndpointGuildInvites(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildRoles(guildID string) (st []*Role, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildRoles(guildID), nil, EndpointGuildRoles(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildRoleCreate(guildID string) (st *Role, err error) {

	body, err := s.RequestWithBucketID("POST", EndpointGuildRoles(guildID), nil, EndpointGuildRoles(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildRoleEdit(guildID, roleID, name string, color int, hoist bool, perm int64, mention bool) (st *Role, err error) {

	if color > 0xFFFFFF {
		err = fmt.Errorf("color value cannot be larger than 0xFFFFFF")
		return nil, err
	}

	data := struct {
		Name        string `json:"name"`
		Color       int    `json:"color"`
		Hoist       bool   `json:"hoist"`
		Permissions int64  `json:"permissions,string"`
		Mentionable bool   `json:"mentionable"`
	}{name, color, hoist, perm, mention}

	body, err := s.RequestWithBucketID("PATCH", EndpointGuildRole(guildID, roleID), data, EndpointGuildRole(guildID, ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildRoleReorder(guildID string, roles []*Role) (st []*Role, err error) {

	body, err := s.RequestWithBucketID("PATCH", EndpointGuildRoles(guildID), roles, EndpointGuildRoles(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildRoleDelete(guildID, roleID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointGuildRole(guildID, roleID), nil, EndpointGuildRole(guildID, ""))

	return
}

func (s *Session) GuildPruneCount(guildID string, days uint32) (count uint32, err error) {
	count = 0

	if days <= 0 {
		err = ErrPruneDaysBounds
		return
	}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	uri := EndpointGuildPrune(guildID) + "?days=" + strconv.FormatUint(uint64(days), 10)
	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointGuildPrune(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &p)
	if err != nil {
		return
	}

	count = p.Pruned

	return
}

func (s *Session) GuildPrune(guildID string, days uint32) (count uint32, err error) {

	count = 0

	if days <= 0 {
		err = ErrPruneDaysBounds
		return
	}

	data := struct {
		days uint32
	}{days}

	p := struct {
		Pruned uint32 `json:"pruned"`
	}{}

	body, err := s.RequestWithBucketID("POST", EndpointGuildPrune(guildID), data, EndpointGuildPrune(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &p)
	if err != nil {
		return
	}

	count = p.Pruned

	return
}

func (s *Session) GuildIntegrations(guildID string) (st []*Integration, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildIntegrations(guildID), nil, EndpointGuildIntegrations(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildIntegrationCreate(guildID, integrationType, integrationID string) (err error) {

	data := struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{integrationType, integrationID}

	_, err = s.RequestWithBucketID("POST", EndpointGuildIntegrations(guildID), data, EndpointGuildIntegrations(guildID))
	return
}

func (s *Session) GuildIntegrationEdit(guildID, integrationID string, expireBehavior, expireGracePeriod int, enableEmoticons bool) (err error) {

	data := struct {
		ExpireBehavior    int  `json:"expire_behavior"`
		ExpireGracePeriod int  `json:"expire_grace_period"`
		EnableEmoticons   bool `json:"enable_emoticons"`
	}{expireBehavior, expireGracePeriod, enableEmoticons}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildIntegration(guildID, integrationID), data, EndpointGuildIntegration(guildID, ""))
	return
}

func (s *Session) GuildIntegrationDelete(guildID, integrationID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointGuildIntegration(guildID, integrationID), nil, EndpointGuildIntegration(guildID, ""))
	return
}

func (s *Session) GuildIntegrationSync(guildID, integrationID string) (err error) {

	_, err = s.RequestWithBucketID("POST", EndpointGuildIntegrationSync(guildID, integrationID), nil, EndpointGuildIntegration(guildID, ""))
	return
}

func (s *Session) GuildIcon(guildID string) (img image.Image, err error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return
	}

	if g.Icon == "" {
		err = ErrGuildNoIcon
		return
	}

	body, err := s.RequestWithBucketID("GET", EndpointGuildIcon(guildID, g.Icon), nil, EndpointGuildIcon(guildID, ""))
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

func (s *Session) GuildSplash(guildID string) (img image.Image, err error) {
	g, err := s.Guild(guildID)
	if err != nil {
		return
	}

	if g.Splash == "" {
		err = ErrGuildNoSplash
		return
	}

	body, err := s.RequestWithBucketID("GET", EndpointGuildSplash(guildID, g.Splash), nil, EndpointGuildSplash(guildID, ""))
	if err != nil {
		return
	}

	img, _, err = image.Decode(bytes.NewReader(body))
	return
}

func (s *Session) GuildEmbed(guildID string) (st *GuildEmbed, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildEmbed(guildID), nil, EndpointGuildEmbed(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildEmbedEdit(guildID string, enabled bool, channelID string) (err error) {

	data := GuildEmbed{enabled, channelID}

	_, err = s.RequestWithBucketID("PATCH", EndpointGuildEmbed(guildID), data, EndpointGuildEmbed(guildID))
	return
}

func (s *Session) GuildAuditLog(guildID, userID, beforeID string, actionType, limit int) (st *GuildAuditLog, err error) {

	uri := EndpointGuildAuditLogs(guildID)

	v := url.Values{}
	if userID != "" {
		v.Set("user_id", userID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if actionType > 0 {
		v.Set("action_type", strconv.Itoa(actionType))
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if len(v) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, v.Encode())
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointGuildAuditLogs(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) GuildEmojis(guildID string) (emoji []*Emoji, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildEmojis(guildID), nil, EndpointGuildEmojis(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &emoji)
	return
}

func (s *Session) GuildEmojiCreate(guildID, name, image string, roles []string) (emoji *Emoji, err error) {

	data := struct {
		Name  string   `json:"name"`
		Image string   `json:"image"`
		Roles []string `json:"roles,omitempty"`
	}{name, image, roles}

	body, err := s.RequestWithBucketID("POST", EndpointGuildEmojis(guildID), data, EndpointGuildEmojis(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &emoji)
	return
}

func (s *Session) GuildEmojiEdit(guildID, emojiID, name string, roles []string) (emoji *Emoji, err error) {

	data := struct {
		Name  string   `json:"name"`
		Roles []string `json:"roles,omitempty"`
	}{name, roles}

	body, err := s.RequestWithBucketID("PATCH", EndpointGuildEmoji(guildID, emojiID), data, EndpointGuildEmojis(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &emoji)
	return
}

func (s *Session) GuildEmojiDelete(guildID, emojiID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointGuildEmoji(guildID, emojiID), nil, EndpointGuildEmojis(guildID))
	return
}


func (s *Session) Channel(channelID string) (st *Channel, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointChannel(channelID), nil, EndpointChannel(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelEdit(channelID, name string) (*Channel, error) {
	return s.ChannelEditComplex(channelID, &ChannelEdit{
		Name: name,
	})
}

func (s *Session) ChannelEditComplex(channelID string, data *ChannelEdit) (st *Channel, err error) {
	body, err := s.RequestWithBucketID("PATCH", EndpointChannel(channelID), data, EndpointChannel(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelDelete(channelID string) (st *Channel, err error) {

	body, err := s.RequestWithBucketID("DELETE", EndpointChannel(channelID), nil, EndpointChannel(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelTyping(channelID string) (err error) {

	_, err = s.RequestWithBucketID("POST", EndpointChannelTyping(channelID), nil, EndpointChannelTyping(channelID))
	return
}

func (s *Session) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (st []*Message, err error) {

	uri := EndpointChannelMessages(channelID)

	v := url.Values{}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}
	if aroundID != "" {
		v.Set("around", aroundID)
	}
	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointChannelMessages(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelMessage(channelID, messageID string) (st *Message, err error) {

	response, err := s.RequestWithBucketID("GET", EndpointChannelMessage(channelID, messageID), nil, EndpointChannelMessage(channelID, ""))
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

func (s *Session) ChannelMessageAck(channelID, messageID, lastToken string) (st *Ack, err error) {

	body, err := s.RequestWithBucketID("POST", EndpointChannelMessageAck(channelID, messageID), &Ack{Token: lastToken}, EndpointChannelMessageAck(channelID, ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelMessageSend(channelID string, content string) (*Message, error) {
	return s.ChannelMessageSendComplex(channelID, &MessageSend{
		Content: content,
	})
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func (s *Session) ChannelMessageSendComplex(channelID string, data *MessageSend) (st *Message, err error) {
	if data.Embed != nil && data.Embed.Type == "" {
		data.Embed.Type = "rich"
	}

	endpoint := EndpointChannelMessages(channelID)

	files := data.Files
	if data.File != nil {
		if files == nil {
			files = []*File{data.File}
		} else {
			err = fmt.Errorf("cannot specify both File and Files")
			return
		}
	}

	var response []byte
	if len(files) > 0 {
		contentType, body, encodeErr := MultipartBodyWithJSON(data, files)
		if encodeErr != nil {
			return st, encodeErr
		}

		response, err = s.request("POST", endpoint, contentType, body, endpoint, 0)
	} else {
		response, err = s.RequestWithBucketID("POST", endpoint, data, endpoint)
	}
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

func (s *Session) ChannelMessageSendTTS(channelID string, content string) (*Message, error) {
	return s.ChannelMessageSendComplex(channelID, &MessageSend{
		Content: content,
		TTS:     true,
	})
}

func (s *Session) ChannelMessageSendEmbed(channelID string, embed *MessageEmbed) (*Message, error) {
	return s.ChannelMessageSendComplex(channelID, &MessageSend{
		Embed: embed,
	})
}

func (s *Session) ChannelMessageSendReply(channelID string, content string, reference *MessageReference) (*Message, error) {
	return s.ChannelMessageSendComplex(channelID, &MessageSend{
		Content:   content,
		Reference: reference,
	})
}

func (s *Session) ChannelMessageEdit(channelID, messageID, content string) (*Message, error) {
	return s.ChannelMessageEditComplex(NewMessageEdit(channelID, messageID).SetContent(content))
}

func (s *Session) ChannelMessageEditComplex(m *MessageEdit) (st *Message, err error) {
	if m.Embed != nil && m.Embed.Type == "" {
		m.Embed.Type = "rich"
	}

	response, err := s.RequestWithBucketID("PATCH", EndpointChannelMessage(m.Channel, m.ID), m, EndpointChannelMessage(m.Channel, ""))
	if err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

func (s *Session) ChannelMessageEditEmbed(channelID, messageID string, embed *MessageEmbed) (*Message, error) {
	return s.ChannelMessageEditComplex(NewMessageEdit(channelID, messageID).SetEmbed(embed))
}

func (s *Session) ChannelMessageDelete(channelID, messageID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointChannelMessage(channelID, messageID), nil, EndpointChannelMessage(channelID, ""))
	return
}

func (s *Session) ChannelMessagesBulkDelete(channelID string, messages []string) (err error) {

	if len(messages) == 0 {
		return
	}

	if len(messages) == 1 {
		err = s.ChannelMessageDelete(channelID, messages[0])
		return
	}

	if len(messages) > 100 {
		messages = messages[:100]
	}

	data := struct {
		Messages []string `json:"messages"`
	}{messages}

	_, err = s.RequestWithBucketID("POST", EndpointChannelMessagesBulkDelete(channelID), data, EndpointChannelMessagesBulkDelete(channelID))
	return
}

func (s *Session) ChannelMessagePin(channelID, messageID string) (err error) {

	_, err = s.RequestWithBucketID("PUT", EndpointChannelMessagePin(channelID, messageID), nil, EndpointChannelMessagePin(channelID, ""))
	return
}

func (s *Session) ChannelMessageUnpin(channelID, messageID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointChannelMessagePin(channelID, messageID), nil, EndpointChannelMessagePin(channelID, ""))
	return
}

func (s *Session) ChannelMessagesPinned(channelID string) (st []*Message, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointChannelMessagesPins(channelID), nil, EndpointChannelMessagesPins(channelID))

	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelFileSend(channelID, name string, r io.Reader) (*Message, error) {
	return s.ChannelMessageSendComplex(channelID, &MessageSend{File: &File{Name: name, Reader: r}})
}

func (s *Session) ChannelInvites(channelID string) (st []*Invite, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointChannelInvites(channelID), nil, EndpointChannelInvites(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelInviteCreate(channelID string, i Invite) (st *Invite, err error) {

	data := struct {
		MaxAge    int  `json:"max_age"`
		MaxUses   int  `json:"max_uses"`
		Temporary bool `json:"temporary"`
		Unique    bool `json:"unique"`
	}{i.MaxAge, i.MaxUses, i.Temporary, i.Unique}

	body, err := s.RequestWithBucketID("POST", EndpointChannelInvites(channelID), data, EndpointChannelInvites(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelPermissionSet(channelID, targetID string, targetType PermissionOverwriteType, allow, deny int64) (err error) {

	data := struct {
		ID    string                  `json:"id"`
		Type  PermissionOverwriteType `json:"type"`
		Allow int64                   `json:"allow,string"`
		Deny  int64                   `json:"deny,string"`
	}{targetID, targetType, allow, deny}

	_, err = s.RequestWithBucketID("PUT", EndpointChannelPermission(channelID, targetID), data, EndpointChannelPermission(channelID, ""))
	return
}

func (s *Session) ChannelPermissionDelete(channelID, targetID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointChannelPermission(channelID, targetID), nil, EndpointChannelPermission(channelID, ""))
	return
}

func (s *Session) ChannelMessageCrosspost(channelID, messageID string) (st *Message, err error) {

	endpoint := EndpointChannelMessageCrosspost(channelID, messageID)

	body, err := s.RequestWithBucketID("POST", endpoint, nil, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) ChannelNewsFollow(channelID, targetID string) (st *ChannelFollow, err error) {

	endpoint := EndpointChannelFollow(channelID)

	data := struct {
		WebhookChannelID string `json:"webhook_channel_id"`
	}{targetID}

	body, err := s.RequestWithBucketID("POST", endpoint, data, endpoint)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}


func (s *Session) Invite(inviteID string) (st *Invite, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointInvite(inviteID), nil, EndpointInvite(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) InviteWithCounts(inviteID string) (st *Invite, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointInvite(inviteID)+"?with_counts=true", nil, EndpointInvite(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) InviteDelete(inviteID string) (st *Invite, err error) {

	body, err := s.RequestWithBucketID("DELETE", EndpointInvite(inviteID), nil, EndpointInvite(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) InviteAccept(inviteID string) (st *Invite, err error) {

	body, err := s.RequestWithBucketID("POST", EndpointInvite(inviteID), nil, EndpointInvite(""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}


func (s *Session) VoiceRegions() (st []*VoiceRegion, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointVoiceRegions, nil, EndpointVoiceRegions)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) VoiceICE() (st *VoiceICE, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointVoiceIce, nil, EndpointVoiceIce)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}


func (s *Session) Gateway() (gateway string, err error) {

	response, err := s.RequestWithBucketID("GET", EndpointGateway, nil, EndpointGateway)
	if err != nil {
		return
	}

	temp := struct {
		URL string `json:"url"`
	}{}

	err = unmarshal(response, &temp)
	if err != nil {
		return
	}

	gateway = temp.URL

	if !strings.HasSuffix(gateway, "/") {
		gateway += "/"
	}

	return
}


func (s *Session) WebhookCreate(channelID, name, avatar string) (st *Webhook, err error) {

	data := struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	body, err := s.RequestWithBucketID("POST", EndpointChannelWebhooks(channelID), data, EndpointChannelWebhooks(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) ChannelWebhooks(channelID string) (st []*Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointChannelWebhooks(channelID), nil, EndpointChannelWebhooks(channelID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) GuildWebhooks(guildID string) (st []*Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointGuildWebhooks(guildID), nil, EndpointGuildWebhooks(guildID))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) Webhook(webhookID string) (st *Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointWebhook(webhookID), nil, EndpointWebhooks)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) WebhookWithToken(webhookID, token string) (st *Webhook, err error) {

	body, err := s.RequestWithBucketID("GET", EndpointWebhookToken(webhookID, token), nil, EndpointWebhookToken("", ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) WebhookEdit(webhookID, name, avatar, channelID string) (st *Role, err error) {

	data := struct {
		Name      string `json:"name,omitempty"`
		Avatar    string `json:"avatar,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	}{name, avatar, channelID}

	body, err := s.RequestWithBucketID("PATCH", EndpointWebhook(webhookID), data, EndpointWebhooks)
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) WebhookEditWithToken(webhookID, token, name, avatar string) (st *Role, err error) {

	data := struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	}{name, avatar}

	body, err := s.RequestWithBucketID("PATCH", EndpointWebhookToken(webhookID, token), data, EndpointWebhookToken("", ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) WebhookDelete(webhookID string) (err error) {

	_, err = s.RequestWithBucketID("DELETE", EndpointWebhook(webhookID), nil, EndpointWebhooks)

	return
}

func (s *Session) WebhookDeleteWithToken(webhookID, token string) (st *Webhook, err error) {

	body, err := s.RequestWithBucketID("DELETE", EndpointWebhookToken(webhookID, token), nil, EndpointWebhookToken("", ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)

	return
}

func (s *Session) WebhookExecute(webhookID, token string, wait bool, data *WebhookParams) (st *Message, err error) {
	uri := EndpointWebhookToken(webhookID, token)

	if wait {
		uri += "?wait=true"
	}

	var response []byte
	if len(data.Files) > 0 {
		contentType, body, encodeErr := MultipartBodyWithJSON(data, data.Files)
		if encodeErr != nil {
			return st, encodeErr
		}

		response, err = s.request("POST", uri, contentType, body, uri, 0)
	} else {
		response, err = s.RequestWithBucketID("POST", uri, data, uri)
	}
	if !wait || err != nil {
		return
	}

	err = unmarshal(response, &st)
	return
}

func (s *Session) WebhookMessage(webhookID, token, messageID string) (message *Message, err error) {
	uri := EndpointWebhookMessage(webhookID, token, messageID)

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointWebhookToken("", ""))
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &message)

	return
}

func (s *Session) WebhookMessageEdit(webhookID, token, messageID string, data *WebhookEdit) (err error) {
	uri := EndpointWebhookMessage(webhookID, token, messageID)
	if len(data.Files) > 0 {
		contentType, body, err := MultipartBodyWithJSON(data, data.Files)
		if err != nil {
			return err
		}

		_, err = s.request("PATCH", uri, contentType, body, uri, 0)
	} else {
		_, err = s.RequestWithBucketID("PATCH", uri, data, EndpointWebhookToken("", ""))
	}
	return
}

func (s *Session) WebhookMessageDelete(webhookID, token, messageID string) (err error) {
	uri := EndpointWebhookMessage(webhookID, token, messageID)

	_, err = s.RequestWithBucketID("DELETE", uri, nil, EndpointWebhookToken("", ""))
	return
}

func (s *Session) MessageReactionAdd(channelID, messageID, emojiID string) error {

	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID("PUT", EndpointMessageReaction(channelID, messageID, emojiID, "@me"), nil, EndpointMessageReaction(channelID, "", "", ""))

	return err
}

func (s *Session) MessageReactionRemove(channelID, messageID, emojiID, userID string) error {

	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID("DELETE", EndpointMessageReaction(channelID, messageID, emojiID, userID), nil, EndpointMessageReaction(channelID, "", "", ""))

	return err
}

func (s *Session) MessageReactionsRemoveAll(channelID, messageID string) error {

	_, err := s.RequestWithBucketID("DELETE", EndpointMessageReactionsAll(channelID, messageID), nil, EndpointMessageReactionsAll(channelID, messageID))

	return err
}

func (s *Session) MessageReactionsRemoveEmoji(channelID, messageID, emojiID string) error {

	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	_, err := s.RequestWithBucketID("DELETE", EndpointMessageReactions(channelID, messageID, emojiID), nil, EndpointMessageReactions(channelID, messageID, emojiID))

	return err
}

func (s *Session) MessageReactions(channelID, messageID, emojiID string, limit int, beforeID, afterID string) (st []*User, err error) {
	emojiID = strings.Replace(emojiID, "#", "%23", -1)
	uri := EndpointMessageReactions(channelID, messageID, emojiID)

	v := url.Values{}

	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}

	if afterID != "" {
		v.Set("after", afterID)
	}
	if beforeID != "" {
		v.Set("before", beforeID)
	}

	if len(v) > 0 {
		uri += "?" + v.Encode()
	}

	body, err := s.RequestWithBucketID("GET", uri, nil, EndpointMessageReaction(channelID, "", "", ""))
	if err != nil {
		return
	}

	err = unmarshal(body, &st)
	return
}

func (s *Session) RelationshipsGet() (r []*Relationship, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointRelationships(), nil, EndpointRelationships())
	if err != nil {
		return
	}

	err = unmarshal(body, &r)
	return
}

func (s *Session) relationshipCreate(userID string, relationshipType int) (err error) {
	data := struct {
		Type int `json:"type"`
	}{relationshipType}

	_, err = s.RequestWithBucketID("PUT", EndpointRelationship(userID), data, EndpointRelationships())
	return
}

func (s *Session) RelationshipFriendRequestSend(userID string) (err error) {
	err = s.relationshipCreate(userID, 4)
	return
}

func (s *Session) RelationshipFriendRequestAccept(userID string) (err error) {
	err = s.relationshipCreate(userID, 1)
	return
}

func (s *Session) RelationshipUserBlock(userID string) (err error) {
	err = s.relationshipCreate(userID, 2)
	return
}

func (s *Session) RelationshipDelete(userID string) (err error) {
	_, err = s.RequestWithBucketID("DELETE", EndpointRelationship(userID), nil, EndpointRelationships())
	return
}

func (s *Session) RelationshipsMutualGet(userID string) (mf []*User, err error) {
	body, err := s.RequestWithBucketID("GET", EndpointRelationshipsMutual(userID), nil, EndpointRelationshipsMutual(userID))
	if err != nil {
		return
	}

	err = unmarshal(body, &mf)
	return
}
