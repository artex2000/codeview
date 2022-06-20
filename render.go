package main

import (
        "fmt"

	"github.com/faiface/mainthread"
        "github.com/artex2000/codeview/thirdparty/glhf"
       )

type Render struct {
        Shader          *glhf.Shader
}

func NewRender(vert, frag string) (*Render, error) {
        r := &Render{}

        var shader *glhf.Shader
        mainthread.Call(func() {
                var err error
                shader, err = glhf.NewShader(vert, frag, nil)
                if err != nil {
                        panic(fmt.Sprintf("Shader compile error - %w", err))
                }
        })

        r.Shader = shader
        return r, nil
}
