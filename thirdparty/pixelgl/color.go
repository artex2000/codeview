package pixelgl

import (
        "image/color"
)

type RGBA struct {
        R, G, B, A float32
}

var (
        black  = color.RGBA{0x00, 0x00, 0x00, 0xFF}
        red    = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
        green  = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
        blue   = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
        yellow = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
)

var (
        Black  RGBA
        Red    RGBA
        Green  RGBA
        Blue   RGBA
        Yellow RGBA
)

func init() {
        Black  = toRGBA(black)
        Red    = toRGBA(red)
        Green  = toRGBA(green)
        Blue   = toRGBA(blue)
        Yellow = toRGBA(yellow)
}

func toRGBA(c color.RGBA) RGBA {
        return RGBA{R: float32(c.R)/255, G: float32(c.G)/255, B: float32(c.B)/255, A: float32(c.A)/255}
}
