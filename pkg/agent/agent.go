package agent

import (
	"context"
	"newproxy/pkg/logger"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/kpture/agent/api/capture"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CaptureSocket Represent a single target capture.
type CaptureSocket struct {
	logger    *logrus.Entry
	AgentInfo *Info
}

// NewCaptureSocket creates a grpc agent client
func NewCaptureSocket(m Metadata) *CaptureSocket {
	c := &CaptureSocket{
		logger: logger.NewLogger("capture"),
		AgentInfo: &Info{
			Metadata: m,
		},
	}
	c.HealthCheck()

	return c
}

func (c *CaptureSocket) Info() *Info {
	return c.AgentInfo
}

// Packets handle the packet reception.
func (c *CaptureSocket) Packets(
	captureCtx context.Context,
	request *capture.CaptureDescriptor,
) (chan gopacket.Packet, error) {
	var err error

	pkch := make(chan gopacket.Packet, bufChanSize)

	conn, err := grpc.Dial(c.AgentInfo.Metadata.TargetURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.WithMessagef(err, "error dialing with target %s", c.AgentInfo.Metadata.Name)
	}

	socketCapture, err := capture.NewKptureClient(conn).Capture(captureCtx, request)
	if err != nil {
		return nil, errors.WithMessagef(err, "error starting capture %s", c.AgentInfo.Metadata.Name)
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
					c.AgentInfo.Errors = append(c.AgentInfo.Errors, err.Error())

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
				c.AgentInfo.PacketNb++
				pkch <- pk
			}
		}
	}()

	return pkch, nil
}

func (c *CaptureSocket) HealthCheck() {
	conn, err := grpc.Dial(c.AgentInfo.Metadata.TargetURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.AgentInfo.Status = StatusDown
		c.AgentInfo.Errors = append(c.AgentInfo.Errors, err.Error())

		return
	}

	if _, err := capture.NewKptureClient(conn).Health(context.Background(), &capture.Empty{}); err != nil {
		c.AgentInfo.Status = StatusDown
		c.AgentInfo.Errors = append(c.AgentInfo.Errors, err.Error())

		return
	}

	c.AgentInfo.Status = StatusUP
}
