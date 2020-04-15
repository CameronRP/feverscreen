// thermal-recorder - record thermal video footage of warm moving objects
//  Copyright (C) 2018, The Cacophony Project
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"time"

	goconfig "github.com/TheCacophonyProject/go-config"
)

type Config struct {
	SPISpeed    int64
	PowerPin    string
	FrameOutput string
	FFCPeriod   time.Duration
}

func ParseConfig(configFolder string) (*Config, error) {
	configRW, err := goconfig.New(configFolder)
	if err != nil {
		return nil, err
	}

	gpio := goconfig.DefaultGPIO()
	if err := configRW.Unmarshal(goconfig.GPIOKey, &gpio); err != nil {
		return nil, err
	}

	lepton := goconfig.DefaultLepton()
	if err := configRW.Unmarshal(goconfig.LeptonKey, &lepton); err != nil {
		return nil, err
	}

	return &Config{
		SPISpeed:    lepton.SPISpeed,
		PowerPin:    gpio.ThermalCameraPower,
		FrameOutput: lepton.FrameOutput,
		FFCPeriod:   lepton.FFCPeriod,
	}, nil
}
