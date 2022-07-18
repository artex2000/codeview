package shaper

import (
//        "math"
        "github.com/artex2000/codeview/font"
        "github.com/artex2000/codeview/thirdparty/pixelgl"
)

type Shaper struct {
        Render          *pixelgl.Render
        Font            *font.Monofont
}

func NewShaper(r *pixelgl.Render, f *font.Monofont) *Shaper {
        return &Shaper{Render: r, Font: f}
}

func (s *Shaper) DrawQuad(x, y, w, h float32, c pixelgl.RGBA) {
        quad := []float32{
                      x, y + h,     0, c.R, c.G, c.B, -1, -1,       //left top
                      x, y,         0, c.R, c.G, c.B, -1, -1,       //left bottom
                      x + w, y,     0, c.R, c.G, c.B, -1, -1,       //right bottom
                      x + w, y + h, 0, c.R, c.G, c.B, -1, -1,       //right top
              }
        s.Render.PushQuad(quad)
}

func (s *Shaper) DrawPointRight(x, y, w, h float32, c pixelgl.RGBA) {
        rect := pixelgl.R(float64(x), float64(y), float64(x+w), float64(y+h))
        points := EquilateralTriangle(rect, East)
        triangle := []float32{
                      float32(points[0].X), float32(points[0].Y), 0, c.R, c.G, c.B, -1, -1,
                      float32(points[1].X), float32(points[1].Y), 0, c.R, c.G, c.B, -1, -1,
                      float32(points[2].X), float32(points[2].Y), 0, c.R, c.G, c.B, -1, -1,
                  }
        s.Render.PushTriangle(triangle)
}

func (s *Shaper) DrawPointDown(x, y, w, h float32, c pixelgl.RGBA) {
        rect := pixelgl.R(float64(x), float64(y), float64(x+w), float64(y+h))
        points := EquilateralTriangle(rect, South)
        triangle := []float32{
                      float32(points[0].X), float32(points[0].Y), 0, c.R, c.G, c.B, -1, -1,
                      float32(points[1].X), float32(points[1].Y), 0, c.R, c.G, c.B, -1, -1,
                      float32(points[2].X), float32(points[2].Y), 0, c.R, c.G, c.B, -1, -1,
                  }
        s.Render.PushTriangle(triangle)
}

func (s *Shaper) DrawString(str string, penX, penY float32, c pixelgl.RGBA) {
        for _, t := range str {
                if t == rune(' ') {
                        penX += float32(s.Font.SpaceAdvance)
                        continue
                }

                g, ok := s.Font.Glyphs[rune(t)]
                if !ok {
                        continue
                }
                quad := []float32{
                        penX + float32(g.OffsetX), penY + float32(g.OffsetY),                      0, c.R, c.G, c.B, g.TexS0, g.TexT0,        //left top
                        penX + float32(g.OffsetX), penY + float32(g.OffsetY - g.Height),           0, c.R, c.G, c.B, g.TexS0, g.TexT1,        //left bottom
                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY - g.Height), 0, c.R, c.G, c.B, g.TexS1, g.TexT1,        //right bottom
                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY),            0, c.R, c.G, c.B, g.TexS1, g.TexT0,        //right top
                }
                s.Render.PushQuad(quad)
                penX += float32(g.Advance)
        }
}


