package main

import (
//        "fmt"

        "github.com/artex2000/codeview/thirdparty/pixelgl"
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
        bounds          pixelgl.Rect
        widgets         []Widget
        widgetErr       chan error
        focusIndex      int
        palette         map[string]pixelgl.RGBA
}

func NewContainer(rect pixelgl.Rect) *Container {
        out := &Container{bounds: rect}
        out.widgets = make([]Widget, 5)
        out.widgetErr = make(chan error, 5)
        return out
}

func (c *Container) AddTree(root string) {
        w := Widget{Type: Tree, Errors: c.widgetErr}
        xmin, ymin := c.bounds.Min.X, c.bounds.Min.Y
        xmax, ymax := xmin + c.bounds.W() / 3, c.bounds.Max.Y
        w.Bounds = pixelgl.R(xmin, ymin, xmax, ymax)
        
        tree := &TreeData{}
        tree.RootEntry = &DirEntry{Name: root, Parent: nil, Level: 0, Flags: (Directory)}
        tree.Cache     = make(map[DirEntry][]*DirEntry, 5)
        tree.ExpandEntry(tree.RootEntry, true)
        w.Data = tree

        c.widgets = append(c.widgets)
}

func (c *Container) DrawTree(w Widget) {
        t := w.Data.(*TreeData)
        tree := t.Lines[t.ScrollIndex:]

        //Set up drawing origin and calculate visible lines count
        numLines := 1

        for i, d := range tree {
                if i == numLines {
                        break
                }

                if i == t.CurrentLine {
                        //Change background color and draw background box
                }

                //Calculate horizontal indent

                if d.IsDir() {
                        if (d.Flags & Expanded) != 0 {
                                //draw PointDown triangle
                        } else {
                                //draw PointRight triangle
                        }
                } else {
                }
                //Update drawing coordinates
        }
}
