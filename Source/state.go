package discordgoself

import (
	"errors"
	"sort"
	"sync"
)

var ErrNilState = errors.New("state not instantiated, please use discordgoself.New() or assign Session.State")
var ErrStateNotFound = errors.New("state cache not found")
var ErrMessageIncompletePermissions = errors.New("message incomplete, unable to determine permissions")

type State struct {
	sync.RWMutex
	Ready
	MaxMessageCount int
	TrackChannels   bool
	TrackEmojis     bool
	TrackMembers    bool
	TrackRoles      bool
	TrackVoice      bool
	TrackPresences  bool
	guildMap        map[string]*Guild
	channelMap      map[string]*Channel
	memberMap       map[string]map[string]*Member
}

func NewState() *State {
	return &State{
		Ready: Ready{
			PrivateChannels: []*Channel{},
			Guilds:          []*Guild{},
		},
		TrackChannels:  true,
		TrackEmojis:    true,
		TrackMembers:   true,
		TrackRoles:     true,
		TrackVoice:     true,
		TrackPresences: true,
		guildMap:       make(map[string]*Guild),
		channelMap:     make(map[string]*Channel),
		memberMap:      make(map[string]map[string]*Member),
	}
}

func (s *State) createMemberMap(guild *Guild) {
	members := make(map[string]*Member)
	for _, m := range guild.Members {
		members[m.User.ID] = m
	}
	s.memberMap[guild.ID] = members
}

func (s *State) GuildAdd(guild *Guild) error {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	for _, c := range guild.Channels {
		s.channelMap[c.ID] = c
	}

	if guild.Members != nil {
		s.createMemberMap(guild)
	} else if _, ok := s.memberMap[guild.ID]; !ok {
		s.memberMap[guild.ID] = make(map[string]*Member)
	}

	if g, ok := s.guildMap[guild.ID]; ok {
		if guild.MemberCount == 0 {
			guild.MemberCount = g.MemberCount
		}
		if guild.Roles == nil {
			guild.Roles = g.Roles
		}
		if guild.Emojis == nil {
			guild.Emojis = g.Emojis
		}
		if guild.Members == nil {
			guild.Members = g.Members
		}
		if guild.Presences == nil {
			guild.Presences = g.Presences
		}
		if guild.Channels == nil {
			guild.Channels = g.Channels
		}
		if guild.VoiceStates == nil {
			guild.VoiceStates = g.VoiceStates
		}
		*g = *guild
		return nil
	}

	s.Guilds = append(s.Guilds, guild)
	s.guildMap[guild.ID] = guild

	return nil
}

