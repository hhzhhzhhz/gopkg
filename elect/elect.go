package election

type FunType string

const (
	Follower FunType = "Follower"
	Leader   FunType = "Leader"
)

type Election interface {
	IsLeader() (bool, string)
	Register(ft FunType, f func())
	Run() error
	Close() error
}
