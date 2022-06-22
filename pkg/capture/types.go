package capture

type KptureStatus string

const (
	KptureStatusRunning    KptureStatus = "running"
	KptureStatusStopped    KptureStatus = "stopped"
	KptureStatusWriting    KptureStatus = "writing"
	KptureStatusError      KptureStatus = "error"
	KptureStatusTerminated KptureStatus = "terminated"
)
