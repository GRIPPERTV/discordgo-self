package discordgoself

const (
	channelCreateEventType            = "CHANNEL_CREATE"
	channelDeleteEventType            = "CHANNEL_DELETE"
	channelPinsUpdateEventType        = "CHANNEL_PINS_UPDATE"
	channelUpdateEventType            = "CHANNEL_UPDATE"
	connectEventType                  = "__CONNECT__"
	disconnectEventType               = "__DISCONNECT__"
	eventEventType                    = "__EVENT__"
	guildBanAddEventType              = "GUILD_BAN_ADD"
	guildBanRemoveEventType           = "GUILD_BAN_REMOVE"
	guildCreateEventType              = "GUILD_CREATE"
	guildDeleteEventType              = "GUILD_DELETE"
	guildEmojisUpdateEventType        = "GUILD_EMOJIS_UPDATE"
	guildIntegrationsUpdateEventType  = "GUILD_INTEGRATIONS_UPDATE"
	guildMemberAddEventType           = "GUILD_MEMBER_ADD"
	guildMemberRemoveEventType        = "GUILD_MEMBER_REMOVE"
	guildMemberUpdateEventType        = "GUILD_MEMBER_UPDATE"
	guildMembersChunkEventType        = "GUILD_MEMBERS_CHUNK"
	guildRoleCreateEventType          = "GUILD_ROLE_CREATE"
	guildRoleDeleteEventType          = "GUILD_ROLE_DELETE"
	guildRoleUpdateEventType          = "GUILD_ROLE_UPDATE"
	guildUpdateEventType              = "GUILD_UPDATE"
	messageAckEventType               = "MESSAGE_ACK"
	messageCreateEventType            = "MESSAGE_CREATE"
	messageDeleteEventType            = "MESSAGE_DELETE"
	messageDeleteBulkEventType        = "MESSAGE_DELETE_BULK"
	messageReactionAddEventType       = "MESSAGE_REACTION_ADD"
	messageReactionRemoveEventType    = "MESSAGE_REACTION_REMOVE"
	messageReactionRemoveAllEventType = "MESSAGE_REACTION_REMOVE_ALL"
	messageUpdateEventType            = "MESSAGE_UPDATE"
	presenceUpdateEventType           = "PRESENCE_UPDATE"
	presencesReplaceEventType         = "PRESENCES_REPLACE"
	rateLimitEventType                = "__RATE_LIMIT__"
	readyEventType                    = "READY"
	relationshipAddEventType          = "RELATIONSHIP_ADD"
	relationshipRemoveEventType       = "RELATIONSHIP_REMOVE"
	resumedEventType                  = "RESUMED"
	typingStartEventType              = "TYPING_START"
	userGuildSettingsUpdateEventType  = "USER_GUILD_SETTINGS_UPDATE"
	userNoteUpdateEventType           = "USER_NOTE_UPDATE"
	userSettingsUpdateEventType       = "USER_SETTINGS_UPDATE"
	userUpdateEventType               = "USER_UPDATE"
	voiceServerUpdateEventType        = "VOICE_SERVER_UPDATE"
	voiceStateUpdateEventType         = "VOICE_STATE_UPDATE"
	webhooksUpdateEventType           = "WEBHOOKS_UPDATE"
)

type channelCreateEventHandler func(*Session, *ChannelCreate)

func (eh channelCreateEventHandler) Type() string {
	return channelCreateEventType
}

func (eh channelCreateEventHandler) New() interface{} {
	return &ChannelCreate{}
}

func (eh channelCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelCreate); ok {
		eh(s, t)
	}
}

type channelDeleteEventHandler func(*Session, *ChannelDelete)

func (eh channelDeleteEventHandler) Type() string {
	return channelDeleteEventType
}

func (eh channelDeleteEventHandler) New() interface{} {
	return &ChannelDelete{}
}

func (eh channelDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelDelete); ok {
		eh(s, t)
	}
}

type channelPinsUpdateEventHandler func(*Session, *ChannelPinsUpdate)

func (eh channelPinsUpdateEventHandler) Type() string {
	return channelPinsUpdateEventType
}

func (eh channelPinsUpdateEventHandler) New() interface{} {
	return &ChannelPinsUpdate{}
}

func (eh channelPinsUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelPinsUpdate); ok {
		eh(s, t)
	}
}

type channelUpdateEventHandler func(*Session, *ChannelUpdate)

