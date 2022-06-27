package pixelgl

import (
	"fmt"
	//"image"
	"runtime"

	"github.com/artex2000/codeview/thirdparty/glhf"
	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/pkg/errors"
)

// WindowConfig is a structure for specifying all possible properties of a Window. Properties are
// chosen in such a way, that you usually only need to set a few of them - defaults (zeros) should
// usually be sensible.
//
// Note that you always need to set the Bounds of a Window.
type WindowConfig struct {
	// Title at the top of the Window.
	Title string

	// Icon specifies the icon images available to be used by the window. This is usually
	// displayed in the top bar of the window or in the task bar of the desktop environment.
	//
	// If passed one image, it will use that image, if passed an array of images those of or
	// closest to the sizes desired by the system are selected. The desired image sizes varies
	// depending on platform and system settings. The selected images will be rescaled as
	// needed. Good sizes include 16x16, 32x32 and 48x48.
	//
	// Note: Setting this value doesn't have an effect on OSX. You'll need to set the icon when
	// bundling your application for release.
	//Icon []Picture

	// Bounds specify the bounds of the Window in pixels.
	Bounds Rect

	// Initial window position
	Position Vec

	// If set to nil, the Window will be windowed. Otherwise it will be fullscreen on the
	// specified Monitor.
	Monitor *Monitor

	// Resizable specifies whether the window will be resizable by the user.
	Resizable bool

	// Undecorated Window omits the borders and decorations (close button, etc.).
	Undecorated bool

	// NoIconify specifies whether fullscreen windows should not automatically
	// iconify (and restore the previous video mode) on focus loss.
	NoIconify bool

	// AlwaysOnTop specifies whether the windowed mode window will be floating
	// above other regular windows, also called topmost or always-on-top.
	// This is intended primarily for debugging purposes and cannot be used to
	// implement proper full screen windows.
	AlwaysOnTop bool

	// TransparentFramebuffer specifies whether the window framebuffer will be
	// transparent. If enabled and supported by the system, the window
	// framebuffer alpha channel will be used to combine the framebuffer with
	// the background. This does not affect window decorations.
	TransparentFramebuffer bool

	// VSync (vertical synchronization) synchronizes Window's framerate with the framerate of
	// the monitor.
	VSync bool

	// Maximized specifies whether the window is maximized.
	Maximized bool

	// Invisible specifies whether the window will be initially hidden.
	// You can make the window visible later using Window.Show().
	Invisible bool

	//SamplesMSAA specifies the level of MSAA to be used. Must be one of 0, 2, 4, 8, 16. 0 to disable.
	SamplesMSAA int
}

// Window is a window handler. Use this type to manipulate a window (input, drawing, etc.).
type Window struct {
	window *glfw.Window

        canvas             *Render
	bounds             Rect
	vsync              bool
	cursorVisible      bool
	cursorInsideWindow bool

	// need to save these to correctly restore a fullscreen window
	restore struct {
		xpos, ypos, width, height int
	}

	prevInp, currInp, tempInp struct {
		mouse   Vec
		buttons [KeyLast + 1]bool
		repeat  [KeyLast + 1]bool
		scroll  Vec
		typed   string
	}

	pressEvents, tempPressEvents     [KeyLast + 1]bool
	releaseEvents, tempReleaseEvents [KeyLast + 1]bool

	prevJoy, currJoy, tempJoy joystickState
}

var currWin *Window

