package main

import (
        "fmt"
        //"os"
        //"image"
        //"image/png"
        //"image/draw"
        "math"
        "math/rand"
        "time"

	"github.com/go-gl/mathgl/mgl32"
        "github.com/artex2000/codeview/thirdparty/pixelgl"
        "github.com/artex2000/codeview/thirdparty/glhf"
        "github.com/artex2000/codeview/font"
       )

var fontFile string = "C:/Windows/Fonts/VictorMono-Regular.ttf"

func main () {
        pixelgl.Run(run)
}

func run() {
        cfg := pixelgl.WindowConfig{
                Title:     "Codeview",
                Bounds:    pixelgl.R(0, 0, 1920, 1080),
                Resizable: true,
                VSync:     true,
        }

        font, err := font.InitFontFromFile(fontFile, 18)
        if err != nil {
                panic (err)
        }

//        fmt.Printf("Asc %d, Dsc %d, Lg %d, Adv %d\n", font.Ascender, font.Descender, font.Linegap, font.SpaceAdvance)

        /*
        //Load PNG for testing
        pic, err := os.Open("./assets/Squares.png")
        if err != nil {
                panic(err)
        }
        img, err := png.Decode(pic)
        if err != nil {
                panic(err)
        }
        rgba := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
        draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
        */


        win, err := pixelgl.NewWindow(cfg)
        if err != nil {
                panic (err)
        }

        Uniforms := pixelgl.VariableList{
                        { Name: "Model",        Type: glhf.Mat4 },
                        { Name: "View",         Type: glhf.Mat4 },
                        { Name: "Projection",   Type: glhf.Mat4 },
                    }

        Attributes := pixelgl.VariableList{
                        { Name: "VertexPosition", Type: glhf.Vec3 },
                        { Name: "VertexColor",    Type: glhf.Vec3 },
                        { Name: "TextureCoord",   Type: glhf.Vec2 },
                    }

        r := pixelgl.NewRender(VertexShader, FragmentShader, Uniforms, Attributes, font.Atlas)

        //Camera setup
        eye    := mgl32.Vec3{ 0.0, 0.0, 2.0 }
        center := mgl32.Vec3{ 0.0, 0.0, 0.0 }
        up     := mgl32.Vec3{ 0.0, 1.0, 0.0 }

        //Projection setup (will match viewport width/height)
        w, h := float32(win.Bounds().W()), float32(win.Bounds().H())
        iw, ih := int(math.Floor(win.Bounds().W())), int(math.Floor(win.Bounds().H()))

        //Transform setup
        r.Model      = mgl32.Ident4()
        r.View       = mgl32.LookAtV(eye, center, up)
        r.Projection = mgl32.Ortho(0, w, 0, h, 0.1, 5)

        r.SetTransform(true, true, true)
        //r.SetTexture("Texture")

        border := 2     //apron size in pixels
        lineSpace := font.Ascender + font.Descender + font.Linegap
        //Calculate how many symbols we can have on the screen with 5px bounds
        cols := (iw - 2 * border) / font.SpaceAdvance
        rows := (ih - 2 * border) / lineSpace

        win.SetCanvas(r)

        frameDt := int64(0)
        var hello []string
        freeze := false

        closed := true

        for !win.Closed() {
                start := time.Now()

                //Push glyph quads
                //hello := "The quick brown fox jumps over the lazy dog. 1234567890"
                extras := 10
                if len(hello) == 0 || !freeze {
                        hello = generateGlyphs(rows + extras, cols + extras)
                }

                penX, penY := float32(border), float32(ih - border - font.Ascender)
                for _, s := range hello {
                        drawString(r, font, s, penX, penY, pixelgl.Black)
                        penX = float32(border)
                        penY -= float32(lineSpace)
                }

                //draw frame duration in top right corner
                if frameDt > 0 {
                        s := fmt.Sprintf("%dms", frameDt)
                        fW := float32(len(s) * font.SpaceAdvance)
                        fH := float32(font.Ascender)
                        //Put it in the right top corner
                        fX := w - fW
                        fY := h - fH
                        drawQuad(r, fX, fY, fW, fH, pixelgl.Yellow)
                        drawString(r, font, s, fX, fY, pixelgl.Red)
                }

                //draw directory triangle thing in top left corner
                qW, qH := float32(lineSpace - 2), float32(lineSpace - 2)
                qX, qY := float32(2), h - qH - 2
                drawQuad(r, qX, qY, qW, qH, pixelgl.Yellow)
                if closed {
                        drawPointRight(r, qX, qY, qW, qH, pixelgl.Blue) 
                } else {
                        drawPointDown(r, qX, qY, qW, qH, pixelgl.Blue) 
                }



                r.SetVertices()
                win.Update()
                r.ResetVertices()

                //process events
                for _, e := range win.Events() {
                        switch e.Type {
                        case pixelgl.KeyPress:
                                if e.Key == pixelgl.KeySpace {
                                        freeze = !freeze
                                        closed = !closed
                                } else if e.Key == pixelgl.KeyEscape {
                                        win.SetClosed(true)
                                }
                        }
                }

                elapsed := time.Since(start)
                frameDt = elapsed.Milliseconds()
        }
}

