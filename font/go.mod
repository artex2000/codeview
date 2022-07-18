module github.com/artex2000/codeview/font

go 1.18

require (
	github.com/artex2000/codeview/thirdparty/freetype v0.0.0
	github.com/artex2000/codeview/thirdparty/binpack v0.0.0
)

replace (
	github.com/artex2000/codeview/thirdparty/binpack => ../thirdparty/binpack
	github.com/artex2000/codeview/thirdparty/freetype => ../thirdparty/freetype
)
