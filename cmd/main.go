package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	"time"
)

func main() {
	fmt.Println("Hello World")

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Use mcu/GPIO pin 10, corresponds to physical pin 19 on the pi
	blinkLED(10, 10)

	// Use mcu/GPIO pin 19 (used for PWM1), corresponds to physical pin 35 on the pi
	fadeLED(19, 2)

	softPWMpin := rpio.Pin(17)
	softPWMpin.Output()
	softPWMpin.High()
	time.Sleep(200 * time.Millisecond)

	fadeSoftPWM(softPWMpin, 10)

	fadedSoftPWM(softPWMpin)

	fmt.Println("End of the World")
}

func fadedSoftPWM(softPWMpin rpio.Pin) {
	softPwm2 := NewSoftPwm(500, 1, func(state rpio.State) { softPWMpin.Write(state) })
	softPwm2.StartBalancedStream(0.50)
	time.Sleep(5 * 1000 * time.Millisecond)
	softPwm2.Stop()
	softPWMpin.Low()
}

func fadeSoftPWM(softPWMpin rpio.Pin, amountSteps int) {
	softPwm := NewSoftPwm(1000, 1, func(state rpio.State) { softPWMpin.Write(state) })
	softPwm.StartBalancedStream(0.0)
	for i := 0; i <= amountSteps; i++ {
		percent := float64(i) / float64(amountSteps)

		softPwm.Duty(percent)
		time.Sleep(500 * time.Millisecond)
	}
	softPwm.Stop()
	softPWMpin.Low()
}

func fadeLED(pin rpio.Pin, amount int) {
	pin.Mode(rpio.Pwm)
	pin.Freq(64000)
	pin.DutyCycleWithPwmMode(0, 32, rpio.Balanced)
	// the LED will be blinking at 2000Hz
	// (source frequency divided by cycle length => 64000/32 = 2000)

	// five times smoothly fade in and out
	for i := 0; i < amount; i++ {
		for i := uint32(0); i < 32; i++ { // increasing brightness
			pin.DutyCycleWithPwmMode(i, 32, rpio.Balanced)
			time.Sleep(time.Second / 32)
		}
		for i := uint32(32); i > 0; i-- { // decreasing brightness
			pin.DutyCycleWithPwmMode(i, 32, rpio.Balanced)
			time.Sleep(time.Second / 32)
		}
	}

	pin.DutyCycleWithPwmMode(0, 32, rpio.Balanced)
}

func blinkLED(pin rpio.Pin, amount int) {

	// Set pin to output mode
	pin.Output()

	// Toggle pin amount of times
	for x := 0; x < amount; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 10)
	}
}