func drawString(r *pixelgl.Render, font *font.Monofont, s string, penX, penY float32, c pixelgl.RGBA) {
        for _, t := range s {
                if t == rune(' ') {
                        penX += float32(font.SpaceAdvance)
                        continue
                }

                g, ok := font.Glyphs[rune(t)]
                if !ok {
                        continue
                }
                quad := []float32{
                        penX + float32(g.OffsetX), penY + float32(g.OffsetY),                      0, c.R, c.G, c.B, g.TexS0, g.TexT0,        //left top
                        penX + float32(g.OffsetX), penY + float32(g.OffsetY - g.Height),           0, c.R, c.G, c.B, g.TexS0, g.TexT1,        //left bottom
                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY - g.Height), 0, c.R, c.G, c.B, g.TexS1, g.TexT1,        //right bottom
                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY),            0, c.R, c.G, c.B, g.TexS1, g.TexT0,        //right top
                }
                r.PushQuad(quad)
                penX += float32(g.Advance)
        }
}

func drawQuad(r *pixelgl.Render, x, y, w, h float32, c pixelgl.RGBA) {
        quad := []float32{
                      x, y + h,     0, c.R, c.G, c.B, -1, -1,       //left top
                      x, y,         0, c.R, c.G, c.B, -1, -1,       //left bottom
                      x + w, y,     0, c.R, c.G, c.B, -1, -1,       //right bottom
                      x + w, y + h, 0, c.R, c.G, c.B, -1, -1,       //right top
              }
        r.PushQuad(quad)
}

func drawPointRight(r *pixelgl.Render, x, y, w, h float32, c pixelgl.RGBA) {
        triangle := []float32{
                      x, y + h,         0, c.R, c.G, c.B, -1, -1,   //left top
                      x, y,             0, c.R, c.G, c.B, -1, -1,   //left bottom
                      x + w, y + h / 2, 0, c.R, c.G, c.B, -1, -1,   //right middle
                  }
        r.PushTriangle(triangle)
}

func drawPointDown(r *pixelgl.Render, x, y, w, h float32, c pixelgl.RGBA) {
        triangle := []float32{
                      x, y + h,     0, c.R, c.G, c.B, -1, -1,   //left top
                      x + w / 2, y, 0, c.R, c.G, c.B, -1, -1,   //middle bottom
                      x + w, y + h, 0, c.R, c.G, c.B, -1, -1,   //right top
                  }
        r.PushTriangle(triangle)
}

func generateGlyphs(rows, cols int) []string {
        out := make([]string, rows)
        for i, _ := range out {
                s := make([]byte, cols)
                for j, _ := range s {
                        s[j] = byte(rand.Intn(0x7f - 0x20) + 0x20)
                }
                out[i] = string(s)
        }
        return out
}
