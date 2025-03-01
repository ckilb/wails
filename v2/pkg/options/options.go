package options

import (
	"context"
	"html"
	"io/fs"
	"log"
	"net/http"
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/wailsapp/wails/v2/pkg/menu"

	"github.com/imdario/mergo"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type WindowStartState int

const (
	Normal     WindowStartState = 0
	Maximised  WindowStartState = 1
	Minimised  WindowStartState = 2
	Fullscreen WindowStartState = 3
)

type Experimental struct {
}

// App contains options for creating the App
type App struct {
	Title             string
	Width             int
	Height            int
	DisableResize     bool
	Fullscreen        bool
	Frameless         bool
	MinWidth          int
	MinHeight         int
	MaxWidth          int
	MaxHeight         int
	StartHidden       bool
	HideWindowOnClose bool
	AlwaysOnTop       bool
	// BackgroundColour is the background colour of the window
	// You can use the options.NewRGB and options.NewRGBA functions to create a new colour
	BackgroundColour *RGBA
	// RGBA is deprecated. Please use BackgroundColour
	RGBA               *RGBA
	Assets             fs.FS
	AssetsHandler      http.Handler
	Menu               *menu.Menu
	Logger             logger.Logger `json:"-"`
	LogLevel           logger.LogLevel
	LogLevelProduction logger.LogLevel
	OnStartup          func(ctx context.Context)                `json:"-"`
	OnDomReady         func(ctx context.Context)                `json:"-"`
	OnShutdown         func(ctx context.Context)                `json:"-"`
	OnBeforeClose      func(ctx context.Context) (prevent bool) `json:"-"`
	Bind               []interface{}
	WindowStartState   WindowStartState

	// CSS property to test for draggable elements. Default "--wails-draggable"
	CSSDragProperty string

	// The CSS Value that the CSSDragProperty must have to be draggable, EG: "drag"
	CSSDragValue string

	Windows *windows.Options
	Mac     *mac.Options
	Linux   *linux.Options

	// Experimental options
	Experimental *Experimental
}

type RGBA struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

// NewRGBA creates a new RGBA struct with the given values
func NewRGBA(r, g, b, a uint8) *RGBA {
	return &RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

// NewRGB creates a new RGBA struct with the given values and Alpha set to 255
func NewRGB(r, g, b uint8) *RGBA {
	return &RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
}

// MergeDefaults will set the minimum default values for an application
func MergeDefaults(appoptions *App) {

	// Do default merge
	err := mergo.Merge(appoptions, Default)
	if err != nil {
		log.Fatal(err)
	}

	// Default colour. Doesn't work well with mergo
	if appoptions.BackgroundColour == nil {
		appoptions.BackgroundColour = &RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		}
	}

	// Ensure max and min are valid
	processMinMaxConstraints(appoptions)

	// Default menus
	processMenus(appoptions)

	// Process Drag Options
	processDragOptions(appoptions)
}

func processMenus(appoptions *App) {
	switch runtime.GOOS {
	case "darwin":
		if appoptions.Menu == nil {
			appoptions.Menu = defaultMacMenu
		}
	}
}

func processMinMaxConstraints(appoptions *App) {
	if appoptions.MinWidth > 0 && appoptions.MaxWidth > 0 {
		if appoptions.MinWidth > appoptions.MaxWidth {
			appoptions.MinWidth = appoptions.MaxWidth
		}
	}
	if appoptions.MinHeight > 0 && appoptions.MaxHeight > 0 {
		if appoptions.MinHeight > appoptions.MaxHeight {
			appoptions.MinHeight = appoptions.MaxHeight
		}
	}
	// Ensure width and height are limited if max/min is set
	if appoptions.Width < appoptions.MinWidth {
		appoptions.Width = appoptions.MinWidth
	}
	if appoptions.MaxWidth > 0 && appoptions.Width > appoptions.MaxWidth {
		appoptions.Width = appoptions.MaxWidth
	}
	if appoptions.Height < appoptions.MinHeight {
		appoptions.Height = appoptions.MinHeight
	}
	if appoptions.MaxHeight > 0 && appoptions.Height > appoptions.MaxHeight {
		appoptions.Height = appoptions.MaxHeight
	}
}

func processDragOptions(appoptions *App) {
	appoptions.CSSDragProperty = html.EscapeString(appoptions.CSSDragProperty)
	appoptions.CSSDragValue = html.EscapeString(appoptions.CSSDragValue)
}
