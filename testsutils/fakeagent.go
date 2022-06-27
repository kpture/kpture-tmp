package testutils

import (
	"context"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/kpture/agent/api/capture"
)

type FakeAgent struct {
	file *os.File
	name string
}

func NewFakeAgent(name string) *FakeAgent {
	return &FakeAgent{name: name}
}

func (fa *FakeAgent) OpenCapture(filepath string, ctx context.Context, req *capture.CaptureDescriptor) error {
	var err error
	fa.file, err = os.Open(filepath)

	return err
}

func (fa *FakeAgent) Name() string {
	return fa.name
}

func (fa *FakeAgent) Packets(ctx context.Context, req *capture.CaptureDescriptor) (chan gopacket.Packet, error) {
	h, err := pcap.OpenOfflineFile(fa.file)
	if err != nil {
		return nil, err
	}

	packetSource := gopacket.NewPacketSource(h, h.LinkType())
	packetSource.DecodeOptions.Lazy = true
	packetSource.DecodeOptions.NoCopy = true

	ch := packetSource.Packets()

	go func() {
		<-ctx.Done()

		fa.file.Close()
	}()

	return ch, nil
}