func (eh channelUpdateEventHandler) Type() string {
	return channelUpdateEventType
}

func (eh channelUpdateEventHandler) New() interface{} {
	return &ChannelUpdate{}
}

func (eh channelUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelUpdate); ok {
		eh(s, t)
	}
}

type connectEventHandler func(*Session, *Connect)

func (eh connectEventHandler) Type() string {
	return connectEventType
}

func (eh connectEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Connect); ok {
		eh(s, t)
	}
}

type disconnectEventHandler func(*Session, *Disconnect)

func (eh disconnectEventHandler) Type() string {
	return disconnectEventType
}

func (eh disconnectEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Disconnect); ok {
		eh(s, t)
	}
}

type eventEventHandler func(*Session, *Event)

func (eh eventEventHandler) Type() string {
	return eventEventType
}

func (eh eventEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Event); ok {
		eh(s, t)
	}
}

type guildBanAddEventHandler func(*Session, *GuildBanAdd)

func (eh guildBanAddEventHandler) Type() string {
	return guildBanAddEventType
}

func (eh guildBanAddEventHandler) New() interface{} {
	return &GuildBanAdd{}
}

func (eh guildBanAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildBanAdd); ok {
		eh(s, t)
	}
}

type guildBanRemoveEventHandler func(*Session, *GuildBanRemove)

func (eh guildBanRemoveEventHandler) Type() string {
	return guildBanRemoveEventType
}

func (eh guildBanRemoveEventHandler) New() interface{} {
	return &GuildBanRemove{}
}

func (eh guildBanRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildBanRemove); ok {
		eh(s, t)
	}
}

type guildCreateEventHandler func(*Session, *GuildCreate)

func (eh guildCreateEventHandler) Type() string {
	return guildCreateEventType
}

func (eh guildCreateEventHandler) New() interface{} {
	return &GuildCreate{}
}

func (eh guildCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildCreate); ok {
		eh(s, t)
	}
}

type guildDeleteEventHandler func(*Session, *GuildDelete)

func (eh guildDeleteEventHandler) Type() string {
	return guildDeleteEventType
}

func (eh guildDeleteEventHandler) New() interface{} {
	return &GuildDelete{}
}

func (eh guildDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildDelete); ok {
		eh(s, t)
	}
}

type guildEmojisUpdateEventHandler func(*Session, *GuildEmojisUpdate)

func (eh guildEmojisUpdateEventHandler) Type() string {
	return guildEmojisUpdateEventType
}

func (eh guildEmojisUpdateEventHandler) New() interface{} {
	return &GuildEmojisUpdate{}
}

func (eh guildEmojisUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildEmojisUpdate); ok {
		eh(s, t)
	}
}

type guildIntegrationsUpdateEventHandler func(*Session, *GuildIntegrationsUpdate)

func (eh guildIntegrationsUpdateEventHandler) Type() string {
	return guildIntegrationsUpdateEventType
}

func (eh guildIntegrationsUpdateEventHandler) New() interface{} {
	return &GuildIntegrationsUpdate{}
}

func (eh guildIntegrationsUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildIntegrationsUpdate); ok {
		eh(s, t)
	}
}

type guildMemberAddEventHandler func(*Session, *GuildMemberAdd)

func (eh guildMemberAddEventHandler) Type() string {
	return guildMemberAddEventType
}

func (eh guildMemberAddEventHandler) New() interface{} {
	return &GuildMemberAdd{}
}

func (eh guildMemberAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMemberAdd); ok {
		eh(s, t)
	}
}

type guildMemberRemoveEventHandler func(*Session, *GuildMemberRemove)

func (eh guildMemberRemoveEventHandler) Type() string {
	return guildMemberRemoveEventType
}

func (eh guildMemberRemoveEventHandler) New() interface{} {
	return &GuildMemberRemove{}
}

func (eh guildMemberRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMemberRemove); ok {
		eh(s, t)
	}
}

type guildMemberUpdateEventHandler func(*Session, *GuildMemberUpdate)

func (eh guildMemberUpdateEventHandler) Type() string {
	return guildMemberUpdateEventType
}

func (eh guildMemberUpdateEventHandler) New() interface{} {
	return &GuildMemberUpdate{}
}

func (eh guildMemberUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMemberUpdate); ok {
		eh(s, t)
	}
}

