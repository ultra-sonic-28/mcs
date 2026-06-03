// Package spectrumrom provides the embedded Spectrum ROM data.
package spectrumrom

import _ "embed"

// Rom48 contains the binary data of the Spectrum 48K ROM.
//
//go:embed 48.rom
var Rom48 []byte
