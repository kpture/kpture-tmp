package agent

import (
	"context"
	"fmt"
	"newproxy/pkg/logger"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/kpture/agent/api/capture"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CaptureSocket Represent a single target capture.
type CaptureSocket struct {
	logger    *logrus.Entry
	AgentInfo Info
}

const (
	AgentTypek8s int = iota
	AgentTypeContainer
)

type Info struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Type      int    `json:"type,omitempty"`
	TargetURL string `json:"targetUrl,omitempty"`
}

func NewCaptureSocket(info Info) *CaptureSocket {
	return &CaptureSocket{
		logger:    logger.NewLogger("capture"),
		AgentInfo: info,
	}
}

func (c *CaptureSocket) Info() Info {
	return c.AgentInfo
}

// Packets handle the packet reception.
func (c *CaptureSocket) Packets(
	captureCtx context.Context,
	request *capture.CaptureDescriptor,
) (chan gopacket.Packet, error) {

	var err error
	pkch := make(chan gopacket.Packet, 1024)
	conn, err := grpc.Dial(c.AgentInfo.TargetURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error dialing with target %s %w", c.AgentInfo.TargetURL, err)
	}
	socketCapture, err := capture.NewKptureClient(conn).Capture(captureCtx, request)
	if err != nil {
		return nil, fmt.Errorf("error starting capture %s %w", c.AgentInfo.TargetURL, err)
	}

	go func() {
		for {
			select {
			case <-captureCtx.Done():
				conn.Close()
				return
			default:
				packet, err := socketCapture.Recv()
				if err != nil {
					return
				}
				info := gopacket.CaptureInfo{
					Timestamp:      time.Now(),
					CaptureLength:  int(packet.CaptureInfo.CaptureLength),
					Length:         int(packet.CaptureInfo.Length),
					InterfaceIndex: int(packet.CaptureInfo.InterfaceIndex),
				}
				pk := gopacket.NewPacket(packet.Data, layers.LayerTypeEthernet, gopacket.NoCopy)
				pk.Metadata().CaptureInfo = info
				pkch <- pk
			}
		}
	}()

	return pkch, nil
}
