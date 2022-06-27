package capture

type KptureStatus string

const (
	KptureStatusNotStarted KptureStatus = "not_started"
	KptureStatusRunning    KptureStatus = "running"
	KptureStatusStopped    KptureStatus = "stopped"
	KptureStatusWriting    KptureStatus = "writing"
	KptureStatusError      KptureStatus = "error"
	KptureStatusTerminated KptureStatus = "terminated"
)

const (
	bufChanSize    int    = 1024
	snapshotLen    int32  = 1024
	pcapFileHeader uint32 = 1024
)
