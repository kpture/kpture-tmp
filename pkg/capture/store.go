package capture

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"github.com/pkg/errors"
)

func (k *Kpture) storePackets(basepath string, name string, ch chan gopacket.Packet) error {
	err := os.MkdirAll(basepath, fs.ModePerm)
	if err != nil {
		return err
	}

	location := filepath.Join(basepath, name) + ".pcap"

	file, err := os.Create(location)
	if err != nil {
		return errors.WithMessage(err, "error creating file")
	}

	w := pcapgo.NewWriter(file)

	err = w.WriteFileHeader(pcapFileHeader, layers.LinkTypeEthernet)
	if err != nil {
		return err
	}

	go func() {
		for packet := range ch {
			err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			if err != nil {
				k.logger.Error(errors.WithMessage(err, "could not write packet"))
			}
		}

		file.Close()
	}()

	return nil
}

func (k *Kpture) createTar() (*bytes.Buffer, error) {
	err := os.MkdirAll(filepath.Join(k.archivePath, k.UUID), fs.ModePerm)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

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

		if strings.HasPrefix(newStr, "/") {
			newStr = strings.TrimPrefix(newStr, "/")
		}

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
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	// if err := zr.Close(); err != nil {
	// 	return nil, err
	// }

	return buf, err
}

func (k *Kpture) writeFile(buf *bytes.Buffer) error {
	location := filepath.Join(k.archivePath, k.UUID, k.Name+".tar")

	fileToWrite, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR, os.FileMode(fs.ModePerm))
	if err != nil {
		return err
	}

	if _, err := io.Copy(fileToWrite, buf); err != nil {
		return err
	}

	return os.RemoveAll(k.basePath)
}

func (k *Kpture) MarshalDescription() error {
	bytes, err := json.MarshalIndent(k, "", "    ")
	if err != nil {
		return errors.WithMessage(err, "could not mashal kpture to bytes")
	}

	location := filepath.Join(k.archivePath, k.UUID, "descriptor.json")

	fileToWrite, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR, os.FileMode(fs.ModePerm))
	if err != nil {
		panic(err)
	}

	_, err = fileToWrite.Write(bytes)

	return errors.WithMessage(err, "could not write kpture description")
}
