package main

import (
	"github.com/stianeikeland/go-rpio/v4"
	"go.uber.org/atomic"
	"time"
)

type SoftPWM struct {
	callback func(edge rpio.State)

	frequency  int32
	resolution int32

	ticker  *time.Ticker
	pwmStop chan bool
	state   rpio.State

	balancedDuty          atomic.Float32
	balancedBitPercentage float32
	balancedAccumulator   float32
}

func NewSoftPwm(frequency int32, resolution uint8, callback func(state rpio.State)) *SoftPWM {
	return &SoftPWM{
		callback:              callback,
		frequency:             frequency,
		resolution:            int32(resolution),
		ticker:                time.NewTicker(calculateBitDuration(frequency, int32(resolution))),
		pwmStop:               make(chan bool),
		state:                 rpio.Low,
		balancedBitPercentage: 1.0 / float32(resolution),
		balancedAccumulator:   0.0,
	}
}

func calculateBitDuration(frequency int32, resolution int32) time.Duration {
	return time.Duration(int64(time.Second) / int64(frequency*resolution))
}

func (spwm *SoftPWM) StartBalancedStream(dutyPercentage float64) {
	spwm.state = rpio.Low
	spwm.Duty(dutyPercentage)

	go func() {
		rateAccumulator := float32(0.0)

		for {
			select {
			case <-spwm.pwmStop:
				return
			case <-spwm.ticker.C:
				rate := 1.0 / spwm.balancedDuty.Load()
				rateAccumulator += 1.0

				oldState := spwm.state
				var newState rpio.State

				if rateAccumulator >= rate {
					rateAccumulator -= rate
					newState = rpio.High
				} else {
					newState = rpio.Low
				}

				if oldState != newState {
					spwm.state = newState
					spwm.callback(newState)
				}

			}
		}
	}()
}

func (spwm *SoftPWM) StartMarkSpace(dutyPercentage float64) {
	spwm.state = rpio.Low
	spwm.Duty(dutyPercentage)

	go func() {
		for {
			select {
			case <-spwm.pwmStop:
				return
			case <-spwm.ticker.C:
				spwm.balancedAccumulator += spwm.balancedBitPercentage

				duty := spwm.balancedDuty.Load()

				if spwm.balancedAccumulator >= 1.0 {
					spwm.balancedAccumulator = 1.0 - spwm.balancedAccumulator

					if spwm.state == rpio.Low {
						spwm.state = rpio.High
						spwm.callback(spwm.state)
					}

				} else if spwm.balancedAccumulator >= duty {
					if spwm.state == rpio.High {
						spwm.balancedAccumulator -= spwm.balancedAccumulator - duty
						spwm.state = rpio.Low
						spwm.callback(spwm.state)
					}
				}

			}
		}
	}()
}

func (spwm *SoftPWM) Stop() {
	spwm.pwmStop <- true
	spwm.ticker.Stop()
	spwm.state = rpio.Low
}

func (spwm *SoftPWM) Duty(percent float64) {
	if percent < 0.0 {
		percent = 0.0
	} else if percent > 1.0 {
		percent = 1.0
	}

	spwm.balancedDuty.Store(float32(percent))
}
