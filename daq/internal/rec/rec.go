package rec

import (
	"encoding/gob"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const RecordPath = "records.gob"

type MeasurementRecord struct {
	Timestamp     int64
	UploadSpeed   float64
	DownloadSpeed float64
}

type RecordSet struct {
	Records []MeasurementRecord
}

var recordIoMutex = sync.Mutex{}

func (set *RecordSet) Store() error {
	recordIoMutex.Lock()
	defer recordIoMutex.Unlock()

	file, err := os.OpenFile(RecordPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(set)
}

func (set *RecordSet) Append(record MeasurementRecord) error {
	set.Records = append(set.Records, record)
	return set.Store()
}

func (set *RecordSet) Load() error {
	recordIoMutex.Lock()
	defer recordIoMutex.Unlock()

	file, err := os.Open(RecordPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewDecoder(file).Decode(set)
}

func Initialize() (*RecordSet, error) {
	f, err := os.Stat(RecordPath)

	// The file does not exist - Create and store an empty one
	if os.IsNotExist(err) {
		set := &RecordSet{
			Records: []MeasurementRecord{},
		}
		err = set.Store()
		if err != nil {
			return nil, err
		}
		logrus.Infof("Created new records file: %s", RecordPath)
		return set, nil
	}

	if err != nil {
		return nil, err
	}

	if f.IsDir() {
		err := os.RemoveAll(RecordPath)
		if err != nil {
			return nil, err
		}
		logrus.Infof("Records file is a directory. Removing it now: %s", RecordPath)
		return Initialize()
	}

	set := &RecordSet{}
	err = set.Load()
	if err != nil {
		return nil, err
	}
	logrus.Infof("Loading existing records file: %s", RecordPath)
	logrus.Infof("Existing records file contains %d record(s).", len(set.Records))
	return set, nil
}
