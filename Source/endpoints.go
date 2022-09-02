package discordgoself

import "strconv"

var APIVersion = "8"

var (
	// Host
	EndpointAPI    = "https://discord.com/api/v" + APIVersion + "/"
	EndpointCDN    = "https://cdn.discordapp.com/"
	EndpointStatus = "https://status.discord.com/api/v2/"

	// Status
	EndpointSm                         = EndpointStatus + "scheduled-maintenances/"
	EndpointSmActive                   = EndpointSm     + "active.json"
	EndpointSmUpcoming                 = EndpointSm     + "upcoming.json"

	// CDN
	EndpointCDNAttachments             = EndpointCDN + "attachments/"
	EndpointCDNAvatars                 = EndpointCDN + "avatars/"
	EndpointCDNIcons                   = EndpointCDN + "icons/"
	EndpointCDNSplashes                = EndpointCDN + "splashes/"
	EndpointCDNChannelIcons            = EndpointCDN + "channel-icons/"
	EndpointCDNBanners                 = EndpointCDN + "banners/"

	// API
	EndpointGuilds                     = EndpointAPI + "guilds/"
	EndpointChannels                   = EndpointAPI + "channels/"
	EndpointUsers                      = EndpointAPI + "users/"
	EndpointGateway                    = EndpointAPI + "gateway"
	EndpointWebhooks                   = EndpointAPI + "webhooks/"
	EndpointAuth                       = EndpointAPI + "auth/"
	EndpointVoice                      = EndpointAPI + "/voice/"
	EndpointTrack                      = EndpointAPI + "track"
	EndpointSso                        = EndpointAPI + "sso"
	EndpointReport                     = EndpointAPI + "report"
	EndpointIntegrations               = EndpointAPI + "integrations"
	EndpointGuildCreate                = EndpointAPI + "guilds"
	EndpointApplications               = EndpointAPI + "applications"
	EndpointOAuth2                     = EndpointAPI + "oauth2/"
	EndpointTutorial                   = EndpointAPI + "tutorial/"

	// Auth
	EndpointLogin                      = EndpointAuth + "login"
	EndpointLogout                     = EndpointAuth + "logout"
	EndpointVerify                     = EndpointAuth + "verify"
	EndpointVerifyResend               = EndpointAuth + "verify/resend"
	EndpointForgotPassword             = EndpointAuth + "forgot"
	EndpointResetPassword              = EndpointAuth + "reset"
	EndpointRegister                   = EndpointAuth + "register"

	// Voice
	EndpointVoiceRegions               = EndpointVoice + "regions"
	EndpointVoiceIce                   = EndpointVoice + "ice"

	// Tutorial
	EndpointTutorialIndicators         = EndpointTutorial + "indicators"

	// OAuth2
	EndpointOAuth2Applications         = EndpointOAuth2 + "applications"

	// User
	EndpointUser                       = func(uID string)      string { return EndpointUsers + uID                                  }
	EndpointUserProfile                = func(uID string)      string { return EndpointUsers + uID + "/profile"                     }
	EndpointUserSettings               = func(uID string)      string { return EndpointUsers + uID + "/settings"                    }
	EndpointUserGuilds                 = func(uID string)      string { return EndpointUsers + uID + "/guilds"                      }
	EndpointUserGuild                  = func(uID, gID string) string { return EndpointUsers + uID + "/guilds/" + gID               }
	EndpointUserGuildSettings          = func(uID, gID string) string { return EndpointUsers + uID + "/guilds/" + gID + "/settings" }
	EndpointUserChannels               = func(uID string)      string { return EndpointUsers + uID + "/channels"                    }
	EndpointUserDevices                = func(uID string)      string { return EndpointUsers + uID + "/devices"                     }
	EndpointUserConnections            = func(uID string)      string { return EndpointUsers + uID + "/connections"                 }
	EndpointUserAvatar                 = func(uID, aID string) string { return EndpointCDNAvatars + uID + "/" + aID + ".png"        }
	EndpointUserAvatarAnimated         = func(uID, aID string) string { return EndpointCDNAvatars + uID + "/" + aID + ".gif"        }
	EndpointUserNotes                  = func(uID string)      string { return EndpointUsers + "@me/notes/" + uID                   }

	// Guild
	EndpointGuild                = func(gID string)           string { return EndpointGuilds + gID                                       }
	EndpointGuildPreview         = func(gID string)           string { return EndpointGuilds + gID + "/preview"                          }
	EndpointGuildChannels        = func(gID string)           string { return EndpointGuilds + gID + "/channels"                         }
	EndpointGuildMembers         = func(gID string)           string { return EndpointGuilds + gID + "/members"                          }
	EndpointGuildMember          = func(gID, uID string)      string { return EndpointGuilds + gID + "/members/" + uID                   }
	EndpointGuildMemberRole      = func(gID, uID, rID string) string { return EndpointGuilds + gID + "/members/" + uID + "/roles/" + rID }
	EndpointGuildBans            = func(gID string)           string { return EndpointGuilds + gID + "/bans"                             }
	EndpointGuildBan             = func(gID, uID string)      string { return EndpointGuilds + gID + "/bans/" + uID                      }
	EndpointGuildIntegrations    = func(gID string)           string { return EndpointGuilds + gID + "/integrations"                     }
	EndpointGuildIntegration     = func(gID, iID string)      string { return EndpointGuilds + gID + "/integrations/" + iID              }
	EndpointGuildIntegrationSync = func(gID, iID string)      string { return EndpointGuilds + gID + "/integrations/" + iID + "/sync"    }
	EndpointGuildRoles           = func(gID string)           string { return EndpointGuilds + gID + "/roles"                            }
	EndpointGuildRole            = func(gID, rID string)      string { return EndpointGuilds + gID + "/roles/" + rID                     }
	EndpointGuildInvites         = func(gID string)           string { return EndpointGuilds + gID + "/invites"                          }
	EndpointGuildWidget          = func(gID string)           string { return EndpointGuilds + gID + "/widget"                           }
	EndpointGuildPrune           = func(gID string)           string { return EndpointGuilds + gID + "/prune"                            }
	EndpointGuildIcon            = func(gID, hash string)     string { return EndpointCDNIcons + gID + "/" + hash + ".png"               }
	EndpointGuildIconAnimated    = func(gID, hash string)     string { return EndpointCDNIcons + gID + "/" + hash + ".gif"               }
	EndpointGuildSplash          = func(gID, hash string)     string { return EndpointCDNSplashes + gID + "/" + hash + ".png"            }
	EndpointGuildWebhooks        = func(gID string)           string { return EndpointGuilds + gID + "/webhooks"                         }
	EndpointGuildAuditLogs       = func(gID string)           string { return EndpointGuilds + gID + "/audit-logs"                       }
	EndpointGuildEmojis          = func(gID string)           string { return EndpointGuilds + gID + "/emojis"                           }
	EndpointGuildEmoji           = func(gID, eID string)      string { return EndpointGuilds + gID + "/emojis/" + eID                    }
	EndpointGuildBanner          = func(gID, hash string)     string { return EndpointCDNBanners + gID + "/" + hash + ".png"             }
	EndpointGuildEmbed           = EndpointGuildWidget

	// Channel
	EndpointChannel                   = func(cID string)       string { return EndpointChannels + cID                                   }
	EndpointChannelPermissions        = func(cID string)       string { return EndpointChannels + cID + "/permissions"                  }
	EndpointChannelPermission         = func(cID, tID string)  string { return EndpointChannels + cID + "/permissions/" + tID           }
	EndpointChannelInvites            = func(cID string)       string { return EndpointChannels + cID + "/invites"                      }
	EndpointChannelTyping             = func(cID string)       string { return EndpointChannels + cID + "/typing"                       }
	EndpointChannelMessages           = func(cID string)       string { return EndpointChannels + cID + "/messages"                     }
	EndpointChannelMessage            = func(cID, mID string)  string { return EndpointChannels + cID + "/messages/" + mID              }
	EndpointChannelMessageAck         = func(cID, mID string)  string { return EndpointChannels + cID + "/messages/" + mID + "/ack"     }
	EndpointChannelMessagesBulkDelete = func(cID string)       string { return EndpointChannel(cID) + "/messages/bulk-delete"           }
	EndpointChannelMessagesPins       = func(cID string)       string { return EndpointChannel(cID) + "/pins"                           }
	EndpointChannelMessagePin         = func(cID, mID string)  string { return EndpointChannel(cID) + "/pins/" + mID                    }
	EndpointChannelMessageCrosspost   = func(cID, mID string)  string { return EndpointChannel(cID) + "/messages/" + mID + "/crosspost" }
	EndpointChannelFollow             = func(cID string)       string { return EndpointChannel(cID) + "/followers"                      }
	EndpointChannelWebhooks           = func(cID string)       string { return EndpointChannel(cID) + "/webhooks"                       }
	EndpointGroupIcon                 = func(cID, hash string) string { return EndpointCDNChannelIcons + cID + "/" + hash + ".png"      }

	// Webhook
	EndpointWebhook        = func(wID string)                   string { return EndpointWebhooks + wID                                      }
	EndpointWebhookToken   = func(wID, token string)            string { return EndpointWebhooks + wID + "/" + token                        }
	EndpointWebhookMessage = func(wID, token, messageID string) string { return EndpointWebhookToken(wID, token) + "/messages/" + messageID }

	// Reaction
	EndpointMessageReactionsAll = func(cID, mID string)           string { return EndpointChannelMessage(cID, mID) + "/reactions"        }
	EndpointMessageReactions    = func(cID, mID, eID string)      string { return EndpointChannelMessage(cID, mID) + "/reactions/" + eID }
	EndpointMessageReaction     = func(cID, mID, eID, uID string) string { return EndpointMessageReactions(cID, mID, eID) + "/" + uID    }

	// ???
	EndpointApplicationGlobalCommands = func(aID string)           string { return EndpointApplication(aID) + "/commands"                    }
	EndpointApplicationGlobalCommand  = func(aID, cID string)      string { return EndpointApplicationGlobalCommands(aID) + "/" + cID        }
	EndpointApplicationGuildCommands  = func(aID, gID string)      string { return EndpointApplication(aID) + "/guilds/" + gID + "/commands" }
	EndpointApplicationGuildCommand   = func(aID, gID, cID string) string { return EndpointApplicationGuildCommands(aID, gID) + "/" + cID    }

	EndpointFollowupMessage        = func(aID, iToken string)      string { return EndpointWebhookToken(aID, iToken)        }
	EndpointFollowupMessageActions = func(aID, iToken, mID string) string { return EndpointWebhookMessage(aID, iToken, mID) }

	// Relationship
	EndpointRelationships       = func()           string { return EndpointUsers + "@me" + "/relationships" }
	EndpointRelationship        = func(uID string) string { return EndpointRelationships() + "/" + uID      }
	EndpointRelationshipsMutual = func(uID string) string { return EndpointUsers + uID + "/relationships"   }

	// Invite
	EndpointInvite = func(iID string) string { return EndpointAPI + "invites/" + iID }

	// JoinIntegration
	EndpointIntegrationsJoin = func(iID string) string { return EndpointAPI + "integrations/" + iID + "/join" }

	// Emoji
	EndpointEmoji         = func(eID string) string { return EndpointCDN + "emojis/" + eID + ".png" }
	EndpointEmojiAnimated = func(eID string) string { return EndpointCDN + "emojis/" + eID + ".gif" }

	// Application
	EndpointApplication             = func(aID string) string { return EndpointApplications + "/" + aID                   }
	EndpointOAuth2Application       = func(aID string) string { return EndpointOAuth2Applications + "/" + aID             }
	EndpointOAuth2ApplicationsBot   = func(aID string) string { return EndpointOAuth2Applications + "/" + aID + "/bot"    }
	EndpointOAuth2ApplicationAssets = func(aID string) string { return EndpointOAuth2Applications + "/" + aID + "/assets" }

	// Default User Avatar
	EndpointDefaultUserAvatar = func(uDiscriminator string) string {
		uDiscriminatorInt, _ := strconv.Atoi(uDiscriminator)
		return EndpointCDN + "embed/avatars/" + strconv.Itoa(uDiscriminatorInt%5) + ".png"
	}
)
