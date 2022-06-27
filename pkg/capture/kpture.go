package capture

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"newproxy/pkg/agent"
	"newproxy/pkg/logger"

	"github.com/google/gopacket"
	"github.com/google/uuid"
	"github.com/kpture/agent/api/capture"
	"github.com/pkg/errors"
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
	agents      []Agent
	ProfileName string        `json:"profilName,omitempty"`
	AgentsInfos []*agent.Info `json:"agents,omitempty"`

	Name        string `json:"name,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	archivePath string
	CaptureInfo Info         `json:"captureInfo,omitempty"`
	Status      KptureStatus `json:"status,omitempty"`
	logger      *logrus.Entry
	stopCh      chan bool
	basePath    string
	StartTime   time.Time `json:"startTime,omitempty"`
	StopTime    time.Time `json:"stopTime,omitempty"`
}

type Info struct {
	Size     int `json:"size,omitempty"`
	PacketNB int `json:"packetNb,omitempty"`
}

func NewKpture(kptureName, profileName, archivePath string, agents []Agent) (*Kpture, error) {
	uuid := uuid.New().String()

	err := os.MkdirAll(filepath.Join(os.TempDir(), uuid, kptureName), fs.ModePerm)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create kpture directory")
	}

	kapture := &Kpture{
		Name:        kptureName,
		UUID:        uuid,
		logger:      logger.NewLogger("kpture"),
		agents:      agents,
		ProfileName: profileName,
		AgentsInfos: []*agent.Info{},
		stopCh:      make(chan bool),
		archivePath: archivePath,
		basePath:    filepath.Join(os.TempDir(), profileName, uuid, kptureName),
		CaptureInfo: Info{
			Size:     0,
			PacketNB: 0,
		},
		Status: KptureStatusError,
	}

	for _, currA := range agents {
		kapture.AgentsInfos = append(kapture.AgentsInfos, currA.Info())
	}

	return kapture, nil
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
	k.StartTime = time.Now()
}

func (k *Kpture) handleAgents(ctx context.Context, chans ...chan gopacket.Packet) {
	for _, agent := range k.agents {
		agentChan := make(chan gopacket.Packet, bufChanSize)

		err := k.storePackets(filepath.Join(k.basePath, agent.Info().Metadata.Name), agent.Info().Metadata.Name, agentChan)
		if err != nil {
			panic(err)
		}

		packet, err := agent.Packets(ctx, &capture.CaptureDescriptor{
			Timeout:       -1,
			PacketCount:   0,
			InterfaceName: "eth0",
			SnapshotLen:   snapshotLen,
			Promiscuous:   false,
			Filter:        "port not 10000",
		})
		if err != nil {
			panic(err)
		}

		go func() {
			for currPacket := range packet {
				select {
				case <-ctx.Done():
					close(agentChan)

					return
				default:
					for i := 0; i < len(chans); i++ {
						chans[i] <- currPacket
					}

					agentChan <- currPacket
				}
			}
		}()
	}
}

func (k *Kpture) Stop() error {
	k.Status = KptureStatusStopped
	k.StopTime = time.Now()
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
