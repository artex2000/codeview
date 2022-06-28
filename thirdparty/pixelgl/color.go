package pixelgl

type RGBA struct {
        R, G, B, A float32
}

var (
        Black  = RGBA{0, 0, 0, 1}
        Red    = RGBA{1, 0, 0, 1}
        Yellow = RGBA{0, 1, 1, 1}
)
