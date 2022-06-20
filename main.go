package main

import (
        "fmt"

        "github.com/artex2000/codeview/thirdparty/pixelgl"
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

        render, err := NewRender(VertexShader, FragmentShader)
        if err != nil {
                panic (err)
        }

        fmt.Printf("Created shader with id %x\n", render.Shader.ID())

        for !win.Closed() {
//                win.Clear()
                win.Update()
        }
}
