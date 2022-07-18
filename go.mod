module github.com/artex2000/codeview

go 1.18

require (
	github.com/artex2000/codeview/font v0.0.0
	github.com/artex2000/codeview/thirdparty/pixelgl v0.0.0
	github.com/artex2000/codeview/shaper v0.0.0
)

require (
	github.com/artex2000/codeview/thirdparty/binpack v0.0.0 // indirect
	github.com/artex2000/codeview/thirdparty/freetype v0.0.0 // indirect
	github.com/artex2000/codeview/thirdparty/glhf v0.0.0 // indirect
	github.com/faiface/mainthread v0.0.0-20171120011319-8b78f0a41ae3 // indirect
	github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20220516021902-eb3e265c7661 // indirect
	github.com/go-gl/mathgl v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/image v0.0.0-20190321063152-3fc05d484e9f // indirect
)

replace (
	github.com/artex2000/codeview/font => ./font
	github.com/artex2000/codeview/shaper => ./shaper
	github.com/artex2000/codeview/thirdparty/binpack => ./thirdparty/binpack
	github.com/artex2000/codeview/thirdparty/freetype => ./thirdparty/freetype
	github.com/artex2000/codeview/thirdparty/glhf => ./thirdparty/glhf
	github.com/artex2000/codeview/thirdparty/pixelgl => ./thirdparty/pixelgl
)
