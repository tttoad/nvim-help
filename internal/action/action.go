package action

type Action interface {
	Action() string
	Run([]string) *Result
	Usege() string
}
