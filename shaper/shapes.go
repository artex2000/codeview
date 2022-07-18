package shaper

import (
        "math"
        "github.com/artex2000/codeview/thirdparty/pixelgl"
)

type Direction int

const (
        East  = Direction(0)
        West  = Direction(1)
        North = Direction(2)
        South = Direction(3)
)

var halfSquareOfThree   float64

func init() {
        halfSquareOfThree = math.Sqrt(3) / 2
}

func EquilateralTriangle(bounds pixelgl.Rect, direction Direction) [3]pixelgl.Vec {
        var out [3]pixelgl.Vec

        r := bounds.Norm()
        width, height := r.W(), r.H()
        if width != height {
                panic("EquilateralTriangle: Bounding box is not a square")
        }

        radius := r.W() / 2
        middlePoint := pixelgl.Vec{r.Min.X + radius, r.Min.Y + radius}

        //We will find offsets from the middle point in first quadrant of bounding square and
        //flip them as needed, depending on the direction requested
        xOf30GradOff := radius * halfSquareOfThree 
        yOf30GradOff := radius * 0.5
        xOf60GradOff := radius * 0.5
        yOf60GradOff := radius * halfSquareOfThree

        switch direction {
        case East:
                out[0] = middlePoint.Add(pixelgl.Vec{-xOf60GradOff, yOf60GradOff})         //left top
                out[1] = middlePoint.Add(pixelgl.Vec{-xOf60GradOff, -yOf60GradOff})        //left bottom
                out[2] = middlePoint.Add(pixelgl.Vec{radius, 0})                           //right middle
        case West:
                out[0] = middlePoint.Add(pixelgl.Vec{-radius, 0})                          //left middle
                out[1] = middlePoint.Add(pixelgl.Vec{xOf60GradOff, -yOf60GradOff})         //right bottom
                out[2] = middlePoint.Add(pixelgl.Vec{xOf60GradOff, yOf60GradOff})          //right top
        case South:
                out[0] = middlePoint.Add(pixelgl.Vec{-xOf30GradOff, yOf30GradOff})         //left top
                out[1] = middlePoint.Add(pixelgl.Vec{0, -radius})                          //middle bottom
                out[2] = middlePoint.Add(pixelgl.Vec{xOf30GradOff, yOf30GradOff})          //right top
        case North:
                out[0] = middlePoint.Add(pixelgl.Vec{0, radius})                           //middle top
                out[1] = middlePoint.Add(pixelgl.Vec{-xOf30GradOff, -yOf30GradOff})        //left bottom
                out[2] = middlePoint.Add(pixelgl.Vec{xOf30GradOff, -yOf30GradOff})         //right bottom
        }
        return out
}

