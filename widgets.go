package main

import (
//        "fmt"
        "math"

        "github.com/artex2000/codeview/thirdparty/pixelgl"
        "github.com/artex2000/codeview/shaper"
)

type WidgetFlags uint64
const (
        CanHover = WidgetFlags(1 << iota)
)

type WidgetType int
const (
        Button = WidgetType(iota)
        Editor
        Tree
)

type Widget struct {
        Type            WidgetType
        Bounds          pixelgl.Rect
        /*
        Foreground      pixelgl.RGBA
        Background      pixelgl.RGBA
        */
        Flags           WidgetFlags
        Errors          <-chan error
        Data            interface{}
}

type Container struct {
        shaper          *shaper.Shaper
        bounds          pixelgl.Rect
        widgets         []*Widget
        widgetErr       chan error
        focusIndex      int
        palette         map[string]pixelgl.RGBA
}

func NewContainer(s *shaper.Shaper, r pixelgl.Rect) *Container {
        out := &Container{shaper: s, bounds: r}
        out.widgetErr = make(chan error, 5)
        return out
}

func (c *Container) AddTree(root string) {
        w := Widget{Type: Tree, Errors: c.widgetErr}
        xmin, ymin := c.bounds.Min.X, c.bounds.Min.Y
        xmax, ymax := xmin + c.bounds.W() / 3, c.bounds.Max.Y
        w.Bounds = pixelgl.R(xmin, ymin, xmax, ymax)
        
        tree := &TreeData{}
        tree.RootEntry = &DirEntry{Name: root, Parent: nil, Level: 0, Flags: (Directory | Root)}
        tree.Cache     = make(map[DirEntry][]*DirEntry, 5)
        tree.ExpandEntry(tree.RootEntry, true)
        w.Data = tree

        c.widgets = append(c.widgets, &w)
}

func (c *Container) DrawTree(w *Widget) {
        t := w.Data.(*TreeData)
        tree := t.Lines[t.ScrollIndex:]
        rect := w.Bounds
        treeX, treeY := float32(rect.Min.X), float32(rect.Min.Y)
        //treeH := float32(rect.H())
        treeW, treeH := float32(rect.W()), float32(rect.H())

        //Maximum line height of the given font
        lineH := float32(c.shaper.Font.Ascender + c.shaper.Font.Descender + c.shaper.Font.Linegap)

        //origin point describes LEFT-BOTTOM corner of bounding box of the top drawable line in tree
        boxX, boxY := treeX, treeY + (treeH - lineH)

        //origin point describes baseline point of the first drawable character of the top line
        penX, penY := treeX, treeY + (treeH - float32(c.shaper.Font.Ascender))

        maxLines   := int(math.Ceil(float64(treeH / lineH)))

        for i, d := range tree {
                if i == maxLines {
                        break
                }

                if i == t.CurrentLine {
                        //Change background color and draw background box
                        c.shaper.DrawQuad(boxX, boxY, treeW, lineH, pixelgl.SolBase02)
                }

                //Calculate horizontal indent
                //we have lineHeight square for dir triangle, so it's our level step
                indent := float32(d.Level) * lineH
                boxX += indent
                penX += indent

                if d.IsDir() {
                        if (d.Flags & Expanded) != 0 {
                                //draw PointDown triangle
                                c.shaper.DrawEqTriangleDown(boxX, boxY, lineH, pixelgl.SolBase0)
                        } else {
                                //draw PointRight triangle
                                c.shaper.DrawEqTriangleRight(boxX, boxY, lineH, pixelgl.SolBase0)
                        }
                        penX += lineH
                }

                // Draw file/directory name
                c.shaper.DrawString(d.Name, penX, penY, pixelgl.SolBase0)

                //Update drawing coordinates
                boxX = treeX
                penX = treeY
                boxY -= lineH
                penY -= lineH
        }
}

func (c *Container) Draw() {
        for _, w := range c.widgets {
                switch w.Type {
                case Tree:
                        c.DrawTree(w)
                }
        }
}

func (c *Container) Update(ms int64) {
}

func (c *Container) ProcessEvent(e pixelgl.Event) {
}