type guildMembersChunkEventHandler func(*Session, *GuildMembersChunk)

func (eh guildMembersChunkEventHandler) Type() string {
	return guildMembersChunkEventType
}

func (eh guildMembersChunkEventHandler) New() interface{} {
	return &GuildMembersChunk{}
}

func (eh guildMembersChunkEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildMembersChunk); ok {
		eh(s, t)
	}
}

type guildRoleCreateEventHandler func(*Session, *GuildRoleCreate)

func (eh guildRoleCreateEventHandler) Type() string {
	return guildRoleCreateEventType
}

func (eh guildRoleCreateEventHandler) New() interface{} {
	return &GuildRoleCreate{}
}

func (eh guildRoleCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildRoleCreate); ok {
		eh(s, t)
	}
}

type guildRoleDeleteEventHandler func(*Session, *GuildRoleDelete)

func (eh guildRoleDeleteEventHandler) Type() string {
	return guildRoleDeleteEventType
}

func (eh guildRoleDeleteEventHandler) New() interface{} {
	return &GuildRoleDelete{}
}

func (eh guildRoleDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildRoleDelete); ok {
		eh(s, t)
	}
}

type guildRoleUpdateEventHandler func(*Session, *GuildRoleUpdate)

func (eh guildRoleUpdateEventHandler) Type() string {
	return guildRoleUpdateEventType
}

func (eh guildRoleUpdateEventHandler) New() interface{} {
	return &GuildRoleUpdate{}
}

func (eh guildRoleUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildRoleUpdate); ok {
		eh(s, t)
	}
}

type guildUpdateEventHandler func(*Session, *GuildUpdate)

func (eh guildUpdateEventHandler) Type() string {
	return guildUpdateEventType
}

func (eh guildUpdateEventHandler) New() interface{} {
	return &GuildUpdate{}
}

func (eh guildUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GuildUpdate); ok {
		eh(s, t)
	}
}

type messageAckEventHandler func(*Session, *MessageAck)

func (eh messageAckEventHandler) Type() string {
	return messageAckEventType
}

func (eh messageAckEventHandler) New() interface{} {
	return &MessageAck{}
}

func (eh messageAckEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageAck); ok {
		eh(s, t)
	}
}

type messageCreateEventHandler func(*Session, *MessageCreate)

func (eh messageCreateEventHandler) Type() string {
	return messageCreateEventType
}

func (eh messageCreateEventHandler) New() interface{} {
	return &MessageCreate{}
}

func (eh messageCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageCreate); ok {
		eh(s, t)
	}
}

type messageDeleteEventHandler func(*Session, *MessageDelete)

func (eh messageDeleteEventHandler) Type() string {
	return messageDeleteEventType
}

func (eh messageDeleteEventHandler) New() interface{} {
	return &MessageDelete{}
}

func (eh messageDeleteEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageDelete); ok {
		eh(s, t)
	}
}

type messageDeleteBulkEventHandler func(*Session, *MessageDeleteBulk)

func (eh messageDeleteBulkEventHandler) Type() string {
	return messageDeleteBulkEventType
}

func (eh messageDeleteBulkEventHandler) New() interface{} {
	return &MessageDeleteBulk{}
}

func (eh messageDeleteBulkEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageDeleteBulk); ok {
		eh(s, t)
	}
}

type messageReactionAddEventHandler func(*Session, *MessageReactionAdd)

func (eh messageReactionAddEventHandler) Type() string {
	return messageReactionAddEventType
}

func (eh messageReactionAddEventHandler) New() interface{} {
	return &MessageReactionAdd{}
}

func (eh messageReactionAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageReactionAdd); ok {
		eh(s, t)
	}
}

type messageReactionRemoveEventHandler func(*Session, *MessageReactionRemove)

func (eh messageReactionRemoveEventHandler) Type() string {
	return messageReactionRemoveEventType
}

func (eh messageReactionRemoveEventHandler) New() interface{} {
	return &MessageReactionRemove{}
}

func (eh messageReactionRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageReactionRemove); ok {
		eh(s, t)
	}
}

type messageReactionRemoveAllEventHandler func(*Session, *MessageReactionRemoveAll)

func (eh messageReactionRemoveAllEventHandler) Type() string {
	return messageReactionRemoveAllEventType
}

func (eh messageReactionRemoveAllEventHandler) New() interface{} {
	return &MessageReactionRemoveAll{}
}

