package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllDefaults(t *testing.T) {
	conf, err := ParseConfig([]byte(""), []byte(""))
	require.NoError(t, err)
	require.NoError(t, conf.Validate())

	assert.Equal(t, Config{
		DeviceName:   "",
		FrameInput:   "/var/run/lepton-frames",
		OutputDir:    "/var/spool/cptv",
		MinSecs:      10,
		MaxSecs:      600,
		PreviewSecs:  3,
		MinDiskSpace: 200,
		Motion: MotionConfig{
			TempThresh:        2900,
			DeltaThresh:       50,
			CountThresh:       3,
			NonzeroMaxPercent: 50,
			FrameCompareGap:   45,
			UseOneDiffOnly:    true,
			Verbose:           false,
			TriggerFrames:     2,
			WarmerOnly:        true,
		},
		Turret: TurretConfig{
			Active: false,
			PID:    []float64{0.05, 0, 0},
			ServoX: ServoConfig{
				Active:   false,
				Pin:      "17",
				MaxAng:   160,
				MinAng:   20,
				StartAng: 90,
			},
			ServoY: ServoConfig{
				Active:   false,
				Pin:      "18",
				MaxAng:   160,
				MinAng:   20,
				StartAng: 90,
			},
		},
	}, *conf)
}

func TestAllSet(t *testing.T) {
	// All config set at non-default values.
	config := []byte(`
frame-input: "/some/sock"
output-dir: "/some/where"
min-secs: 2
max-secs: 10
preview-secs: 5
window-start: 17:10
window-end: 07:20
min-disk-space: 321
motion:
    temp-thresh: 2000
    delta-thresh: 20
    count-thresh: 1
    nonzero-max-percent: 20
    frame-compare-gap: 90
    one-diff-only: false
    trigger-frames: 1
    verbose: true
    warmer-only: false
leds:
    recording: "RecordingPIN"
    running: "RunningPIN"
turret:
    active: true
    pid:
      - 1
      - 2
      - 3
    servo-x:
      active: false
      pin: "pin"
      min-ang: 0
      max-ang: 180
      start-ang: 30
    servo-y:
      active: true
      pin: "pin1"
      min-ang: 10
      max-ang: 190
      start-ang: 40
`)

	uploaderConfig := []byte(`
device-name: "aDeviceName"
`)

	conf, err := ParseConfig(config, uploaderConfig)
	require.NoError(t, err)
	require.NoError(t, conf.Validate())

	assert.Equal(t, Config{
		DeviceName:   "aDeviceName",
		FrameInput:   "/some/sock",
		OutputDir:    "/some/where",
		MinSecs:      2,
		MaxSecs:      10,
		PreviewSecs:  5,
		WindowStart:  time.Date(0, 1, 1, 17, 10, 0, 0, time.UTC),
		WindowEnd:    time.Date(0, 1, 1, 07, 20, 0, 0, time.UTC),
		MinDiskSpace: 321,
		Motion: MotionConfig{
			TempThresh:        2000,
			DeltaThresh:       20,
			CountThresh:       1,
			NonzeroMaxPercent: 20,
			FrameCompareGap:   90,
			UseOneDiffOnly:    false,
			Verbose:           true,
			TriggerFrames:     1,
			WarmerOnly:        false,
		},
		Turret: TurretConfig{
			Active: true,
			PID:    []float64{1, 2, 3},
			ServoX: ServoConfig{
				Active:   false,
				Pin:      "pin",
				MaxAng:   180,
				MinAng:   0,
				StartAng: 30,
			},
			ServoY: ServoConfig{
				Active:   true,
				Pin:      "pin1",
				MaxAng:   190,
				MinAng:   10,
				StartAng: 40,
			},
		},
	}, *conf)
}

func TestInvalidWindowStart(t *testing.T) {
	conf, err := ParseConfig([]byte("window-start: 25:10"), []byte(""))
	assert.Nil(t, conf)
	assert.EqualError(t, err, "invalid window-start")
}

func TestInvalidWindowEnd(t *testing.T) {
	conf, err := ParseConfig([]byte("window-end: 25:10"), []byte(""))
	assert.Nil(t, conf)
	assert.EqualError(t, err, "invalid window-end")
}

func TestWindowEndWithoutStart(t *testing.T) {
	conf, err := ParseConfig([]byte("window-end: 09:10"), []byte(""))
	assert.Nil(t, conf)
	assert.EqualError(t, err, "window-end is set but window-start isn't")
}

func TestWindowStartWithoutEnd(t *testing.T) {
	conf, err := ParseConfig([]byte("window-start: 09:10"), []byte(""))
	assert.Nil(t, conf)
	assert.EqualError(t, err, "window-start is set but window-end isn't")
}