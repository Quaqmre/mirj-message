package events

type UserConnected struct {
	ClientID int32
	Name     string
	Key      string
}