func (eh messageReactionRemoveAllEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageReactionRemoveAll); ok {
		eh(s, t)
	}
}

type messageUpdateEventHandler func(*Session, *MessageUpdate)

func (eh messageUpdateEventHandler) Type() string {
	return messageUpdateEventType
}

func (eh messageUpdateEventHandler) New() interface{} {
	return &MessageUpdate{}
}

func (eh messageUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageUpdate); ok {
		eh(s, t)
	}
}

type presenceUpdateEventHandler func(*Session, *PresenceUpdate)

func (eh presenceUpdateEventHandler) Type() string {
	return presenceUpdateEventType
}

func (eh presenceUpdateEventHandler) New() interface{} {
	return &PresenceUpdate{}
}

func (eh presenceUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*PresenceUpdate); ok {
		eh(s, t)
	}
}

type presencesReplaceEventHandler func(*Session, *PresencesReplace)

func (eh presencesReplaceEventHandler) Type() string {
	return presencesReplaceEventType
}

func (eh presencesReplaceEventHandler) New() interface{} {
	return &PresencesReplace{}
}

func (eh presencesReplaceEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*PresencesReplace); ok {
		eh(s, t)
	}
}

type rateLimitEventHandler func(*Session, *RateLimit)

func (eh rateLimitEventHandler) Type() string {
	return rateLimitEventType
}

func (eh rateLimitEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*RateLimit); ok {
		eh(s, t)
	}
}

type readyEventHandler func(*Session, *Ready)

func (eh readyEventHandler) Type() string {
	return readyEventType
}

func (eh readyEventHandler) New() interface{} {
	return &Ready{}
}

func (eh readyEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Ready); ok {
		eh(s, t)
	}
}

type relationshipAddEventHandler func(*Session, *RelationshipAdd)

func (eh relationshipAddEventHandler) Type() string {
	return relationshipAddEventType
}

func (eh relationshipAddEventHandler) New() interface{} {
	return &RelationshipAdd{}
}

func (eh relationshipAddEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*RelationshipAdd); ok {
		eh(s, t)
	}
}

type relationshipRemoveEventHandler func(*Session, *RelationshipRemove)

func (eh relationshipRemoveEventHandler) Type() string {
	return relationshipRemoveEventType
}

func (eh relationshipRemoveEventHandler) New() interface{} {
	return &RelationshipRemove{}
}

func (eh relationshipRemoveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*RelationshipRemove); ok {
		eh(s, t)
	}
}

type resumedEventHandler func(*Session, *Resumed)

func (eh resumedEventHandler) Type() string {
	return resumedEventType
}

func (eh resumedEventHandler) New() interface{} {
	return &Resumed{}
}

func (eh resumedEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*Resumed); ok {
		eh(s, t)
	}
}

type typingStartEventHandler func(*Session, *TypingStart)

func (eh typingStartEventHandler) Type() string {
	return typingStartEventType
}

func (eh typingStartEventHandler) New() interface{} {
	return &TypingStart{}
}

func (eh typingStartEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*TypingStart); ok {
		eh(s, t)
	}
}

type userGuildSettingsUpdateEventHandler func(*Session, *UserGuildSettingsUpdate)

func (eh userGuildSettingsUpdateEventHandler) Type() string {
	return userGuildSettingsUpdateEventType
}

func (eh userGuildSettingsUpdateEventHandler) New() interface{} {
	return &UserGuildSettingsUpdate{}
}

func (eh userGuildSettingsUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*UserGuildSettingsUpdate); ok {
		eh(s, t)
	}
}

type userNoteUpdateEventHandler func(*Session, *UserNoteUpdate)

func (eh userNoteUpdateEventHandler) Type() string {
	return userNoteUpdateEventType
}

func (eh userNoteUpdateEventHandler) New() interface{} {
	return &UserNoteUpdate{}
}

func (eh userNoteUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*UserNoteUpdate); ok {
		eh(s, t)
	}
}

type userSettingsUpdateEventHandler func(*Session, *UserSettingsUpdate)

func (eh userSettingsUpdateEventHandler) Type() string {
	return userSettingsUpdateEventType
}

func (eh userSettingsUpdateEventHandler) New() interface{} {
	return &UserSettingsUpdate{}
}

func (eh userSettingsUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*UserSettingsUpdate); ok {
		eh(s, t)
	}
}

