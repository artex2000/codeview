package font

import (
	"fmt"
        "os"
        "image"
        "image/draw"
        "path/filepath"
        "sort"

	"github.com/artex2000/codeview/thirdparty/freetype"
	"github.com/artex2000/codeview/thirdparty/binpack"
       )

type Glyph struct {
        OffsetX         int
        OffsetY         int
        Advance         int
        //TODO: Should the size be as reported by freetype or extended to apron?
        Width           int
        Height          int

        TexS0           float32
        TexT0           float32
        TexS1           float32
        TexT1           float32
}

type Monofont struct {
        Atlas           *image.RGBA
        FontSize        int
        Ascender        int
        Descender       int
        Linegap         int
        SpaceAdvance    int
        Glyphs          map[rune]Glyph
}

type glyph struct {
        char            rune
        metrics         *freetype.Metrics
        image           *image.RGBA
        atlasX          int
        atlasY          int
        atlasW          int
        atlasH          int
}

type font       []*glyph

func (f font) Len() int {
        return len(f)
}

func (f font) Less(i, j int) bool {
        return f[i].atlasH < f[j].atlasH
}

func (f font) Swap(i, j int) {
        f[i], f[j] = f[j], f[i]
}

func (f font) Size(n int) (width, height int) {
        width, height = f[n].atlasW, f[n].atlasH
        return
}

func (f font) Place(n, x, y int) {
        f[n].atlasX, f[n].atlasY = x, y
}

func InitFontFromFile(fname string, size int) (*Monofont, error) {
        mf := &Monofont{}
        mf.FontSize = size
        mf.Glyphs   = make(map[rune]Glyph)

        fileName, _ := filepath.Abs(fname)
        fontData, err := os.ReadFile(fileName)
        if err != nil {
                return nil, fmt.Errorf("InitFont failed - %w", err)
        }

	lib, err := freetype.NewLibrary()
        if err != nil {
                return nil, fmt.Errorf("InitFont failed - %w", err)
        }
        defer lib.Done()

	face, err := freetype.NewFace(lib, fontData, 0)
        if err != nil {
                return nil, fmt.Errorf("InitFont failed - %w", err)
        }
        defer face.Done()

        //TODO: Support custom-supplied DPI
	err = face.Pt(size, 96)
        if err != nil {
                return nil, fmt.Errorf("InitFont failed - %w", err)
        }

        //Deal with space
        _, m, err := face.GlyphEx(0x20)
        if err != nil {
                return nil, fmt.Errorf("InitFont failed - %w", err)
        }

        mf.SpaceAdvance = m.AdvanceWidth

        var f font
        for i := rune(0x21); i < 0x7F; i += 1 {
                img, m, err := face.GlyphEx(i)
                if err != nil {
                        continue
                }

                //we're adding 1px padding on the right and bottom so atlas neighbours won't affect colors
                g := &glyph{char: i,
                            metrics: m, 
                            image : img, 
                            atlasW : img.Bounds().Dx() + 1, 
                            atlasH : img.Bounds().Dy() + 1,
                   }
                f = append(f, g)
                //fmt.Printf("%d - %d ......... %d - %d\n", m.Width, m.Height, img.Bounds().Dx(), img.Bounds().Dy())
        }

        //Sort glyphs by height - it gives good results with binpack
        sort.Sort(sort.Reverse(f))
        //256x256 is enough to pack 32pt font.
        //TODO: Expand resulting atlas if necessary to square multiple of 4 pixels
        w, h := binpack.Pack(f, 256, 256)
        if (w > 256) || (h > 256) {
                return nil, fmt.Errorf("Atlas too small")
        }

        mf.Atlas = image.NewRGBA(image.Rect(0, 0, w, h))

        maxAsc, maxDsc := 0, 0
        for _, g := range f {
                draw.Draw(mf.Atlas, image.Rect(g.atlasX, g.atlasY, g.atlasX + g.atlasW - 1, g.atlasY + g.atlasH - 1), g.image, image.Pt(0, 0), draw.Src)
                
                gl := Glyph{}
                gl.OffsetX = g.metrics.HorizontalBearingX
                gl.OffsetY = g.metrics.HorizontalBearingY
                //gl.Width   = g.metrics.Width
                gl.Width   = g.atlasW - 1
                gl.Height  = g.metrics.Height
                gl.Advance = g.metrics.AdvanceWidth

                gl.TexS0 = float32(g.atlasX) / float32(w)
                gl.TexT0 = float32(g.atlasY) / float32(h)
                gl.TexS1 = float32(g.atlasX + g.atlasW - 1) / float32(w)
                gl.TexT1 = float32(g.atlasY + g.atlasH) / float32(h)

                //Find biggest ascender
                if gl.OffsetY > 0 {
                        if maxAsc < gl.OffsetY {
                                maxAsc = gl.OffsetY
                        }
                }

                //Find biggest descender
                if gl.Height > gl.OffsetY {
                        dsc := gl.Height - gl.OffsetY
                        if maxDsc < dsc {
                                maxDsc = dsc
                        }
                }

                mf.Glyphs[g.char] = gl
        }

        mf.Ascender  = maxAsc
        mf.Descender = maxDsc
        mf.Linegap   = int(0.1 * float32(mf.Ascender + mf.Descender))

        return mf, nil
}


