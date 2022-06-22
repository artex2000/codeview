package pixelgl

import (
        "fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.3-core/gl"
        "github.com/artex2000/codeview/thirdparty/glhf"
       )

type Variable struct {
        Name            string
        Type            glhf.AttrType
}

type VariableList []Variable

type ShaderVar struct {
        Name            string
        Type            glhf.AttrType
        Location        int32
        LocationValid   bool
}

type Render struct {
        Shader          *glhf.Shader
        VertexArray     *glhf.VertexArray
        Uniforms        []*ShaderVar
        count           int32

        Model           mgl32.Mat4
        View            mgl32.Mat4
        Projection      mgl32.Mat4
}

func NewRender(vert, frag string, uniforms, attributes VariableList) *Render {
        r := &Render{}

        var shader *glhf.Shader
        var vertexArray *glhf.VertexArray

        if attributes == nil {
                panic("Can't init shader without vertex attributes")
        }

        mainthread.Call(func() {
                var err error
                shader, err = glhf.NewShader(vert, frag)
                if err != nil {
                        panic(fmt.Sprintf("Shader compile error - %w", err))
                }

                va := make(glhf.AttrFormat, len(attributes))
                for i, a := range attributes {
                        va[i] = glhf.Attr{Name: a.Name, Type: a.Type}
                }

                vertexArray, err = glhf.NewVertexArray(shader, va)
                if err != nil {
                        panic(fmt.Sprintf("Vertex array init error - %w", err))
                }
        })

        if uniforms != nil {
                r.Uniforms = make([]*ShaderVar, len(uniforms))
                for i, u := range uniforms {
                        r.Uniforms[i] = &ShaderVar{Name: u.Name, Type: u.Type}
                }
        }

        r.VertexArray = vertexArray 
        r.Shader      = shader
        return r
}

//Must be called from main thread
func (r *Render) SetUniformByName(Name string, value interface{}) {
        if r.Uniforms == nil {
                return
        }

        for _, v := range r.Uniforms {
                if v.Name == Name {
                        if !v.LocationValid {
                                v.Location = gl.GetUniformLocation(r.Shader.ID(), gl.Str(Name + "\x00"))
                                if v.Location == -1 {
                                        panic(fmt.Sprintf("Name %s is not in the shader", Name))
                                }
                                v.LocationValid = true
                        }
                        r.SetUniform(v, value)
                        break
                }
        }
}

//Must be called from main thread
func (r *Render) SetUniform(v *ShaderVar, value interface{}) {
	switch v.Type {
	case glhf.Int:
		value := value.(int32)
		gl.Uniform1iv(v.Location, 1, &value)
	case glhf.Float:
		value := value.(float32)
		gl.Uniform1fv(v.Location, 1, &value)
	case glhf.Vec2:
		value := value.(mgl32.Vec2)
		gl.Uniform2fv(v.Location, 1, &value[0])
	case glhf.Vec3:
		value := value.(mgl32.Vec3)
		gl.Uniform3fv(v.Location, 1, &value[0])
	case glhf.Vec4:
		value := value.(mgl32.Vec4)
		gl.Uniform4fv(v.Location, 1, &value[0])
	case glhf.Mat2:
		value := value.(mgl32.Mat2)
		gl.UniformMatrix2fv(v.Location, 1, false, &value[0])
	case glhf.Mat23:
		value := value.(mgl32.Mat2x3)
		gl.UniformMatrix2x3fv(v.Location, 1, false, &value[0])
	case glhf.Mat24:
		value := value.(mgl32.Mat2x4)
		gl.UniformMatrix2x4fv(v.Location, 1, false, &value[0])
	case glhf.Mat3:
		value := value.(mgl32.Mat3)
		gl.UniformMatrix3fv(v.Location, 1, false, &value[0])
	case glhf.Mat32:
		value := value.(mgl32.Mat3x2)
		gl.UniformMatrix3x2fv(v.Location, 1, false, &value[0])
	case glhf.Mat34:
		value := value.(mgl32.Mat3x4)
		gl.UniformMatrix3x4fv(v.Location, 1, false, &value[0])
	case glhf.Mat4:
		value := value.(mgl32.Mat4)
		gl.UniformMatrix4fv(v.Location, 1, false, &value[0])
	case glhf.Mat42:
		value := value.(mgl32.Mat4x2)
		gl.UniformMatrix4x2fv(v.Location, 1, false, &value[0])
	case glhf.Mat43:
		value := value.(mgl32.Mat4x3)
		gl.UniformMatrix4x3fv(v.Location, 1, false, &value[0])
	default:
		panic("set uniform attr: invalid attribute type")
	}
}

func (r *Render) SetTranslationMatrix(width, height float32) {
        r.Model = mgl32.Ident4()
        r.View  = mgl32.Ident4()
        r.Projection = mgl32.Ident4()
        //r.Projection = mgl32.Ortho2D(0, width, 0, height)

        mainthread.Call(func() {
                r.SetUniformByName("Model", r.Model)
                r.SetUniformByName("View",  r.View)
                r.SetUniformByName("Projection", r.Projection)
        })
}

func (r *Render) SetVertices(f []float32, i []uint32) {
        mainthread.Call(func() {
                r.VertexArray.Begin()
                r.VertexArray.SetVertexData(f, i)
                r.VertexArray.End()
        })
        r.count = int32(len(i))
}

func (r *Render) Draw() {
        r.Shader.Begin()
        r.VertexArray.Begin()
        r.VertexArray.Draw(r.count)
        r.VertexArray.End()
        r.Shader.End()
}