func (s *State) GuildRemove(guild *Guild) error {
	if s == nil {
		return ErrNilState
	}

	_, err := s.Guild(guild.ID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	delete(s.guildMap, guild.ID)

	for i, g := range s.Guilds {
		if g.ID == guild.ID {
			s.Guilds = append(s.Guilds[:i], s.Guilds[i+1:]...)
			return nil
		}
	}

	return nil
}

func (s *State) Guild(guildID string) (*Guild, error) {
	if s == nil {
		return nil, ErrNilState
	}

	s.RLock()
	defer s.RUnlock()

	if g, ok := s.guildMap[guildID]; ok {
		return g, nil
	}

	return nil, ErrStateNotFound
}

func (s *State) PresenceAdd(guildID string, presence *Presence) error {
	if s == nil {
		return ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for i, p := range guild.Presences {
		if p.User.ID == presence.User.ID {
			guild.Presences[i].Activities = presence.Activities
			if presence.Status != "" {
				guild.Presences[i].Status = presence.Status
			}

			guild.Presences[i].User.ID = presence.User.ID

			if presence.User.Avatar != "" {
				guild.Presences[i].User.Avatar = presence.User.Avatar
			}
			if presence.User.Discriminator != "" {
				guild.Presences[i].User.Discriminator = presence.User.Discriminator
			}
			if presence.User.Email != "" {
				guild.Presences[i].User.Email = presence.User.Email
			}
			if presence.User.Token != "" {
				guild.Presences[i].User.Token = presence.User.Token
			}
			if presence.User.Username != "" {
				guild.Presences[i].User.Username = presence.User.Username
			}

			return nil
		}
	}

	guild.Presences = append(guild.Presences, presence)
	return nil
}

func (s *State) PresenceRemove(guildID string, presence *Presence) error {
	if s == nil {
		return ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for i, p := range guild.Presences {
		if p.User.ID == presence.User.ID {
			guild.Presences = append(guild.Presences[:i], guild.Presences[i+1:]...)
			return nil
		}
	}

	return ErrStateNotFound
}

func (s *State) Presence(guildID, userID string) (*Presence, error) {
	if s == nil {
		return nil, ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, p := range guild.Presences {
		if p.User.ID == userID {
			return p, nil
		}
	}

	return nil, ErrStateNotFound
}


func (s *State) MemberAdd(member *Member) error {
	if s == nil {
		return ErrNilState
	}

	guild, err := s.Guild(member.GuildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	members, ok := s.memberMap[member.GuildID]
	if !ok {
		return ErrStateNotFound
	}

	m, ok := members[member.User.ID]
	if !ok {
		members[member.User.ID] = member
		guild.Members = append(guild.Members, member)
	} else {
		if member.JoinedAt == "" {
			member.JoinedAt = m.JoinedAt
		}
		*m = *member
	}

	return nil
}

func (s *State) MemberRemove(member *Member) error {
	if s == nil {
		return ErrNilState
	}

	guild, err := s.Guild(member.GuildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	members, ok := s.memberMap[member.GuildID]
	if !ok {
		return ErrStateNotFound
	}

	_, ok = members[member.User.ID]
	if !ok {
		return ErrStateNotFound
	}
	delete(members, member.User.ID)

	for i, m := range guild.Members {
		if m.User.ID == member.User.ID {
			guild.Members = append(guild.Members[:i], guild.Members[i+1:]...)
			return nil
		}
	}

	return ErrStateNotFound
}

func (s *State) Member(guildID, userID string) (*Member, error) {
	if s == nil {
		return nil, ErrNilState
	}

	s.RLock()
	defer s.RUnlock()

	members, ok := s.memberMap[guildID]
	if !ok {
		return nil, ErrStateNotFound
	}

	m, ok := members[userID]
	if ok {
		return m, nil
	}

	return nil, ErrStateNotFound
}

func (s *State) RoleAdd(guildID string, role *Role) error {
	if s == nil {
		return ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for i, r := range guild.Roles {
		if r.ID == role.ID {
			guild.Roles[i] = role
			return nil
		}
	}

	guild.Roles = append(guild.Roles, role)
	return nil
}

func (s *State) RoleRemove(guildID, roleID string) error {
	if s == nil {
		return ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for i, r := range guild.Roles {
		if r.ID == roleID {
			guild.Roles = append(guild.Roles[:i], guild.Roles[i+1:]...)
			return nil
		}
	}

	return ErrStateNotFound
}

func (s *State) Role(guildID, roleID string) (*Role, error) {
	if s == nil {
		return nil, ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	s.RLock()
	defer s.RUnlock()

	for _, r := range guild.Roles {
		if r.ID == roleID {
			return r, nil
		}
	}

	return nil, ErrStateNotFound
}

func (s *State) ChannelAdd(channel *Channel) error {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	if c, ok := s.channelMap[channel.ID]; ok {
		if channel.Messages == nil {
			channel.Messages = c.Messages
		}
		if channel.PermissionOverwrites == nil {
			channel.PermissionOverwrites = c.PermissionOverwrites
		}

		*c = *channel
		return nil
	}

	if channel.Type == ChannelTypeDM || channel.Type == ChannelTypeGroupDM {
		s.PrivateChannels = append(s.PrivateChannels, channel)
	} else {
		guild, ok := s.guildMap[channel.GuildID]
		if !ok {
			return ErrStateNotFound
		}

		guild.Channels = append(guild.Channels, channel)
	}

	s.channelMap[channel.ID] = channel

	return nil
}

func (s *State) ChannelRemove(channel *Channel) error {
	if s == nil {
		return ErrNilState
	}

	_, err := s.Channel(channel.ID)
	if err != nil {
		return err
	}

	if channel.Type == ChannelTypeDM || channel.Type == ChannelTypeGroupDM {
		s.Lock()
		defer s.Unlock()

		for i, c := range s.PrivateChannels {
			if c.ID == channel.ID {
				s.PrivateChannels = append(s.PrivateChannels[:i], s.PrivateChannels[i+1:]...)
				break
			}
		}
	} else {
		guild, err := s.Guild(channel.GuildID)
		if err != nil {
			return err
		}

		s.Lock()
		defer s.Unlock()

		for i, c := range guild.Channels {
			if c.ID == channel.ID {
				guild.Channels = append(guild.Channels[:i], guild.Channels[i+1:]...)
				break
			}
		}
	}

	delete(s.channelMap, channel.ID)

	return nil
}

func (s *State) Channel(channelID string) (*Channel, error) {
	if s == nil {
		return nil, ErrNilState
	}

	s.RLock()
	defer s.RUnlock()

	if c, ok := s.channelMap[channelID]; ok {
		return c, nil
	}

	return nil, ErrStateNotFound
}

func (s *State) Emoji(guildID, emojiID string) (*Emoji, error) {
	if s == nil {
		return nil, ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	s.RLock()
	defer s.RUnlock()

	for _, e := range guild.Emojis {
		if e.ID == emojiID {
			return e, nil
		}
	}

	return nil, ErrStateNotFound
}

func (s *State) EmojiAdd(guildID string, emoji *Emoji) error {
	if s == nil {
		return ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for i, e := range guild.Emojis {
		if e.ID == emoji.ID {
			guild.Emojis[i] = emoji
			return nil
		}
	}

	guild.Emojis = append(guild.Emojis, emoji)
	return nil
}

func (s *State) EmojisAdd(guildID string, emojis []*Emoji) error {
	for _, e := range emojis {
		if err := s.EmojiAdd(guildID, e); err != nil {
			return err
		}
	}
	return nil
}

func (s *State) MessageAdd(message *Message) error {
	if s == nil {
		return ErrNilState
	}

	c, err := s.Channel(message.ChannelID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for _, m := range c.Messages {
		if m.ID == message.ID {
			if message.Content != "" {
				m.Content = message.Content
			}
			if message.EditedTimestamp != "" {
				m.EditedTimestamp = message.EditedTimestamp
			}
			if message.Mentions != nil {
				m.Mentions = message.Mentions
			}
			if message.Embeds != nil {
				m.Embeds = message.Embeds
			}
			if message.Attachments != nil {
				m.Attachments = message.Attachments
			}
			if message.Timestamp != "" {
				m.Timestamp = message.Timestamp
			}
			if message.Author != nil {
				m.Author = message.Author
			}
			if message.Components != nil {
				m.Components = message.Components
			}

			return nil
		}
	}

	c.Messages = append(c.Messages, message)

	if len(c.Messages) > s.MaxMessageCount {
		c.Messages = c.Messages[len(c.Messages)-s.MaxMessageCount:]
	}
	return nil
}

func (s *State) MessageRemove(message *Message) error {
	if s == nil {
		return ErrNilState
	}

	return s.messageRemoveByID(message.ChannelID, message.ID)
}

func (s *State) messageRemoveByID(channelID, messageID string) error {
	c, err := s.Channel(channelID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	for i, m := range c.Messages {
		if m.ID == messageID {
			c.Messages = append(c.Messages[:i], c.Messages[i+1:]...)
			return nil
		}
	}

	return ErrStateNotFound
}

func (s *State) voiceStateUpdate(update *VoiceStateUpdate) error {
	guild, err := s.Guild(update.GuildID)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	if update.ChannelID == "" {
		for i, state := range guild.VoiceStates {
			if state.UserID == update.UserID {
				guild.VoiceStates = append(guild.VoiceStates[:i], guild.VoiceStates[i+1:]...)
				return nil
			}
		}
	} else {
		for i, state := range guild.VoiceStates {
			if state.UserID == update.UserID {
				guild.VoiceStates[i] = update.VoiceState
				return nil
			}
		}

		guild.VoiceStates = append(guild.VoiceStates, update.VoiceState)
	}

	return nil
}

func (s *State) VoiceState(guildID, userID string) (*VoiceState, error) {
	if s == nil {
		return nil, ErrNilState
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, state := range guild.VoiceStates {
		if state.UserID == userID {
			return state, nil
		}
	}

	return nil, ErrStateNotFound
}

func (s *State) Message(channelID, messageID string) (*Message, error) {
	if s == nil {
		return nil, ErrNilState
	}

	c, err := s.Channel(channelID)
	if err != nil {
		return nil, err
	}

	s.RLock()
	defer s.RUnlock()

	for _, m := range c.Messages {
		if m.ID == messageID {
			return m, nil
		}
	}

	return nil, ErrStateNotFound
}

func (s *State) onReady(se *Session, r *Ready) (err error) {
	if s == nil {
		return ErrNilState
	}

	s.Lock()
	defer s.Unlock()

	if !se.StateEnabled {
		ready := Ready{
			Version:   r.Version,
			SessionID: r.SessionID,
			User:      r.User,
		}

		s.Ready = ready

		return nil
	}

	s.Ready = *r

	for _, g := range s.Guilds {
		s.guildMap[g.ID] = g
		s.createMemberMap(g)

		for _, c := range g.Channels {
			s.channelMap[c.ID] = c
		}
	}

	for _, c := range s.PrivateChannels {
		s.channelMap[c.ID] = c
	}

	return nil
}

func (s *State) OnInterface(se *Session, i interface{}) (err error) {
	if s == nil {
		return ErrNilState
	}

	r, ok := i.(*Ready)
	if ok {
		return s.onReady(se, r)
	}

	if !se.StateEnabled {
		return nil
	}

	switch t := i.(type) {
	case *GuildCreate:
		err = s.GuildAdd(t.Guild)
	case *GuildUpdate:
		err = s.GuildAdd(t.Guild)
	case *GuildDelete:
		err = s.GuildRemove(t.Guild)
	case *GuildMemberAdd:
		guild, err := s.Guild(t.Member.GuildID)
		if err != nil {
			return err
		}
		guild.MemberCount++

		if s.TrackMembers {
			err = s.MemberAdd(t.Member)
		}
	case *GuildMemberUpdate:
		if s.TrackMembers {
			err = s.MemberAdd(t.Member)
		}
	case *GuildMemberRemove:
		guild, err := s.Guild(t.Member.GuildID)
		if err != nil {
			return err
		}
		guild.MemberCount--

		if s.TrackMembers {
			err = s.MemberRemove(t.Member)
		}
	case *GuildMembersChunk:
		if s.TrackMembers {
			for i := range t.Members {
				t.Members[i].GuildID = t.GuildID
				err = s.MemberAdd(t.Members[i])
			}
		}

		if s.TrackPresences {
			for _, p := range t.Presences {
				err = s.PresenceAdd(t.GuildID, p)
			}
		}
	case *GuildRoleCreate:
		if s.TrackRoles {
			err = s.RoleAdd(t.GuildID, t.Role)
		}
	case *GuildRoleUpdate:
		if s.TrackRoles {
			err = s.RoleAdd(t.GuildID, t.Role)
		}
	case *GuildRoleDelete:
		if s.TrackRoles {
			err = s.RoleRemove(t.GuildID, t.RoleID)
		}
	case *GuildEmojisUpdate:
		if s.TrackEmojis {
			err = s.EmojisAdd(t.GuildID, t.Emojis)
		}
	case *ChannelCreate:
		if s.TrackChannels {
			err = s.ChannelAdd(t.Channel)
		}
	case *ChannelUpdate:
		if s.TrackChannels {
			err = s.ChannelAdd(t.Channel)
		}
	case *ChannelDelete:
		if s.TrackChannels {
			err = s.ChannelRemove(t.Channel)
		}
	case *MessageCreate:
		if s.MaxMessageCount != 0 {
			err = s.MessageAdd(t.Message)
		}
	case *MessageUpdate:
		if s.MaxMessageCount != 0 {
			var old *Message
			old, err = s.Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.MessageAdd(t.Message)
		}
	case *MessageDelete:
		if s.MaxMessageCount != 0 {
			var old *Message
			old, err = s.Message(t.ChannelID, t.ID)
			if err == nil {
				oldCopy := *old
				t.BeforeDelete = &oldCopy
			}

			err = s.MessageRemove(t.Message)
		}
	case *MessageDeleteBulk:
		if s.MaxMessageCount != 0 {
			for _, mID := range t.Messages {
				s.messageRemoveByID(t.ChannelID, mID)
			}
		}
	case *VoiceStateUpdate:
		if s.TrackVoice {
			var old *VoiceState
			old, err = s.VoiceState(t.GuildID, t.UserID)
			if err == nil {
				oldCopy := *old
				t.BeforeUpdate = &oldCopy
			}

			err = s.voiceStateUpdate(t)
		}
	case *PresenceUpdate:
		if s.TrackPresences {
			s.PresenceAdd(t.GuildID, &t.Presence)
		}
		if s.TrackMembers {
			if t.Status == StatusOffline {
				return
			}

			var m *Member
			m, err = s.Member(t.GuildID, t.User.ID)

			if err != nil {
				m = &Member{
					GuildID: t.GuildID,
					User:    t.User,
				}
			} else {
				if t.User.Username != "" {
					m.User.Username = t.User.Username
				}
			}

			err = s.MemberAdd(m)
		}

	}

	return
}

func (s *State) UserChannelPermissions(userID, channelID string) (apermissions int64, err error) {
	if s == nil {
		return 0, ErrNilState
	}

	channel, err := s.Channel(channelID)
	if err != nil {
		return
	}

	guild, err := s.Guild(channel.GuildID)
	if err != nil {
		return
	}

	member, err := s.Member(guild.ID, userID)
	if err != nil {
		return
	}

	return memberPermissions(guild, channel, userID, member.Roles), nil
}

func (s *State) MessagePermissions(message *Message) (apermissions int64, err error) {
	if s == nil {
		return 0, ErrNilState
	}

	if message.Author == nil || message.Member == nil {
		return 0, ErrMessageIncompletePermissions
	}

	channel, err := s.Channel(message.ChannelID)
	if err != nil {
		return
	}

	guild, err := s.Guild(channel.GuildID)
	if err != nil {
		return
	}

	return memberPermissions(guild, channel, message.Author.ID, message.Member.Roles), nil
}

func (s *State) UserColor(userID, channelID string) int {
	if s == nil {
		return 0
	}

	channel, err := s.Channel(channelID)
	if err != nil {
		return 0
	}

	guild, err := s.Guild(channel.GuildID)
	if err != nil {
		return 0
	}

	member, err := s.Member(guild.ID, userID)
	if err != nil {
		return 0
	}

	return firstRoleColorColor(guild, member.Roles)
}

func (s *State) MessageColor(message *Message) int {
	if s == nil {
		return 0
	}

	if message.Member == nil || message.Member.Roles == nil {
		return 0
	}

	channel, err := s.Channel(message.ChannelID)
	if err != nil {
		return 0
	}

	guild, err := s.Guild(channel.GuildID)
	if err != nil {
		return 0
	}

	return firstRoleColorColor(guild, message.Member.Roles)
}

func firstRoleColorColor(guild *Guild, memberRoles []string) int {
	roles := Roles(guild.Roles)
	sort.Sort(roles)

	for _, role := range roles {
		for _, roleID := range memberRoles {
			if role.ID == roleID {
				if role.Color != 0 {
					return role.Color
				}
			}
		}
	}

	for _, role := range roles {
		if role.ID == guild.ID {
			return role.Color
		}
	}

	return 0
}
