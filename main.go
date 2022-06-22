package main

import (
        "fmt"

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

        win, err := pixelgl.NewWindow(cfg)
        if err != nil {
                panic (err)
        }

        /*
        Uniforms := pixelgl.VariableList{
                        { Name: "Model",        Type: glhf.Mat4 },
                        { Name: "View",         Type: glhf.Mat4 },
                        { Name: "Projection",   Type: glhf.Mat4 },
                    }
        */

        Attributes := pixelgl.VariableList{
                        { Name: "VertexPosition", Type: glhf.Vec3 },
                        { Name: "VertexColor",    Type: glhf.Vec3 },
                    }

        render := pixelgl.NewRender(VertexShader, FragmentShader, nil, Attributes)

        vs := []float32 {
                -0.5, -0.5, 0.0, 1.0, 0.0, 0.0,
                 0.5, -0.5, 0.0, 0.0, 1.0, 0.0,
                 0.0,  0.35, 0.0, 0.0, 0.0, 1.0,
        }

        is := []uint32 { 0, 1, 2 }

        render.SetVertices(vs, is)
//        render.SetTranslationMatrix(0, 0)
        win.SetCanvas(render)

        for !win.Closed() {
//                win.Clear()
                win.Update()
        }
}
