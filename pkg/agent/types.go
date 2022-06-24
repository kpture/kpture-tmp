package agent

type Type string

const (
	TypeKubernetes Type = "kubernetes"
	TypeLinux      Type = "linux"
	TypeContainer  Type = "container"
)

type Status string

const (
	StatusUP   Status = "up"
	StatusDown Status = "down"
)

type Info struct {
	Metadata Metadata `json:"metadata"`
	Status   Status   `json:"status"`
	Errors   []string `json:"errors,omitempty"`
	PacketNb uint64   `json:"packetNb"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Type      Type   `json:"system"`
	TargetURL string `json:"targetUrl"`
}

const (
	bufChanSize int = 1024
	errChanSize int = 1000
)
