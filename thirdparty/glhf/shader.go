package glhf

import (
	"fmt"
	"runtime"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.3-core/gl"
)

// Shader is an OpenGL shader program.
type Shader struct {
	program       binder
}

// NewShader creates a new shader program from the specified vertex shader and fragment shader
// sources.
//
// Note that vertexShader and fragmentShader parameters must contain the source code, they're
// not filenames.
func NewShader(vertexShader, fragmentShader string) (*Shader, error) {
	shader := &Shader{
		program: binder{
			restoreLoc: gl.CURRENT_PROGRAM,
			bindFunc: func(obj uint32) {
				gl.UseProgram(obj)
			},
		},
	}

	var vshader, fshader uint32

	// vertex shader
	{
		vshader = gl.CreateShader(gl.VERTEX_SHADER)
		src, free := gl.Strs(vertexShader)
		defer free()
		length := int32(len(vertexShader))
		gl.ShaderSource(vshader, 1, src, &length)
		gl.CompileShader(vshader)

		var success int32
		gl.GetShaderiv(vshader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(vshader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(vshader, logLen, nil, &infoLog[0])
			return nil, fmt.Errorf("error compiling vertex shader: %s", string(infoLog))
		}

		defer gl.DeleteShader(vshader)
	}

	// fragment shader
	{
		fshader = gl.CreateShader(gl.FRAGMENT_SHADER)
		src, free := gl.Strs(fragmentShader)
		defer free()
		length := int32(len(fragmentShader))
		gl.ShaderSource(fshader, 1, src, &length)
		gl.CompileShader(fshader)

		var success int32
		gl.GetShaderiv(fshader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(fshader, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetShaderInfoLog(fshader, logLen, nil, &infoLog[0])
			return nil, fmt.Errorf("error compiling fragment shader: %s", string(infoLog))
		}

		defer gl.DeleteShader(fshader)
	}

	// shader program
	{
		shader.program.obj = gl.CreateProgram()
		gl.AttachShader(shader.program.obj, vshader)
		gl.AttachShader(shader.program.obj, fshader)
		gl.LinkProgram(shader.program.obj)

		var success int32
		gl.GetProgramiv(shader.program.obj, gl.LINK_STATUS, &success)
		if success == gl.FALSE {
			var logLen int32
			gl.GetProgramiv(shader.program.obj, gl.INFO_LOG_LENGTH, &logLen)

			infoLog := make([]byte, logLen)
			gl.GetProgramInfoLog(shader.program.obj, logLen, nil, &infoLog[0])
			return nil, fmt.Errorf("error linking shader program: %s", string(infoLog))
		}
	}

	runtime.SetFinalizer(shader, (*Shader).delete)

	return shader, nil
}

func (s *Shader) delete() {
	mainthread.CallNonBlock(func() {
		gl.DeleteProgram(s.program.obj)
	})
}

// ID returns the OpenGL ID of this Shader.
func (s *Shader) ID() uint32 {
	return s.program.obj
}

// Begin binds the Shader program. This is necessary before using the Shader.
func (s *Shader) Begin() {
	s.program.bind()
}

// End unbinds the Shader program and restores the previous one.
func (s *Shader) End() {
	s.program.restore()
}
