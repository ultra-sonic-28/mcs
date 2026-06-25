// Package toolbarassets provides embedded toolbar icon assets.
package toolbarassets

import _ "embed"

// QuitApp contains the PNG image data for the quit-application toolbar button.
//
//go:embed quit-app.png
var QuitApp []byte
