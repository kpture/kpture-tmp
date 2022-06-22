package capture

import (
	"context"
	"newproxy/pkg/agent"
	"newproxy/pkg/logger"
	"os"
	"path/filepath"

	"github.com/google/gopacket"
	"github.com/google/uuid"
	"github.com/kpture/agent/api/capture"
	"github.com/sirupsen/logrus"
)

type Agent interface {
	Packets(ctx context.Context, req *capture.CaptureDescriptor) (chan gopacket.Packet, error)
	Info() agent.Info
}

//Kpture represent a Kpture
//It can contains multiples targets
type Kpture struct {
	Agents      []Agent `json:"-"`
	Name        string  `json:"name,omitempty"`
	UUID        string  `json:"uuid,omitempty"`
	archivePath string
	CaptureInfo CaptureInfo  `json:"capture_info,omitempty"`
	Status      KptureStatus `json:"status,omitempty"`
	logger      *logrus.Entry
	ctx         context.Context
	ctxCancel   context.CancelFunc
	basePath    string
}

type CaptureInfo struct {
	Size     int `json:"size,omitempty"`
	PacketNB int `json:"packet_nb,omitempty"`
}

func NewKpture(name, archivePath string, agents []Agent) (*Kpture, error) {
	k := &Kpture{}
	k.Name = name
	k.UUID = uuid.New().String()
	k.logger = logger.NewLogger("kpture")
	k.Agents = agents
	k.ctx, k.ctxCancel = context.WithCancel(context.Background())
	k.archivePath = archivePath
	err := os.MkdirAll(filepath.Join(os.TempDir(), k.UUID, k.Name), 0755)
	k.basePath = filepath.Join(os.TempDir(), k.UUID, k.Name)
	return k, err
}

func (k *Kpture) Start() {

	globalChan := make(chan gopacket.Packet, 1024)
	err := k.storePackets(k.basePath, "kpture", globalChan)
	if err != nil {
		panic(err)
	}

	statsChannel := make(chan gopacket.Packet, 1024)
	err = k.stats(statsChannel)
	if err != nil {
		panic(err)
	}

	go func() {
		<-k.ctx.Done()
		close(globalChan)
	}()

	for _, agent := range k.Agents {

		agentChan := make(chan gopacket.Packet, 1024)
		err := k.storePackets(filepath.Join(k.basePath, agent.Info().Name), agent.Info().Name, agentChan)
		if err != nil {
			panic(err)
		}

		packet, err := agent.Packets(k.ctx, &capture.CaptureDescriptor{
			InterfaceName: "eth0",
			SnapshotLen:   1024,
			Promiscuous:   false,
			Filter:        "port not 10000",
		})
		if err != nil {
			panic(err)
		}
		go func() {
			for p := range packet {
				select {
				case <-k.ctx.Done():
					close(agentChan)
					return
				default:
					globalChan <- p
					agentChan <- p
					statsChannel <- p
				}
			}
		}()
	}
	k.Status = KptureStatusRunning
}

func (k *Kpture) Stop() error {
	k.Status = KptureStatusStopped
	k.ctxCancel()
	k.Status = KptureStatusWriting
	err := k.createTar()
	if err != nil {
		k.logger.Error(err)
		return err
	}
	err = k.MarshalDescription()
	if err != nil {
		k.logger.Error(err)
		return err
	}
	return nil
}
