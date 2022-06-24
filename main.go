package main

import (
        "fmt"

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
                    }

        r := pixelgl.NewRender(VertexShader, FragmentShader, Uniforms, Attributes)

        vs := []float32 {
                -0.5, -0.5, 0.0, 1.0, 0.0, 0.0,
                 0.5, -0.5, 0.0, 0.0, 1.0, 0.0,
                 0.0,  0.35, 0.0, 0.0, 0.0, 1.0,
        }

        model := mgl32.Ident4()

        eye    := mgl32.Vec3{ 0.0, 0.0, 2.0 }
        center := mgl32.Vec3{ 0.0, 0.0, 0.0 }
        up     := mgl32.Vec3{ 0.0, 1.0, 0.0 }
        view  := mgl32.LookAtV(eye, center, up)

        projection := mgl32.Ortho(-1, 1, -1, 1, 0.1, 5)
        //projection = mgl32.Ortho(0, width, 0, height, 0.1, 5)

        r.SetModelMatrix(model)
        r.SetViewMatrix(view)
        r.SetProjectionMatrix(projection)
        r.SetTransform(true, true, true)

        r.PushTriangle(vs)
        r.SetVertices()

        win.SetCanvas(r)

        for !win.Closed() {
//                win.Clear()
                win.Update()
        }
}
