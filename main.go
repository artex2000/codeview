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
                Bounds:    pixelgl.R(10, 10, 3700, 2100),
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

        border := 5     //apron size in pixels
        lineSpace := font.Ascender + font.Descender + font.Linegap
        //Calculate how many symbols we can have on the screen with 5px bounds
        cols := (iw - 2 * border) / font.SpaceAdvance
        rows := (ih - 2 * border) / lineSpace


        win.SetCanvas(r)

        for !win.Closed() {
                start := time.Now()

                //Push glyph quads
                //hello := "The quick brown fox jumps over the lazy dog. 1234567890"
                extras := 10
                hello := generateGlyphs(rows + extras, cols + extras)

                penX, penY := float32(border), float32(ih - border)
                for _, s := range hello {
                        for _, c := range s {
                                if c == rune(' ') {
                                        penX += float32(font.SpaceAdvance)
                                        continue
                                }

                                g, ok := font.Glyphs[rune(c)]
                                if !ok {
                                        continue
                                }
                                quad := []float32{
                                        penX + float32(g.OffsetX), penY + float32(g.OffsetY), 0, 0, 0, 0, g.TexS0, g.TexT0,        //left top
                                        penX + float32(g.OffsetX), penY + float32(g.OffsetY - g.Height), 0, 0, 0, 0, g.TexS0, g.TexT1,        //left bottom
                                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY - g.Height), 0, 0, 0, 0, g.TexS1, g.TexT1,        //right bottom
                                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY), 0, 0, 0, 0, g.TexS1, g.TexT0,        //right top
                                }
                                r.PushQuad(quad)
                                penX += float32(g.Advance)
                        }
                        penX = float32(border)
                        penY -= float32(lineSpace)
                }

                r.SetVertices()
                win.Update()
                r.ResetVertices()

                elapsed := time.Since(start)
                fmt.Printf("%d ms\n", elapsed.Milliseconds())
        }
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
