package rec

import (
	"daq/internal/config"
	"encoding/gob"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type MeasurementRecord struct {
	Timestamp     int64
	UploadSpeed   float64
	DownloadSpeed float64
}

type RecordSet struct {
	Records []MeasurementRecord
}

var recordIoMutex = sync.Mutex{}

func (set *RecordSet) Store(cfg *config.PulseConfig) error {
	recordIoMutex.Lock()
	defer recordIoMutex.Unlock()

	file, err := os.OpenFile(cfg.Testing.RecordsPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(set)
}

func (set *RecordSet) Append(record MeasurementRecord, cfg *config.PulseConfig) error {
	set.Records = append(set.Records, record)
	return set.Store(cfg)
}

func (set *RecordSet) Load(cfg *config.PulseConfig) error {
	recordIoMutex.Lock()
	defer recordIoMutex.Unlock()

	file, err := os.Open(cfg.Testing.RecordsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewDecoder(file).Decode(set)
}

func Initialize(cfg *config.PulseConfig) (*RecordSet, error) {
	f, err := os.Stat(cfg.Testing.RecordsPath)

	// The file does not exist - Create and store an empty one
	if os.IsNotExist(err) {
		set := &RecordSet{
			Records: []MeasurementRecord{},
		}
		err = set.Store(cfg)
		if err != nil {
			return nil, err
		}
		logrus.Infof("Created new records file: %s", cfg.Testing.RecordsPath)
		return set, nil
	}

	if err != nil {
		return nil, err
	}

	if f.IsDir() {
		err := os.RemoveAll(cfg.Testing.RecordsPath)
		if err != nil {
			return nil, err
		}
		logrus.Infof("Records file is a directory. Removing it now: %s", cfg.Testing.RecordsPath)
		return Initialize(cfg)
	}

	set := &RecordSet{}
	err = set.Load(cfg)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Loading existing records file: %s", cfg.Testing.RecordsPath)
	logrus.Infof("Existing records file contains %d record(s).", len(set.Records))
	return set, nil
}
