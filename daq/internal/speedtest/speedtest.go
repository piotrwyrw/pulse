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

func invokeCommand(ctx context.Context, cfg *config.PulseConfig, binary string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, binary, arg...)
	logrus.Debugf("Running: %s", cmd.String())
	return cmd
}

func run(ctx context.Context, cfg *config.PulseConfig) (*Result, error) {
	logrus.Infof("Running speedtest...")
	cmd := invokeCommand(ctx, cfg, cfg.Testing.BinaryPath, "-f", "json")
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

func AcceptLicenses(ctx context.Context, cfg *config.PulseConfig) error {
	err := invokeCommand(ctx, cfg, cfg.Testing.BinaryPath, "--accept-license", "--accept-gdpr").Run()
	if err != nil {
		return err
	}
	return nil
}

func (res *Result) Speeds() (downMbps float64, upMbps float64) {
	return float64(res.Download.Bandwidth*8.0) / 1e6, float64(res.Upload.Bandwidth*8.0) / 1e6
}

func StartSpeedTestService(records *rec.RecordSet, cfg *config.PulseConfig, ctx context.Context, wg *sync.WaitGroup) error {
	logrus.Infof("Speedtest service started")
	err := AcceptLicenses(ctx, cfg)
	if err != nil {
		logrus.Errorf("Could not accept license: %v", err)
		return err
	}
	logrus.Infof("Speedtest will run every %d seconds", cfg.Testing.Interval)
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
				// Keep waiting until the set testing interval, though not longer than 500ms to not impede
				// the main loop
				now := time.Now().Unix()
				elapsed := now - last
				if elapsed < cfg.Testing.Interval {
					time.Sleep(min(time.Duration(cfg.Testing.Interval-(now-last))*time.Second, 500*time.Millisecond))
					continue
				}

				// run the speed test
				res, err := run(ctx, cfg)
				last = time.Now().Unix()
				var down, up float64
				if err != nil {
					logrus.Infof("Could not run speedtest: %v", err)
					down, up = 0, 0
					continue
				} else {
					down, up = res.Speeds()
				}
				err = records.Append(rec.MeasurementRecord{
					Timestamp:     now,
					UploadSpeed:   up,
					DownloadSpeed: down,
				}, cfg)
				if err != nil {
					logrus.Error(err)
					continue
				}
				logrus.Infof("Upload %.2f Mbps, Download %.2f Mbps", up, down)
			}
		}
		logrus.Infof("Speedtest service stopped")
	}()

	return nil
}
