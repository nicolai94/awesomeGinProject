package enums

type UserStatus int

const (
	Active UserStatus = iota
	Inactive
	Suspended
	Deleted
)

func (s UserStatus) String() string {
	return [...]string{"Active", "Inactive", "Suspended", "Deleted"}[s]
}
