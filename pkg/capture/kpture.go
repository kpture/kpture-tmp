package capture

import (
	"context"
	"io/fs"
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
	Info() *agent.Info
	HealthCheck()
}

// Kpture represent a Kpture
// It can contains multiples targets.
type Kpture struct {
	agents []Agent `json:"-"`

	AgentsInfos []*agent.Info `json:"agents,omitempty"`

	Name        string `json:"name,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	archivePath string
	CaptureInfo CaptureInfo  `json:"captureInfo,omitempty"`
	Status      KptureStatus `json:"status,omitempty"`
	logger      *logrus.Entry
	stopCh      chan bool
	basePath    string
}

type CaptureInfo struct {
	Size     int `json:"size,omitempty"`
	PacketNB int `json:"packetNb,omitempty"`
}

func NewKpture(name, archivePath string, agents []Agent) (*Kpture, error) {
	k := &Kpture{}
	k.Name = name
	k.UUID = uuid.New().String()
	k.logger = logger.NewLogger("kpture")
	k.agents = agents
	k.AgentsInfos = []*agent.Info{}
	k.stopCh = make(chan bool)
	k.archivePath = archivePath
	err := os.MkdirAll(filepath.Join(os.TempDir(), k.UUID, k.Name), fs.ModePerm)
	k.basePath = filepath.Join(os.TempDir(), k.UUID, k.Name)

	for _, currA := range agents {
		k.AgentsInfos = append(k.AgentsInfos, currA.Info())
	}

	return k, err
}

func (k *Kpture) Start() {
	ctx, ctxCancel := context.WithCancel(context.Background())

	go func() {
		<-k.stopCh
		ctxCancel()
	}()

	globalChan := make(chan gopacket.Packet, bufChanSize)

	err := k.storePackets(k.basePath, "kpture", globalChan)
	if err != nil {
		panic(err)
	}

	statsChan := make(chan gopacket.Packet, bufChanSize)

	err = k.stats(statsChan)
	if err != nil {
		panic(err)
	}

	k.handleAgents(ctx, globalChan, statsChan)

	go func() {
		<-k.stopCh
		close(globalChan)
	}()

	k.Status = KptureStatusRunning
}

func (k *Kpture) handleAgents(ctx context.Context, chans ...chan gopacket.Packet) {
	for _, agent := range k.agents {
		agentChan := make(chan gopacket.Packet, bufChanSize)

		err := k.storePackets(filepath.Join(k.basePath, agent.Info().Name), agent.Info().Name, agentChan)
		if err != nil {
			panic(err)
		}

		packet, err := agent.Packets(ctx, &capture.CaptureDescriptor{
			InterfaceName: "eth0",
			SnapshotLen:   snapshotLen,
			Promiscuous:   false,
			Filter:        "port not 10000",
		})
		if err != nil {
			panic(err)
		}

		go func() {
			for p := range packet {
				select {
				case <-ctx.Done():
					close(agentChan)

					return
				default:
					for i := 0; i < len(chans); i++ {
						chans[i] <- p
					}

					agentChan <- p
				}
			}
		}()
	}
}

func (k *Kpture) Stop() error {
	k.Status = KptureStatusStopped
	k.stopCh <- true
	k.logger.Debug("kpture stopped")
	k.Status = KptureStatusWriting

	buf, err := k.createTar()
	if err != nil {
		k.logger.Error(err)

		return err
	}

	if err := k.writeFile(buf); err != nil {
		k.logger.Error(err)
		k.Status = KptureStatusError

		return err
	}

	k.Status = KptureStatusTerminated

	err = k.MarshalDescription()
	if err != nil {
		k.logger.Error(err)
		k.Status = KptureStatusError

		return err
	}

	return nil
}