// NewWindow creates a new Window with it's properties specified in the provided config.
//
// If Window creation fails, an error is returned (e.g. due to unavailable graphics device).
func NewWindow(cfg WindowConfig) (*Window, error) {
	bool2int := map[bool]int{
		true:  glfw.True,
		false: glfw.False,
	}

	w := &Window{bounds: cfg.Bounds, cursorVisible: true}

	flag := false
	for _, v := range []int{0, 2, 4, 8, 16} {
		if cfg.SamplesMSAA == v {
			flag = true
			break
		}
	}
	if !flag {
		return nil, fmt.Errorf("invalid value '%v' for msaaSamples", cfg.SamplesMSAA)
	}

	err := mainthread.CallErr(func() error {
		var err error

		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

		glfw.WindowHint(glfw.Resizable, bool2int[cfg.Resizable])
		glfw.WindowHint(glfw.Decorated, bool2int[!cfg.Undecorated])
		glfw.WindowHint(glfw.Floating, bool2int[cfg.AlwaysOnTop])
		glfw.WindowHint(glfw.AutoIconify, bool2int[!cfg.NoIconify])
		glfw.WindowHint(glfw.TransparentFramebuffer, bool2int[cfg.TransparentFramebuffer])
		glfw.WindowHint(glfw.Maximized, bool2int[cfg.Maximized])
		glfw.WindowHint(glfw.Visible, bool2int[!cfg.Invisible])
		glfw.WindowHint(glfw.Samples, cfg.SamplesMSAA)

		if cfg.Position.X != 0 || cfg.Position.Y != 0 {
			glfw.WindowHint(glfw.Visible, glfw.False)
		}

		var share *glfw.Window
		if currWin != nil {
			share = currWin.window
		}
		_, _, width, height := cfg.Bounds.IntBounds()
		w.window, err = glfw.CreateWindow(
			width,
			height,
			cfg.Title,
			nil,
			share,
		)
		if err != nil {
			return err
		}

		if cfg.Position.X != 0 || cfg.Position.Y != 0 {
			w.window.SetPos(int(cfg.Position.X), int(cfg.Position.Y))
			w.window.Show()
		}

		// enter the OpenGL context
		w.begin()
		glhf.Init()
		gl.Enable(gl.MULTISAMPLE)
		w.end()

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating window failed")
	}

        /*
	if len(cfg.Icon) > 0 {
		imgs := make([]image.Image, len(cfg.Icon))
		for i, icon := range cfg.Icon {
			pic := PictureDataFromPicture(icon)
			imgs[i] = pic.Image()
		}
		mainthread.Call(func() {
			w.window.SetIcon(imgs)
		})
	}
        */

	w.SetVSync(cfg.VSync)

	w.initInput()

	mainthread.Call(func() {
                w.window.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
                        //fmt.Printf("Framebuffer size changed to %d : %d\n", width, height)
                        glhf.Bounds(0, 0, width, height)
                        glhf.Clear(0.8, 0.8, 0.9, 0.0)
                        if w.canvas != nil {
                                w.canvas.Draw()
                        }
                        w.window.SwapBuffers()
                })

		if w.vsync {
			glfw.SwapInterval(1)
		} else {
			glfw.SwapInterval(0)
		}

		framebufferWidth, framebufferHeight := w.window.GetFramebufferSize()
		glhf.Bounds(0, 0, framebufferWidth, framebufferHeight)

                /*
                w.window.SetSizeCallback(func(_ *glfw.Window, width, height int) {
                        fmt.Printf("Window size changed to %d : %d\n", width, height)
                })
                w.window.SetRefreshCallback(func(_ *glfw.Window) {
                        fmt.Println("Refresh callback")
                        w.SwapBuffers()
                })
                */
	})

	w.SetMonitor(cfg.Monitor)

	w.Update()

	runtime.SetFinalizer(w, (*Window).Destroy)

	return w, nil
}

// Destroy destroys the Window. The Window can't be used any further.
func (w *Window) Destroy() {
	mainthread.Call(func() {
		w.window.Destroy()
	})
}

// Update swaps buffers and polls events. Call this method at the end of each frame.
func (w *Window) Update() {
	w.SwapBuffers()
	w.UpdateInput()
}

// ClipboardText returns the current value of the systems clipboard.
func (w *Window) ClipboardText() string {
	return w.window.GetClipboardString()
}

// SetClipboardText passes the given string to the underlying glfw window to set the
//	systems clipboard.
func (w *Window) SetClipboardText(text string) {
	w.window.SetClipboardString(text)
}

// SwapBuffers swaps buffers. Call this to swap buffers without polling window events.
// Note that Update invokes SwapBuffers.
func (w *Window) SwapBuffers() {
        /*
	mainthread.Call(func() {
		_, _, oldW, oldH := w.bounds.IntBounds()
		newW, newH := w.window.GetSize()
		w.bounds = w.bounds.ResizedMin(w.bounds.Size().Add(V(
			float64(newW-oldW),
			float64(newH-oldH),
		)))
	})
        */

	mainthread.Call(func() {
		w.begin() //make context current if needed

		glhf.Clear(0.0, 0.0, 0.0, 0.0)
                if w.canvas != nil {
                        w.canvas.Draw()
                }
		w.window.SwapBuffers()
		w.end()
	})
}

// SetClosed sets the closed flag of the Window.
//
// This is useful when overriding the user's attempt to close the Window, or just to close the
// Window from within the program.
func (w *Window) SetClosed(closed bool) {
	mainthread.Call(func() {
		w.window.SetShouldClose(closed)
	})
}

// Closed returns the closed flag of the Window, which reports whether the Window should be closed.
//
// The closed flag is automatically set when a user attempts to close the Window.
func (w *Window) Closed() bool {
	var closed bool
	mainthread.Call(func() {
		closed = w.window.ShouldClose()
	})
	return closed
}

// SetTitle changes the title of the Window.
func (w *Window) SetTitle(title string) {
	mainthread.Call(func() {
		w.window.SetTitle(title)
	})
}

// SetBounds sets the bounds of the Window in pixels. Bounds can be fractional, but the actual size
// of the window will be rounded to integers.
func (w *Window) SetBounds(bounds Rect) {
	w.bounds = bounds
	mainthread.Call(func() {
		_, _, width, height := bounds.IntBounds()
		w.window.SetSize(width, height)
	})
}

