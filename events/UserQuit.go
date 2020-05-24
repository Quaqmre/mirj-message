package events

type UserQuit struct {
	ClientID int32
	Name     string
	Key      string
}
