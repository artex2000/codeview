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

        sol_base03 = color.RGBA{0x00, 0x2B, 0x36, 0xFF} 
        sol_base02 = color.RGBA{0x07, 0x36, 0x42, 0xFF} 
        sol_base01 = color.RGBA{0x58, 0x6E, 0x75, 0xFF} 
        sol_base00 = color.RGBA{0x65, 0x7B, 0x83, 0xFF} 
        sol_base0  = color.RGBA{0x83, 0x94, 0x96, 0xFF} 
        sol_base1  = color.RGBA{0x93, 0xA1, 0xA1, 0xFF} 
        sol_base2  = color.RGBA{0xEE, 0xE8, 0xD5, 0xFF} 
        sol_base3  = color.RGBA{0xFD, 0xF6, 0xE3, 0xFF} 

        sol_yellow  = color.RGBA{0xB5, 0x89, 0x00, 0xFF}
        sol_orange  = color.RGBA{0xCB, 0x4B, 0x16, 0xFF}
        sol_red     = color.RGBA{0xDC, 0x32, 0x2F, 0xFF}
        sol_magenta = color.RGBA{0xD3, 0x36, 0x82, 0xFF}
        sol_violet  = color.RGBA{0x6C, 0x71, 0xC4, 0xFF}
        sol_blue    = color.RGBA{0x26, 0x8B, 0xD2, 0xFF}
        sol_cyan    = color.RGBA{0x2A, 0xA1, 0x98, 0xFF}
        sol_green   = color.RGBA{0x85, 0x99, 0x00, 0xFF}
)

var (
        Black  RGBA
        Red    RGBA
        Green  RGBA
        Blue   RGBA
        Yellow RGBA

        SolBase00 RGBA
        SolBase01 RGBA
        SolBase02 RGBA
        SolBase03 RGBA
        SolBase0  RGBA
        SolBase1  RGBA
        SolBase2  RGBA
        SolBase3  RGBA

        SolYellow  RGBA
        SolOrange  RGBA
        SolRed     RGBA
        SolMagenta RGBA
        SolViolet  RGBA
        SolBlue    RGBA
        SolCyan    RGBA
        SolGreen   RGBA
)

func init() {
        Black  = toRGBA(black)
        Red    = toRGBA(red)
        Green  = toRGBA(green)
        Blue   = toRGBA(blue)
        Yellow = toRGBA(yellow)

        SolBase00 = toRGBA(sol_base00)
        SolBase01 = toRGBA(sol_base01)
        SolBase02 = toRGBA(sol_base02)
        SolBase03 = toRGBA(sol_base03)
        SolBase0  = toRGBA(sol_base0)
        SolBase1  = toRGBA(sol_base1)
        SolBase2  = toRGBA(sol_base2)
        SolBase3  = toRGBA(sol_base3)

        SolYellow  = toRGBA(sol_yellow)
        SolOrange  = toRGBA(sol_orange)
        SolRed     = toRGBA(sol_red)
        SolMagenta = toRGBA(sol_magenta)
        SolViolet  = toRGBA(sol_violet)
        SolBlue    = toRGBA(sol_blue)
        SolCyan    = toRGBA(sol_cyan)
        SolGreen   = toRGBA(sol_green)
}

func toRGBA(c color.RGBA) RGBA {
        return RGBA{R: float32(c.R)/255, G: float32(c.G)/255, B: float32(c.B)/255, A: float32(c.A)/255}
}
