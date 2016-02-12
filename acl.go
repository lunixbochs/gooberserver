package main

import (
	"errors"
)

type AclFlag int

func (a AclFlag) IsChannel() bool {
	return (a == AclChanUser ||
		a == AclChanOp ||
		a == AclChanOwner)
}

func (a AclFlag) IsBattle() bool {
	return a == AclBattleUser || a == AclBattleHost
}

const (
	AclDisabled AclFlag = iota
	AclEveryone
	AclUnauthed
	AclAgreement
	AclUser
	AclMod
	AclAdmin

	AclBattleUser
	AclBattleHost

	AclChanUser
	AclChanOp
	AclChanOwner
)

// This is the default ACL mapping for the server.
// TODO: test to warn if the same command is in two ACLs
var DefaultAcl = map[AclFlag][]string{
	// banned or pre-login acls
	AclDisabled: {},
	AclEveryone: {
		"EXIT",
		"PING",
		"LISTCOMPFLAGS",

		// encryption
		"GETPUBLICKEY",
		"GETSIGNEDMSG",
		"SETSHAREDKEY",
		"ACKSHAREDKEY",
	},
	AclUnauthed: {
		"LOGIN",
		"REGISTER",
	},
	AclAgreement: {
		"CONFIRMAGREEMENT",
	},

	// logged-in acls below this line
	AclUser: {
		// channel
		"CHANNELS",
		"JOIN",
		"LEAVE",
		// other users
		"SAYPRIVATE",
		"SAYPRIVATEEX",
		// ignore
		"IGNORE",
		"UNIGNORE",
		"IGNORELIST",
		// friend
		"FRIENDREQUEST",
		"ACCEPTFRIENDREQUEST",
		"DECLINEFRIENDREQUEST",
		"UNFRIEND",
		"FRIENDLIST",
		"FRIENDREQUESTLIST",
		// channel subscriptions
		"SUBSCRIBE",
		"UNSUBSCRIBE",
		"LISTSUBSCRIPTIONS",
		// meta
		"CHANGEEMAIL",
		"CHANGEPASSWORD",
		"GETINGAMETIME",
		"GETREGISTRATIONDATE",
		"MYSTATUS",
		"PORTTEST",
		"RENAMEACCOUNT",
		// outside battle
		"OPENBATTLE",
		"JOINBATTLE",
	},

	// channel-specific acls
	AclChanUser: {
		"MUTELIST",
		"SAY",
		"SAYEX",
	},
	AclChanOp: {
		"MUTE",
		"CHANNELMESSAGE",
		"FORCELEAVECHANNEL",
		"CHANNELTOPIC",
		"SETCHANNELKEY",
		"UNMUTE",
	},
	// channel owner and op have the same commands
	AclChanOwner: {},

	// battle acls below this line
	AclBattleUser: {
		"ADDBOT",
		"HANDICAP",
		"LEAVEBATTLE",
		"MYBATTLESTATUS",
		"REMOVEBOT",
		"RING",
		"SAYBATTLE",
		"SAYBATTLEEX",
		"SAYBATTLEPRIVATE",
		"SAYBATTLEPRIVATEEX",
		"UPDATEBOT",
	},
	AclBattleHost: {
		"ADDSTARTRECT",
		"DISABLEUNITS",
		"ENABLEUNITS",
		"ENABLEALLUNITS",
		"FORCEALLYNO",
		"FORCESPECTATORMODE",
		"FORCETEAMCOLOR",
		"FORCETEAMNO",
		"FORCEJOINBATTLE",
		"JOINBATTLEACCEPT",
		"JOINBATTLEDENY",
		"KICKFROMBATTLE",
		"OPENBATTLE",
		"REMOVESCRIPTTAGS",
		"REMOVESTARTRECT",
		"SETSCRIPTTAGS",
		"UPDATEBATTLEINFO",
	},

	// mod/admin stuff below this line
	AclMod: {
		"BAN",
		"BANIP",
		"UNBAN",
		"UNBANIP",
		"BANLIST",
		"CHANGEACCOUNTPASS",
		"KICKUSER",
		"FINDIP",
		"GETIP",
		"GETLASTLOGINTIME",
		"GETUSERID",
		"SETBOTMODE",
		"GETLOBBYVERSION",
	},
	AclAdmin: {
		// server
		"ADMINBROADCAST",
		"BROADCAST",
		"BROADCASTEX",
		"RELOAD",
		"CLEANUP",
		"SETLATESTSPRINGVERSION",
		// users
		"GETLASTLOGINTIME",
		"GETACCOUNTACCESS",
		"FORCEJOIN",
		"SETACCESS",
	},
}

var DefaultReverseAcl = make(map[string]AclFlag)

func init() {
	for flag, commands := range DefaultAcl {
		for _, cmd := range commands {
			DefaultReverseAcl[cmd] = flag
		}
	}
}

var CommandNotFoundErr = errors.New("command not found")
var AccessDeniedErr = errors.New("access denied")
var NotInChanErr = errors.New("not in channel")
var NotInBattleErr = errors.New("not in battle")

type AclFlags []AclFlag

// nil means "access granted"
// channel is the Acl for the referenced channel (if appropriate)
// battle is the Acl flag for the referenced battle (if appropriate)
func (acl AclFlags) Check(cmd string, channel AclFlag, battle AclFlag) error {
	if flag, ok := DefaultReverseAcl[cmd]; !ok {
		return CommandNotFoundErr
	} else {
		if flag.IsChannel() {
			if channel == 0 {
				return NotInChanErr
			}
			// if our channel permissions meet or exceed required, acl passes
			if channel >= flag || acl.IsMod() {
				return nil
			}
		} else if flag.IsBattle() {
			// exception so mods can kick hosts/users from battle
			if cmd != "KICKFROMBATTLE" && battle == 0 {
				return NotInBattleErr
			}
			// if our battle permissions meet or exceed required, acl passes
			if battle >= flag || acl.IsMod() {
				return nil
			}
		} else {
			for _, v := range acl {
				if flag == v {
					return nil
				}
			}
		}
		return AccessDeniedErr
	}
}

func (acl AclFlags) Search(flag AclFlag) bool {
	for _, v := range acl {
		if v == flag {
			return true
		}
	}
	return false
}

func (acl AclFlags) IsAdmin() bool {
	return acl.Search(AclAdmin)
}

func (acl AclFlags) IsMod() bool {
	return acl.IsAdmin() || acl.Search(AclMod)
}
