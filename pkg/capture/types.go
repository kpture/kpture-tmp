package capture

type Pod struct {
	PodMetadata PodMetadata `json:"metadata"`
	PodStatus   PodStatus   `json:"status"`
	AgentStatus string      `json:"agentStatus"`
}

type PodStatus struct {
	Phase   string `json:"phase,omitempty"`
	Started string `json:"started,omitempty"`
}

type PodMetadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Filter    string `json:"filter,omitempty"`
}

type CaptureStatus struct {
	CaptureState CaptureState `json:"capture_state,omitempty"`
	Desciption   string       `json:"description,omitempty"`
}

type CaptureState string

const (
	CaptureStatusNotStarted CaptureState = "KPTURE_NOT_STARTED"
	CaptureStatusStarted    CaptureState = "KPTURE_RUNNING"
	CaptureStatusWriting    CaptureState = "KPTURE_WRITING_ARCHIVE"
	CaptureStatusReady      CaptureState = "KPTURE_READY"
	CaptureStatusError      CaptureState = "KPTURE_ERROR"
)
