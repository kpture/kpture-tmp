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
	Name      string  `json:"name,omitempty"`
	Namespace string  `json:"namespace,omitempty"`
	Type      Type    `json:"system,omitempty"`
	TargetURL string  `json:"targetUrl,omitempty"`
	Status    Status  `json:"status,omitempty"`
	Errors    []error `json:"errors,omitempty"`
}

const (
	bufChanSize int = 1024
)
