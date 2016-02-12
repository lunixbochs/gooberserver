package main

type Battle struct {
	Config     Config
	Type       string
	NatType    string
	Password   string
	Port       int
	MaxPlayers int
	Spectators int
	Hashcode   string
	Rank       int
	MapHash    int
	Map        string
	Mod        string
	Users      []*User
	Host       *User
	StartRects []string

	PendingUsers []*User
	AuthedUsers  []*User

	Engine  string
	Version string

	Bots                []string
	ScriptTags          []string
	ReplayScript        string
	Replay              string
	SendingReplayScript bool
	Locked              bool
}
