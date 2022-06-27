package main

import (
        "fmt"
        //"os"
        //"image"
        //"image/png"
        //"image/draw"
        "image/color"

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
                Bounds:    pixelgl.R(100, 100, 1600, 1600),
                Resizable: true,
                VSync:     true,
        }

        font, err := font.InitFontFromFile(fontFile, 32)
        if err != nil {
                panic (err)
        }

        fmt.Printf("Asc %d, Dsc %d, Lg %d, Adv %d\n", font.Ascender, font.Descender, font.Linegap, font.SpaceAdvance)

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
                        { Name: "Foreground",   Type: glhf.Vec4 },
                        { Name: "Background",   Type: glhf.Vec4 },
                    }

        Attributes := pixelgl.VariableList{
                        { Name: "VertexPosition", Type: glhf.Vec3 },
                        { Name: "TextureCoord",   Type: glhf.Vec2 },
                    }

        r := pixelgl.NewRender(VertexShader, FragmentShader, Uniforms, Attributes, font.Atlas)

        //Camera setup
        eye    := mgl32.Vec3{ 0.0, 0.0, 2.0 }
        center := mgl32.Vec3{ 0.0, 0.0, 0.0 }
        up     := mgl32.Vec3{ 0.0, 1.0, 0.0 }

        //Projection setup (will match viewport width/height)
        w, h := float32(win.Bounds().W()), float32(win.Bounds().H())

        //Transform setup
        r.Model      = mgl32.Ident4()
        r.View       = mgl32.LookAtV(eye, center, up)
        r.Projection = mgl32.Ortho(0, w, 0, h, 0.1, 5)

        r.Foreground = color.RGBA {R: 0xff, G: 0xff, B: 0xff, A: 0x1}
        r.Background = color.RGBA {R: 0x0, G: 0x0, B: 0x0, A: 0x1}
        r.SetTransform(true, true, true)
        r.SetColors()
        r.SetTexture("Texture")

        //Push glyph quads
        hello := "Hello Master"

        //Find starting point
        midY := (h - float32(font.Ascender + font.Descender + font.Linegap)) / 2
        midX := (w - float32(len(hello) * font.SpaceAdvance)) / 2

        penX, penY := midX, midY
        for _, s := range hello {
                if s == rune(' ') {
                        penX += float32(font.SpaceAdvance)
                        continue
                }

                g, ok := font.Glyphs[rune(s)]
                if !ok {
                        continue
                }
                quad := []float32{
                        penX + float32(g.OffsetX), penY + float32(g.OffsetY), 0, g.TexS0, g.TexT0,        //left top
                        penX + float32(g.OffsetX), penY + float32(g.OffsetY - g.Height), 0, g.TexS0, g.TexT1,        //left bottom
                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY - g.Height), 0, g.TexS1, g.TexT1,        //right bottom
                        penX + float32(g.OffsetX + g.Width), penY + float32(g.OffsetY), 0, g.TexS1, g.TexT0,        //right top
                      }
                r.PushQuad(quad)
                penX += float32(g.Advance)
        }
        /*
        quad := []float32{
                100, 800, 0, 0, 0,
                100, 800 - 256, 0, 0, 1,
                356, 800 - 256, 0, 1, 1,
                356, 800, 0, 1, 0,
        }
        r.PushQuad(quad)
        */

        r.SetVertices()

        win.SetCanvas(r)

        for !win.Closed() {
//                win.Clear()
                win.Update()
        }
}
