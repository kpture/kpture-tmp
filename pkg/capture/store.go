package capture

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func (k *Kpture) storePackets(basepath string, name string, ch chan gopacket.Packet) error {
	err := os.MkdirAll(basepath, 0755)
	if err != nil {
		return err
	}
	location := filepath.Join(basepath, name) + ".pcap"
	file, err := os.Create(location)
	if err != nil {
		fmt.Println("error creating file", err)
		return nil
	}

	w := pcapgo.NewWriter(file)
	err = w.WriteFileHeader(1024, layers.LinkTypeEthernet)
	if err != nil {
		return err
	}

	go func() {
		for packet := range ch {
			fmt.Println("packet", name)

			err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			if err != nil {
				fmt.Println("error writing packet", err)
			}
		}
		file.Close()
	}()

	return nil
}

func (k *Kpture) createTar() error {

	err := os.MkdirAll(filepath.Join(k.archivePath, k.UUID), 0755)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	// walk through every file in the folder
	err = filepath.Walk(k.basePath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		var header *tar.Header // generate tar header
		header, err = tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		newStr := strings.Replace(file, filepath.Join(os.TempDir(), k.UUID), "", -1)
		header.Name = newStr

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}

	fileToWrite, err := os.OpenFile(filepath.Join(k.archivePath, k.UUID, k.Name+".tar.gzip"), os.O_CREATE|os.O_RDWR, os.FileMode(0755))
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(fileToWrite, buf); err != nil {
		panic(err)
	}

	k.Status = KptureStatusTerminated
	//os.RemoveAll(k.UUID)

	return nil
}

func (k *Kpture) MarshalDescription() error {
	bytes, err := json.MarshalIndent(k, "", "    ")
	if err != nil {
		return err
	}
	fileToWrite, err := os.OpenFile(filepath.Join(k.archivePath, k.UUID, "descriptor.json"), os.O_CREATE|os.O_RDWR, os.FileMode(0755))
	if err != nil {
		panic(err)
	}
	_, err = fileToWrite.Write(bytes)
	return err
}