// SetPos sets the position, in screen coordinates, of the upper-left corner
// of the client area of the window. Position can be fractional, but the actual position
// of the window will be rounded to integers.
//
// If it is a full screen window, this function does nothing.
func (w *Window) SetPos(pos Vec) {
	mainthread.Call(func() {
		left, top := int(pos.X), int(pos.Y)
		w.window.SetPos(left, top)
	})
}

// GetPos gets the position, in screen coordinates, of the upper-left corner
// of the client area of the window. The position is rounded to integers.
func (w *Window) GetPos() Vec {
	var v Vec
	mainthread.Call(func() {
		x, y := w.window.GetPos()
		v = V(float64(x), float64(y))
	})
	return v
}

// Bounds returns the current bounds of the Window.
func (w *Window) Bounds() Rect {
	return w.bounds
}

func (w *Window) setFullscreen(monitor *Monitor) {
	mainthread.Call(func() {
		w.restore.xpos, w.restore.ypos = w.window.GetPos()
		w.restore.width, w.restore.height = w.window.GetSize()

		mode := monitor.monitor.GetVideoMode()

		w.window.SetMonitor(
			monitor.monitor,
			0,
			0,
			mode.Width,
			mode.Height,
			mode.RefreshRate,
		)
	})
}

func (w *Window) setWindowed() {
	mainthread.Call(func() {
		w.window.SetMonitor(
			nil,
			w.restore.xpos,
			w.restore.ypos,
			w.restore.width,
			w.restore.height,
			0,
		)
	})
}

// SetMonitor sets the Window fullscreen on the given Monitor. If the Monitor is nil, the Window
// will be restored to windowed state instead.
//
// The Window will be automatically set to the Monitor's resolution. If you want a different
// resolution, you will need to set it manually with SetBounds method.
func (w *Window) SetMonitor(monitor *Monitor) {
	if w.Monitor() != monitor {
		if monitor != nil {
			w.setFullscreen(monitor)
		} else {
			w.setWindowed()
		}
	}
}

// Monitor returns a monitor the Window is fullscreen on. If the Window is not fullscreen, this
// function returns nil.
func (w *Window) Monitor() *Monitor {
	var monitor *glfw.Monitor
	mainthread.Call(func() {
		monitor = w.window.GetMonitor()
	})
	if monitor == nil {
		return nil
	}
	return &Monitor{
		monitor: monitor,
	}
}

// Focused returns true if the Window has input focus.
func (w *Window) Focused() bool {
	var focused bool
	mainthread.Call(func() {
		focused = w.window.GetAttrib(glfw.Focused) == glfw.True
	})
	return focused
}

// SetVSync sets whether the Window's Update should synchronize with the monitor refresh rate.
func (w *Window) SetVSync(vsync bool) {
	w.vsync = vsync
}

// VSync returns whether the Window is set to synchronize with the monitor refresh rate.
func (w *Window) VSync() bool {
	return w.vsync
}

// SetCursorVisible sets the visibility of the mouse cursor inside the Window client area.
func (w *Window) SetCursorVisible(visible bool) {
	w.cursorVisible = visible
	mainthread.Call(func() {
		if visible {
			w.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		} else {
			w.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
		}
	})
}

// SetCursorDisabled hides the cursor and provides unlimited virtual cursor movement
// make cursor visible using SetCursorVisible
func (w *Window) SetCursorDisabled() {
	w.cursorVisible = false
	mainthread.Call(func() {
		w.window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	})
}

// CursorVisible returns the visibility status of the mouse cursor.
func (w *Window) CursorVisible() bool {
	return w.cursorVisible
}

// Note: must be called inside the main thread.
func (w *Window) begin() {
	if currWin != w {
		w.window.MakeContextCurrent()
		currWin = w
	}
}

// Note: must be called inside the main thread.
func (w *Window) end() {
	// nothing, really
}

// Show makes the window visible, if it was previously hidden. If the window is
// already visible or is in full screen mode, this function does nothing.
func (w *Window) Show() {
	mainthread.Call(func() {
		w.window.Show()
	})
}

// Clipboard returns the contents of the system clipboard.
func (w *Window) Clipboard() string {
	var clipboard string
	mainthread.Call(func() {
		clipboard = w.window.GetClipboardString()
	})
	return clipboard
}

// SetClipboardString sets the system clipboard to the specified UTF-8 encoded string.
func (w *Window) SetClipboard(str string) {
	mainthread.Call(func() {
		w.window.SetClipboardString(str)
	})
}

// SetCanvas sets whether the Window's Update should synchronize with the monitor refresh rate.
func (w *Window) SetCanvas(canvas *Render) {
	w.canvas = canvas
}
