package routers

type Router interface {
	AddPort(config PortConfig) error
	RemovePort(port int) error
	CheckPort(port int) (bool, error)
}

type PortConfig struct {
	Name      string
	Enabled   bool
	Interface string
	SrcPort   int
	DstPort   int
	SrcIp     string
	Protocol  string
}
