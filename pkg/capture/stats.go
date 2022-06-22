package capture

import "github.com/google/gopacket"

func (k *Kpture) stats(ch chan gopacket.Packet) error {

	go func() {
		for packet := range ch {
			k.CaptureInfo.PacketNB++
			k.CaptureInfo.Size += len(packet.Data())
		}
	}()

	return nil
}