type userUpdateEventHandler func(*Session, *UserUpdate)

func (eh userUpdateEventHandler) Type() string {
	return userUpdateEventType
}

func (eh userUpdateEventHandler) New() interface{} {
	return &UserUpdate{}
}

func (eh userUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*UserUpdate); ok {
		eh(s, t)
	}
}

type voiceServerUpdateEventHandler func(*Session, *VoiceServerUpdate)

func (eh voiceServerUpdateEventHandler) Type() string {
	return voiceServerUpdateEventType
}

func (eh voiceServerUpdateEventHandler) New() interface{} {
	return &VoiceServerUpdate{}
}

func (eh voiceServerUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*VoiceServerUpdate); ok {
		eh(s, t)
	}
}

type voiceStateUpdateEventHandler func(*Session, *VoiceStateUpdate)

func (eh voiceStateUpdateEventHandler) Type() string {
	return voiceStateUpdateEventType
}

func (eh voiceStateUpdateEventHandler) New() interface{} {
	return &VoiceStateUpdate{}
}

func (eh voiceStateUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*VoiceStateUpdate); ok {
		eh(s, t)
	}
}

type webhooksUpdateEventHandler func(*Session, *WebhooksUpdate)

func (eh webhooksUpdateEventHandler) Type() string {
	return webhooksUpdateEventType
}

func (eh webhooksUpdateEventHandler) New() interface{} {
	return &WebhooksUpdate{}
}

func (eh webhooksUpdateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*WebhooksUpdate); ok {
		eh(s, t)
	}
}

func handlerForInterface(handler interface{}) EventHandler {
	switch v := handler.(type) {
	case func(*Session, interface{}):
		return interfaceEventHandler(v)
	case func(*Session, *ChannelCreate):
		return channelCreateEventHandler(v)
	case func(*Session, *ChannelDelete):
		return channelDeleteEventHandler(v)
	case func(*Session, *ChannelPinsUpdate):
		return channelPinsUpdateEventHandler(v)
	case func(*Session, *ChannelUpdate):
		return channelUpdateEventHandler(v)
	case func(*Session, *Connect):
		return connectEventHandler(v)
	case func(*Session, *Disconnect):
		return disconnectEventHandler(v)
	case func(*Session, *Event):
		return eventEventHandler(v)
	case func(*Session, *GuildBanAdd):
		return guildBanAddEventHandler(v)
	case func(*Session, *GuildBanRemove):
		return guildBanRemoveEventHandler(v)
	case func(*Session, *GuildCreate):
		return guildCreateEventHandler(v)
	case func(*Session, *GuildDelete):
		return guildDeleteEventHandler(v)
	case func(*Session, *GuildEmojisUpdate):
		return guildEmojisUpdateEventHandler(v)
	case func(*Session, *GuildIntegrationsUpdate):
		return guildIntegrationsUpdateEventHandler(v)
	case func(*Session, *GuildMemberAdd):
		return guildMemberAddEventHandler(v)
	case func(*Session, *GuildMemberRemove):
		return guildMemberRemoveEventHandler(v)
	case func(*Session, *GuildMemberUpdate):
		return guildMemberUpdateEventHandler(v)
	case func(*Session, *GuildMembersChunk):
		return guildMembersChunkEventHandler(v)
	case func(*Session, *GuildRoleCreate):
		return guildRoleCreateEventHandler(v)
	case func(*Session, *GuildRoleDelete):
		return guildRoleDeleteEventHandler(v)
	case func(*Session, *GuildRoleUpdate):
		return guildRoleUpdateEventHandler(v)
	case func(*Session, *GuildUpdate):
		return guildUpdateEventHandler(v)
	case func(*Session, *MessageAck):
		return messageAckEventHandler(v)
	case func(*Session, *MessageCreate):
		return messageCreateEventHandler(v)
	case func(*Session, *MessageDelete):
		return messageDeleteEventHandler(v)
	case func(*Session, *MessageDeleteBulk):
		return messageDeleteBulkEventHandler(v)
	case func(*Session, *MessageReactionAdd):
		return messageReactionAddEventHandler(v)
	case func(*Session, *MessageReactionRemove):
		return messageReactionRemoveEventHandler(v)
	case func(*Session, *MessageReactionRemoveAll):
		return messageReactionRemoveAllEventHandler(v)
	case func(*Session, *MessageUpdate):
		return messageUpdateEventHandler(v)
	case func(*Session, *PresenceUpdate):
		return presenceUpdateEventHandler(v)
	case func(*Session, *PresencesReplace):
		return presencesReplaceEventHandler(v)
	case func(*Session, *RateLimit):
		return rateLimitEventHandler(v)
	case func(*Session, *Ready):
		return readyEventHandler(v)
	case func(*Session, *RelationshipAdd):
		return relationshipAddEventHandler(v)
	case func(*Session, *RelationshipRemove):
		return relationshipRemoveEventHandler(v)
	case func(*Session, *Resumed):
		return resumedEventHandler(v)
	case func(*Session, *TypingStart):
		return typingStartEventHandler(v)
	case func(*Session, *UserGuildSettingsUpdate):
		return userGuildSettingsUpdateEventHandler(v)
	case func(*Session, *UserNoteUpdate):
		return userNoteUpdateEventHandler(v)
	case func(*Session, *UserSettingsUpdate):
		return userSettingsUpdateEventHandler(v)
	case func(*Session, *UserUpdate):
		return userUpdateEventHandler(v)
	case func(*Session, *VoiceServerUpdate):
		return voiceServerUpdateEventHandler(v)
	case func(*Session, *VoiceStateUpdate):
		return voiceStateUpdateEventHandler(v)
	case func(*Session, *WebhooksUpdate):
		return webhooksUpdateEventHandler(v)
	}

	return nil
}

