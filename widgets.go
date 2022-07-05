package main

import (
        "fmt"

        "github.com/artex2000/codeview/thirdparty/pixelgl"
)

type WidgetType int
type WidgetFlags uint64

const (
        Button = WidgetType(0)
        Editor = WidgetType(1)
        Tree   = WidgetType(2)
)

type Widget struct {
        Type            WidgetType
        Bounds          pixelgl.Rect
        Foreground      pixelgl.RGBA
        Background      pixelgl.RGBA
        Flags           WidgetFlags
        Data            interface{}
}

type Container struct {
        widgets         []Widget
        focusIndex      int
}

type TreeData struct {
        Root    string
}

func InitContainer() error {
        return fmt.Errorf("Not implemented")
}
