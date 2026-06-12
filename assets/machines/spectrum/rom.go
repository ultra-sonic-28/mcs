// Package spectrumrom provides the embedded Spectrum ROM data.
package spectrumrom

import _ "embed"

// Rom48 contains the binary data of the Spectrum 48K ROM.
//
//go:embed 48.rom
var Rom48 []byte

// Rom128_0 contains the binary data of the Spectrum 128K ROM 0 (128K Editor/Menu).
//
//go:embed 128-0.rom
var Rom128_0 []byte

// Rom128_1 contains the binary data of the Spectrum 128K ROM 1 (48K BASIC).
//
//go:embed 128-1.rom
var Rom128_1 []byte
