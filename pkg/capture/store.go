package capture

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"newproxy/pkg/logger"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

type store struct {
	BasePath           string
	captureInfos       []*kpture
	descriptorLocation string
	logger             *logrus.Entry
}

func newStore(basePath string) *store {
	return &store{
		BasePath:           basePath,
		captureInfos:       []*kpture{},
		descriptorLocation: filepath.Join(basePath, "descriptor.json"),
		logger:             logger.NewLogger("store"),
	}
}

type AdditionalFile struct {
	Path string
	Data string
}

var wireSharkpreference AdditionalFile = AdditionalFile{
	Path: filepath.Join("profiles", "preferences"),
	Data: `nameres.mac_name: FALSE
nameres.network_name: TRUE
nameres.use_custom_dns_servers: TRUE`,
}

func (s *store) Save(kpture *kpture, additionnalFiles ...AdditionalFile) error {
	additionnalFiles = append(additionnalFiles, wireSharkpreference)

	cptureLoc := filepath.Join(s.BasePath, kpture.UUID)

	s.logger.Debug("Saving capture to ", cptureLoc)
	err := os.MkdirAll(cptureLoc, 0755)
	if err != nil {
		s.logger.Error("Error writing archive:", err)
		return err
	}

	location := filepath.Join(s.BasePath, kpture.UUID, kpture.Name) + ".tar.gz"
	s.logger.Debug("Create tar archive ", kpture.Name+".tar.gz")
	out, err := os.Create(location)
	if err != nil {
		s.logger.Error("Error writing archive:", err)
		return err
	}
	gw := gzip.NewWriter(out)
	tw := tar.NewWriter(gw)

	for i, capture := range kpture.captures {
		kpture.Status.Desciption = fmt.Sprintf("Saving capturing pod %d/%d", i+1, len(kpture.captures))
		buf := capture.GetFileBuf()
		if buf != nil {
			hdr := &tar.Header{
				Name: filepath.Join(kpture.Name, "pods", capture.Pod.Name) + ".pcap",
				Mode: 0600,
				Size: int64(len(buf.Bytes())),
			}

			if err := tw.WriteHeader(hdr); err != nil {
				s.logger.Error(err)
				return err
			}
			if _, err := tw.Write(buf.Bytes()); err != nil {
				s.logger.Error(err)
				return err
			}
		}
	}

	for _, file := range additionnalFiles {
		hdr := &tar.Header{
			Name: filepath.Join(kpture.Name, file.Path),
			Mode: 0600,
			Size: int64(len(file.Data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			s.logger.Error(err)
			return err
		}
		if _, err := tw.Write([]byte(file.Data)); err != nil {
			s.logger.Error(err)
			return err
		}
	}

	kpture.ArchiveLocation = filepath.Join(kpture.UUID, kpture.Name) + ".tar.gz"
	tw.Close()
	gw.Close()
	out.Close()

	file, err := os.Open(location)
	if err != nil {
		s.logger.Error(err)
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		s.logger.Error(err)
		return err
	}
	file.Close()

	kpture.Size = humanize.Bytes(uint64(fi.Size()))
	kpture.CleanBuffer()
	kpture.Status.CaptureState = CaptureStatusReady
	kpture.Status.Desciption = "Capture is ready for download"
	return s.writeDescriptor(kpture)
}

func (s *store) writeDescriptor(kp *kpture) error {
	data, err := ioutil.ReadFile(s.descriptorLocation)
	if err != nil {
		if !os.IsNotExist(err) {
			s.logger.Error("Error reading file descriptor", err)
			return err
		}
	}
	descriptor := make(map[string]*kpture)
	if len(data) > 0 {
		if err := json.Unmarshal(data, &descriptor); err != nil {
			s.logger.Error("Error Unmarshaling data descriptor", err)
			return err
		}
	}
	descriptor[kp.UUID] = kp
	d, err := json.Marshal(descriptor)
	if err != nil {
		s.logger.Error("Error Marshaling data descriptor", err)
		return err
	}
	return ioutil.WriteFile(s.descriptorLocation, d, 0644)
}

func (s *store) GetStoreKpture() (map[string]*kpture, error) {
	descriptor := make(map[string]*kpture)
	data, err := ioutil.ReadFile(s.descriptorLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return descriptor, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &descriptor); err != nil {
		s.logger.Error("Error Unmarshaling data descriptor", err)
		return nil, err
	}

	s.logger.Info("Loaded ", len(descriptor), " captures")
	return descriptor, nil
}
