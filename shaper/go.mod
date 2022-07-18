module github.com/artex2000/codeview/shaper

go 1.18

require (
	github.com/artex2000/codeview/font v0.0.0
	github.com/artex2000/codeview/thirdparty/pixelgl v0.0.0
)

replace (
	github.com/artex2000/codeview/font => ../font
	github.com/artex2000/codeview/thirdparty/pixelgl => ../thirdparty/pixelgl
)
