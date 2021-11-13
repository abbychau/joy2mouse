package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/0xcafed00d/joystick"
	"github.com/go-vgo/robotgo"
	"github.com/nsf/termbox-go"
)

var prevZero = false
var prevOne = false

func printAt(x, y int, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

func readJoystick(js joystick.Joystick) {
	jinfo, err := js.Read()
	if err != nil {
		printAt(1, 5, "Error: "+err.Error())
		return
	}
	boostFactor := 1
	printAt(1, 5, "Buttons:")
	zeroed := false
	oneed := false
	for button := 0; button < js.ButtonCount(); button++ {
		if jinfo.Buttons&(1<<uint32(button)) != 0 {
			if button == 0 {
				zeroed = true
			}
			if button == 1 {
				oneed = true
			}
			if button == 4 {
				boostFactor *= 2
			}
			if button == 5 {
				boostFactor *= 2
			}
			printAt(10+button, 5, "X")
		} else {
			printAt(10+button, 5, ".")
		}
	}
	if !prevZero && zeroed {
		robotgo.MouseToggle("down")
	}
	if prevZero && !zeroed {
		robotgo.MouseToggle("up")
	}
	if !prevOne && oneed {
		robotgo.MouseToggle("down", "right")
	}
	if prevOne && !oneed {
		robotgo.MouseToggle("up", "right")
	}
	prevZero = zeroed
	prevOne = oneed
	for button := 0; button < js.ButtonCount(); button++ {
	}
	for axis := 0; axis < js.AxisCount(); axis++ {
		if axis == 5 {
			if jinfo.AxisData[axis] < 0 {
				robotgo.ScrollMouse(1, "left")
			}
			if jinfo.AxisData[axis] > 0 {
				robotgo.ScrollMouse(1, "right")
			}
		}
		if axis == 6 {
			if jinfo.AxisData[axis] < 0 {
				robotgo.ScrollMouse(1, "up")
			}
			if jinfo.AxisData[axis] > 0 {
				robotgo.ScrollMouse(1, "down")
			}
		}
		printAt(1, axis+7, fmt.Sprintf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis]))
	}
	robotgo.MoveRelative(
		(jinfo.AxisData[0]/2000+jinfo.AxisData[4]/20000)*boostFactor,
		(jinfo.AxisData[1]/2000+jinfo.AxisData[3]/20000)*boostFactor,
	)
}
func main() {
	jsid := 0
	if len(os.Args) > 1 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		jsid = i
	}
	js, jserr := joystick.Open(jsid)
	if jserr != nil {
		fmt.Println(jserr)
		return
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 40)

	for doQuit := false; !doQuit; {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Ch == 'q' {
					doQuit = true
				}
			}
			if ev.Type == termbox.EventResize {
				termbox.Flush()
			}
		case <-ticker.C:
			printAt(1, 0, "-- Press 'q' to Exit --")
			printAt(1, 1, fmt.Sprintf("Joystick Name: %s", js.Name()))
			printAt(1, 2, fmt.Sprintf("   Axis Count: %d", js.AxisCount()))
			printAt(1, 3, fmt.Sprintf(" Button Count: %d", js.ButtonCount()))

			readJoystick(js)
			termbox.Flush()
		}
	}
}
