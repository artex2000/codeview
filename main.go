package main

import (
//        "fmt"
        //"os"
        //"image"
        //"image/png"
        //"image/draw"
        "time"

	"github.com/go-gl/mathgl/mgl32"
        "github.com/artex2000/codeview/thirdparty/pixelgl"
        "github.com/artex2000/codeview/thirdparty/glhf"
        "github.com/artex2000/codeview/font"
        "github.com/artex2000/codeview/shaper"
       )

var fontFile string = "C:/Windows/Fonts/VictorMono-Regular.ttf"

func main () {
        pixelgl.Run(run)
}

func run() {
        cfg := pixelgl.WindowConfig{
                Title:     "Codeview",
                Bounds:     pixelgl.R(0, 0, 1920, 1080),
                Resizable:  true,
                VSync:      true,
                ClearColor: pixelgl.SolBase03,
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

        render := pixelgl.NewRender(VertexShader, FragmentShader, Uniforms, Attributes, font.Atlas)
        shaper := shaper.NewShaper(render, font)

        //Projection setup (will match viewport width/height)
        winRect := win.Bounds()
        w, h := float32(winRect.W()), float32(winRect.H())

        //Camera setup
        eye    := mgl32.Vec3{ 0.0, 0.0, 2.0 }
        center := mgl32.Vec3{ 0.0, 0.0, 0.0 }
        up     := mgl32.Vec3{ 0.0, 1.0, 0.0 }

        //Transform setup
        render.Model      = mgl32.Ident4()
        render.View       = mgl32.LookAtV(eye, center, up)
        render.Projection = mgl32.Ortho(0, w, 0, h, 0.1, 5)

        render.SetTransform(true, true, true)

        win.SetCanvas(render)

        container := NewContainer(shaper, winRect)
        container.AddTree("C:/go/work/src/github.com/artex2000/codeview")

        frameDt := int64(0)

        for !win.Closed() {
                start := time.Now()

                container.Draw()

                render.SetVertices()
                win.Update()
                render.ResetVertices()

                //process events
                for _, e := range win.Events() {
                        handled := false
                        switch e.Type {
                        case pixelgl.KeyPress:
                                if e.Key == pixelgl.KeySpace {
                                        handled = true
                                } else if e.Key == pixelgl.KeyEscape {
                                        win.SetClosed(true)
                                        handled = true
                                }
                        }
                        if !handled {
                                container.ProcessEvent(e)
                        }
                }

                elapsed := time.Since(start)
                frameDt = elapsed.Milliseconds()
                container.Update(frameDt)
        }
}


/*
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
                //Push glyph quads
                //hello := "The quick brown fox jumps over the lazy dog. 1234567890"
                extras := 10
                if len(hello) == 0 || !freeze {
                        hello = generateGlyphs(rows + extras, cols + extras)
                }

                penX, penY := float32(border), float32(ih - border - font.Ascender)
                for _, s := range hello {
                        shaper.DrawString(s, penX, penY, pixelgl.SolBase1)
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
                        shaper.DrawQuad(fX, fY, fW, fH, pixelgl.SolCyan)
                        shaper.DrawString(s, fX, fY, pixelgl.SolBase3)
                }

                //draw directory triangle thing in top left corner
                qW, qH := float32(lineSpace - 2), float32(lineSpace - 2)
                qX, qY := float32(2), h - qH - 2
                shaper.DrawQuad(qX, qY, qW, qH, pixelgl.Yellow)
                if closed {
                        shaper.DrawPointRight(qX, qY, qW, qH, pixelgl.SolRed) 
                } else {
                        shaper.DrawPointDown(qX, qY, qW, qH, pixelgl.SolRed) 
                }
*/