func init() {
	registerInterfaceProvider(channelCreateEventHandler(nil))
	registerInterfaceProvider(channelDeleteEventHandler(nil))
	registerInterfaceProvider(channelPinsUpdateEventHandler(nil))
	registerInterfaceProvider(channelUpdateEventHandler(nil))
	registerInterfaceProvider(guildBanAddEventHandler(nil))
	registerInterfaceProvider(guildBanRemoveEventHandler(nil))
	registerInterfaceProvider(guildCreateEventHandler(nil))
	registerInterfaceProvider(guildDeleteEventHandler(nil))
	registerInterfaceProvider(guildEmojisUpdateEventHandler(nil))
	registerInterfaceProvider(guildIntegrationsUpdateEventHandler(nil))
	registerInterfaceProvider(guildMemberAddEventHandler(nil))
	registerInterfaceProvider(guildMemberRemoveEventHandler(nil))
	registerInterfaceProvider(guildMemberUpdateEventHandler(nil))
	registerInterfaceProvider(guildMembersChunkEventHandler(nil))
	registerInterfaceProvider(guildRoleCreateEventHandler(nil))
	registerInterfaceProvider(guildRoleDeleteEventHandler(nil))
	registerInterfaceProvider(guildRoleUpdateEventHandler(nil))
	registerInterfaceProvider(guildUpdateEventHandler(nil))
	registerInterfaceProvider(messageAckEventHandler(nil))
	registerInterfaceProvider(messageCreateEventHandler(nil))
	registerInterfaceProvider(messageDeleteEventHandler(nil))
	registerInterfaceProvider(messageDeleteBulkEventHandler(nil))
	registerInterfaceProvider(messageReactionAddEventHandler(nil))
	registerInterfaceProvider(messageReactionRemoveEventHandler(nil))
	registerInterfaceProvider(messageReactionRemoveAllEventHandler(nil))
	registerInterfaceProvider(messageUpdateEventHandler(nil))
	registerInterfaceProvider(presenceUpdateEventHandler(nil))
	registerInterfaceProvider(presencesReplaceEventHandler(nil))
	registerInterfaceProvider(readyEventHandler(nil))
	registerInterfaceProvider(relationshipAddEventHandler(nil))
	registerInterfaceProvider(relationshipRemoveEventHandler(nil))
	registerInterfaceProvider(resumedEventHandler(nil))
	registerInterfaceProvider(typingStartEventHandler(nil))
	registerInterfaceProvider(userGuildSettingsUpdateEventHandler(nil))
	registerInterfaceProvider(userNoteUpdateEventHandler(nil))
	registerInterfaceProvider(userSettingsUpdateEventHandler(nil))
	registerInterfaceProvider(userUpdateEventHandler(nil))
	registerInterfaceProvider(voiceServerUpdateEventHandler(nil))
	registerInterfaceProvider(voiceStateUpdateEventHandler(nil))
	registerInterfaceProvider(webhooksUpdateEventHandler(nil))
}
