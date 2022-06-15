package capture

import (
	"bytes"
	"context"
	"fmt"
	"newproxy/pkg/logger"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/kpture/agent/api/capture"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//captureSocket Represent a single target capture
type captureSocket struct {
	//Pod informations
	Pod PodMetadata `json:"pod"`
	//Number of packet captured
	PacketsNumber int64 `json:"packetsNumber"`
	//capture size
	size uint64

	//packet buffer
	packets *bytes.Buffer `json:"-"`
	//packet writer
	packetWriter *pcapgo.Writer `json:"-"`

	//capture context
	ctx context.Context

	//grpc client toward the select pod
	socketCapture capture.Kpture_CaptureClient `json:"-"`
	conn          *grpc.ClientConn

	//logger
	logger *logrus.Entry
}

func NewCaptureSocket(pod PodMetadata, podIP string) (*captureSocket, error) {
	c := &captureSocket{}
	c.Pod = pod
	c.logger = logger.NewLogger("capture")
	c.logger.Info("Created")
	c.packets = bytes.NewBuffer([]byte{})

	c.packetWriter = pcapgo.NewWriter(c.packets)
	c.packetWriter.WriteFileHeader(1024, layers.LinkTypeEthernet)
	var err error
	c.conn, err = grpc.Dial(fmt.Sprintf("%s:%d", podIP, 10000), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return c, nil
}

//OpenSocket open the socket
func (c *captureSocket) Open(captureCtx context.Context, request *capture.CaptureDescriptor) (err error) {
	c.ctx = captureCtx
	c.socketCapture, err = capture.NewKptureClient(c.conn).Capture(captureCtx, request)
	if err != nil {
		return err
	}
	return nil
}

//GetFileBuf get the packet buffer
func (c *captureSocket) GetFileBuf() *bytes.Buffer {
	return c.packets
}

//ReceivePackets handle the packet reception
func (c *captureSocket) ListenSocket(pkCh chan gopacket.Packet, wg *sync.WaitGroup) {
	c.logger.Info("Listening")
	defer wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			c.packets.Reset()
			c.conn.Close()
			return
		default:
			s, err := c.socketCapture.Recv()
			if err != nil {
				return
			}
			c.packetWriter.WritePacket(gopacket.CaptureInfo{
				Timestamp:      time.Now(),
				CaptureLength:  int(s.CaptureInfo.CaptureLength),
				Length:         int(s.CaptureInfo.Length),
				InterfaceIndex: int(s.CaptureInfo.InterfaceIndex),
			}, s.Data)

			if pkCh != nil {
				curPkt := gopacket.NewPacket(s.Data, layers.LayerTypeEthernet, gopacket.NoCopy)
				select {
				case pkCh <- curPkt:
				default:
					fmt.Println("no message sent")
				}
			}
			c.PacketsNumber++
			c.size += uint64(len(s.Data))
		}
	}
}
