
## wow
[![GoDoc](https://godoc.org/github.com/gernest/wow?status.svg)](https://godoc.org/github.com/gernest/wow)

Beautiful spinners for Go commandline apps

![wow](static/wow.gif)

## Install
    go get -u github.com/gernest/wow

## Usage

```go
package main

import (
	"os"
	"time"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

func main() {
	w := wow.New(os.Stdout, spin.Get(spin.Dots), "Such Spins")
	w.Start()
	time.Sleep(2 * time.Second)
	w.Text("Very emojis").Spinner(spin.Get(spin.Hearts))
	time.Sleep(2 * time.Second)
	w.PersistWith(spin.Spinner{Frames: []string{"👍"}}, " Wow!")
}
```

## Available spinners

  Name  | What it looks like 
--------|---------------------
 `Star2` | ![star2](static/star2.gif)
 `GrowHorizontal` | ![growHorizontal](static/growHorizontal.gif)
 `Squish` | ![squish](static/squish.gif)
 `Toggle12` | ![toggle12](static/toggle12.gif)
 `Smiley` | ![smiley](static/smiley.gif)
 `Hearts` | ![hearts](static/hearts.gif)
 `Dots3` | ![dots3](static/dots3.gif)
 `Dots11` | ![dots11](static/dots11.gif)
 `Balloon` | ![balloon](static/balloon.gif)
 `Clock` | ![clock](static/clock.gif)
 `Dots4` | ![dots4](static/dots4.gif)
 `SquareCorners` | ![squareCorners](static/squareCorners.gif)
 `CircleHalves` | ![circleHalves](static/circleHalves.gif)
 `Star` | ![star](static/star.gif)
 `Arc` | ![arc](static/arc.gif)
 `Toggle13` | ![toggle13](static/toggle13.gif)
 `BoxBounce` | ![boxBounce](static/boxBounce.gif)
 `Line2` | ![line2](static/line2.gif)
 `Pipe` | ![pipe](static/pipe.gif)
 `Triangle` | ![triangle](static/triangle.gif)
 `Shark` | ![shark](static/shark.gif)
 `Line` | ![line](static/line.gif)
 `Arrow` | ![arrow](static/arrow.gif)
 `Earth` | ![earth](static/earth.gif)
 `Dots5` | ![dots5](static/dots5.gif)
 `Toggle11` | ![toggle11](static/toggle11.gif)
 `CircleQuarters` | ![circleQuarters](static/circleQuarters.gif)
 `Toggle9` | ![toggle9](static/toggle9.gif)
 `Dots9` | ![dots9](static/dots9.gif)
 `Bounce` | ![bounce](static/bounce.gif)
 `Toggle2` | ![toggle2](static/toggle2.gif)
 `Toggle7` | ![toggle7](static/toggle7.gif)
 `Arrow3` | ![arrow3](static/arrow3.gif)
 `Moon` | ![moon](static/moon.gif)
 `Dots6` | ![dots6](static/dots6.gif)
 `Christmas` | ![christmas](static/christmas.gif)
 `Dots10` | ![dots10](static/dots10.gif)
 `Hamburger` | ![hamburger](static/hamburger.gif)
 `BoxBounce2` | ![boxBounce2](static/boxBounce2.gif)
 `BouncingBar` | ![bouncingBar](static/bouncingBar.gif)
 `Flip` | ![flip](static/flip.gif)
 `Dots8` | ![dots8](static/dots8.gif)
 `Dots12` | ![dots12](static/dots12.gif)
 `Noise` | ![noise](static/noise.gif)
 `Toggle3` | ![toggle3](static/toggle3.gif)
 `Toggle6` | ![toggle6](static/toggle6.gif)
 `Runner` | ![runner](static/runner.gif)
 `Dqpb` | ![dqpb](static/dqpb.gif)
 `Dots` | ![dots](static/dots.gif)
 `Toggle4` | ![toggle4](static/toggle4.gif)
 `Monkey` | ![monkey](static/monkey.gif)
 `Dots7` | ![dots7](static/dots7.gif)
 `SimpleDots` | ![simpleDots](static/simpleDots.gif)
 `GrowVertical` | ![growVertical](static/growVertical.gif)
 `Circle` | ![circle](static/circle.gif)
 `Toggle` | ![toggle](static/toggle.gif)
 `Toggle5` | ![toggle5](static/toggle5.gif)
 `Arrow2` | ![arrow2](static/arrow2.gif)
 `Dots2` | ![dots2](static/dots2.gif)
 `Toggle8` | ![toggle8](static/toggle8.gif)
 `Toggle10` | ![toggle10](static/toggle10.gif)
 `BouncingBall` | ![bouncingBall](static/bouncingBall.gif)
 `SimpleDotsScrolling` | ![simpleDotsScrolling](static/simpleDotsScrolling.gif)
 `Pong` | ![pong](static/pong.gif)
 `Weather` | ![weather](static/weather.gif)
 `Balloon2` | ![balloon2](static/balloon2.gif)
