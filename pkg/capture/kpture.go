package capture

import (
	"context"
	"fmt"
	"newproxy/pkg/logger"
	"sync"

	humanize "github.com/dustin/go-humanize"
	"github.com/google/gopacket"
	"github.com/google/uuid"
	"github.com/kpture/agent/api/capture"
	"github.com/sirupsen/logrus"
)

//kpture represent a kpture
//It can contains multiples targets
type kpture struct {
	captures        []*captureSocket
	Pods            []PodMetadata `json:"pods,omitempty"`
	Name            string        `json:"name,omitempty"`
	UUID            string        `json:"uuid,omitempty"`
	Status          CaptureStatus `json:"status,omitempty"`
	ArchiveLocation string        `json:"archive_location,omitempty"`
	PacketNB        int64         `json:"packet_nb,omitempty"`
	BufferSizeStr   string        `json:"buffer_size,omitempty"`

	Size       string `json:"size,omitempty"`
	PacketMiss uint64 `json:"packet_miss,omitempty"`

	logger *logrus.Entry
	//Graceful grpc streaming stop function
	kpturectx  context.Context
	grpcCancel context.CancelFunc

	// Asynchronous handling for all the socket
	wg *sync.WaitGroup

	//Global packet channel to handle packets incoming from all targets
	globalChannel    chan gopacket.Packet
	websocketChannel chan gopacket.Packet
}

func NewKpture(name string) *kpture {
	k := &kpture{}
	k.Name = name
	k.UUID = uuid.New().String()
	k.Status.CaptureState = CaptureStatusNotStarted
	k.logger = logger.NewLogger("kpture")
	k.globalChannel = make(chan gopacket.Packet, 1024)
	k.websocketChannel = make(chan gopacket.Packet, 1024)

	return k
}

func (k *kpture) Weboscket() chan gopacket.Packet {
	return k.websocketChannel
}

func (k *kpture) HandleGLobalPacket() {
	for {
		select {
		case packet := <-k.globalChannel:
			k.websocketChannel <- packet
		case <-k.kpturectx.Done():
			return
		}
	}
}

func (k *kpture) AddCapture(capture *captureSocket) {
	//Todo validate capture (conainerID, interfaceName)
	k.captures = append(k.captures, capture)
	k.Pods = append(k.Pods, capture.Pod)
}

func (k *kpture) UpdatePkNumber() int64 {
	//Todo validate capture (conainerID, interfaceName)
	var t int64
	var size uint64
	var capturemiss uint64
	for _, capture := range k.captures {
		t += capture.PacketsNumber
		size += capture.size
	}
	k.PacketNB = t

	k.BufferSizeStr = humanize.Bytes(size)
	k.PacketMiss = capturemiss
	return k.PacketNB
}

func (k *kpture) OpenCaptureSockets() error {
	var err error

	//Create a global cancel func for all socket within this kpture
	k.kpturectx, k.grpcCancel = context.WithCancel(context.Background())

	for _, curr := range k.captures {
		newreq := capture.CaptureDescriptor{
			InterfaceName: "eth0",
			SnapshotLen:   1024,
			Promiscuous:   false,
			Filter:        curr.Pod.Filter,
		}

		fmt.Println("Opening capture socket ", newreq)
		err = curr.Open(k.kpturectx, &newreq)
		if err != nil {
			k.logger.Error("Error opening socket", err)
			return err
		}
	}
	return nil
}

func (k *kpture) StartListenSockets() {
	k.wg = new(sync.WaitGroup)

	go k.HandleGLobalPacket()
	for _, capture := range k.captures {
		k.wg.Add(1)
		go capture.ListenSocket(k.globalChannel, k.wg)
	}
	k.Status.CaptureState = CaptureStatusStarted
}

func (k *kpture) StopKptures() {

	k.grpcCancel()
	k.logger.Debug("Stopping kpture, waiting for all sockets to stop")
	k.wg.Wait()

	k.logger.Debug("All goroutine channel stoped")

	close(k.globalChannel)

}

func (k *kpture) CleanBuffer() {
	k.captures = nil
}
