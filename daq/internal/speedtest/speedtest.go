package speedtest

import (
	"context"
	"daq/internal/config"
	"daq/internal/rec"
	"encoding/json"
	"os/exec"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Result struct {
	Download struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"download"`

	Upload struct {
		Bandwidth int64 `json:"bandwidth"`
	} `json:"upload"`
}

func Run(ctx context.Context) (*Result, error) {
	logrus.Infof("Running speedtest...")
	cmd := exec.CommandContext(ctx, "speedtest", "-f", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var result Result
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (res *Result) Speeds() (downMbps float64, upMbps float64) {
	return float64(res.Download.Bandwidth*8.0) / 1e6, float64(res.Upload.Bandwidth*8.0) / 1e6
}

func StartSpeedTestService(records *rec.RecordSet, cfg *config.PulseConfig, ctx context.Context, wg *sync.WaitGroup) {
	logrus.Infof("Speedtest service started")
	logrus.Infof("Speedtest will run every %d seconds", cfg.TestInterval)
	last := time.Now().Unix()
	isRunning := true
	go func() {
		wg.Add(1)
		defer wg.Done()
		for isRunning {
			select {
			case <-ctx.Done():
				isRunning = false
				continue
			default:
				now := time.Now().Unix()
				elapsed := now - last
				if elapsed < cfg.TestInterval {
					time.Sleep(min(time.Duration(cfg.TestInterval-(now-last))*time.Second, 500*time.Millisecond))
					continue
				}
				res, err := Run(ctx)
				last = time.Now().Unix()
				if err != nil {
					continue
				}
				down, up := res.Speeds()
				err = records.Append(rec.MeasurementRecord{
					Timestamp:     now,
					UploadSpeed:   up,
					DownloadSpeed: down,
				})
				if err != nil {
					logrus.Error(err)
					continue
				}
				logrus.Infof("Upload %.2f Mbps, Download %.2f Mbps", up, down)
			}
		}
		logrus.Infof("Speedtest service stopped")
	}()
}
