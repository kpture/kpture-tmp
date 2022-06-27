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

func (k *Kpture) storePackets(basepath string, name string, channel chan gopacket.Packet) error {
	err := os.MkdirAll(basepath, fs.ModePerm)
	if err != nil {
		return errors.WithMessage(err, "could not create directory")
	}

	location := filepath.Join(basepath, name) + ".pcap"

	file, err := os.Create(location)
	if err != nil {
		return errors.WithMessage(err, "error creating file")
	}

	pcapWriter := pcapgo.NewWriter(file)

	err = pcapWriter.WriteFileHeader(pcapFileHeader, layers.LinkTypeEthernet)
	if err != nil {
		return errors.WithMessage(err, "could not write pcap file header")
	}

	go func() {
		for packet := range channel {
			err := pcapWriter.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			if err != nil {
				k.logger.Error(errors.WithMessage(err, "could not write packet"))
			}
		}

		file.Close()
	}()

	return nil
}

func (k *Kpture) createTar() (*bytes.Buffer, error) {
	err := os.MkdirAll(filepath.Join(k.archivePath, k.ProfileName, k.UUID), fs.ModePerm)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create tar directory")
	}

	buf := new(bytes.Buffer)
	tarWriter := tar.NewWriter(buf)

	// walk through every file in the folder
	err = filepath.Walk(k.basePath, func(file string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		var header *tar.Header // generate tar header
		header, err = tar.FileInfoHeader(fileInfo, file)
		if err != nil {
			return errors.WithMessage(err, "could not create tar info header")
		}
		newStr := strings.ReplaceAll(file, filepath.Join(os.TempDir(), k.ProfileName, k.UUID), "")
		newStr = strings.TrimPrefix(newStr, "/")

		header.Name = newStr

		// write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return errors.WithMessage(err, "could not write tar header")
		}
		// if not a dir, write file content
		if !fileInfo.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return errors.WithMessage(err, "error opening pcap file")
			}
			if _, err := io.Copy(tarWriter, data); err != nil {
				return errors.WithMessage(err, "error copying file info")
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "error walking recursively in directory")
	}

	if err := tarWriter.Close(); err != nil {
		return nil, errors.WithMessage(err, "error closing tarWriter")
	}

	return buf, nil
}

func (k *Kpture) writeFile(buf *bytes.Buffer) error {
	location := filepath.Join(k.archivePath, k.ProfileName, k.UUID, k.Name+".tar")

	fileToWrite, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		return errors.WithMessage(err, "error writing file")
	}

	if _, err := io.Copy(fileToWrite, buf); err != nil {
		return errors.WithMessage(err, "error copying pcap buffer to file")
	}

	if err := os.RemoveAll(k.basePath); err != nil {
		return errors.WithMessage(err, "error removing temporary file ")
	}

	return nil
}

func (k *Kpture) MarshalDescription() error {
	bytes, err := json.MarshalIndent(k, "", "    ")
	if err != nil {
		return errors.WithMessage(err, "could not mashal kpture to bytes")
	}

	location := filepath.Join(k.archivePath, k.ProfileName, k.UUID, "descriptor.json")

	fileToWrite, err := os.OpenFile(location, os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = fileToWrite.Write(bytes)

	return errors.WithMessage(err, "could not write kpture description")
}
