package glhf

import (
        "fmt"
	"runtime"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.3-core/gl"
)

type VertexArray struct {
	vao, vbo, veo   binder
        indices         []uint32
        vertices        []float32
        stride          int
}

func NewVertexArray(shader *Shader, attributes AttrFormat) (*VertexArray, error) {
	va := &VertexArray{
		vao: binder{
			restoreLoc: gl.VERTEX_ARRAY_BINDING,
			bindFunc: func(obj uint32) {
				gl.BindVertexArray(obj)
			},
		},
		vbo: binder{
			restoreLoc: gl.ARRAY_BUFFER_BINDING,
			bindFunc: func(obj uint32) {
				gl.BindBuffer(gl.ARRAY_BUFFER, obj)
			},
		},
		veo: binder{
			restoreLoc: gl.ELEMENT_ARRAY_BUFFER_BINDING,
			bindFunc: func(obj uint32) {
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, obj)
			},
		},
	}

	gl.GenVertexArrays(1, &va.vao.obj)

	va.vao.bind()

	gl.GenBuffers(1, &va.veo.obj)

	gl.GenBuffers(1, &va.vbo.obj)
	defer va.vbo.bind().restore()


        va.stride = attributes.Size()
        offset := 0
	for _, attr := range attributes {
		loc := gl.GetAttribLocation(shader.program.obj, gl.Str(attr.Name+"\x00"))

		var size int32
		switch attr.Type {
		case Float:
			size = 1
		case Vec2:
			size = 2
		case Vec3:
			size = 3
		case Vec4:
			size = 4
                default:
                        return nil, fmt.Errorf("Invalid attribute type")
		}

		gl.VertexAttribPointerWithOffset(
			uint32(loc),
			size,
			gl.FLOAT,
			false,
			int32(va.stride),
			uintptr(offset),
		)
		gl.EnableVertexAttribArray(uint32(loc))

                offset += attr.Type.Size()
	}

	va.vao.restore()

	runtime.SetFinalizer(va, (*VertexArray).delete)

	return va, nil
}

func (va *VertexArray) delete() {
	mainthread.CallNonBlock(func() {
		gl.DeleteVertexArrays(1, &va.vao.obj)
		gl.DeleteBuffers(1, &va.vbo.obj)
		gl.DeleteBuffers(1, &va.veo.obj)
	})
}

func (va *VertexArray) Begin() {
	va.vao.bind()
	va.vbo.bind()
	va.veo.bind()
}

func (va *VertexArray) End() {
	va.veo.restore()
	va.vbo.restore()
	va.vao.restore()
}

func (va *VertexArray) Draw(count int32) {
	gl.DrawElements(gl.TRIANGLES, count, gl.UNSIGNED_INT, gl.Ptr(uintptr(0)))
}

func (va *VertexArray) SetVertexData(v []float32, i []uint32) {
	gl.BufferData(gl.ARRAY_BUFFER, len(v)*4, gl.Ptr(v), gl.DYNAMIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(i)*4, gl.Ptr(i), gl.DYNAMIC_DRAW)
}

func (va *VertexArray) Stride() int {
        return va.stride
}
