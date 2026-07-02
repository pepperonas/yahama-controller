// yamaha.go — pure helper functions for the Yamaha RX-V577 protocol.
//
// All functions here are side-effect-free and operate only on their
// arguments, making them straightforwardly unit-testable without a
// receiver or network connection.
package main

import (
	"fmt"
	"math"
	"strings"
)

// ---------------------------------------------------------------------------
// Volume conversion
// ---------------------------------------------------------------------------

// The Yamaha RX-V577 internal volume unit is tenths of a dB.
// Val=-800 → -80.0 dB (minimum / silence)
// Val= 165 → +16.5 dB (maximum)
const (
	VolMinDB  = -80.0
	VolMaxDB  = 16.5
	VolMinRaw = -800
	VolMaxRaw = 165
)

// DBtoRaw converts a dB value to the integer raw value sent in the XML <Val>.
// The raw value is dB × 10 (e.g. -20.5 dB → -205).
// The result is clamped to the valid receiver range.
func DBtoRaw(db float64) int {
	if db < VolMinDB {
		db = VolMinDB
	}
	if db > VolMaxDB {
		db = VolMaxDB
	}
	return int(math.Round(db * 10))
}

// RawToDB converts the raw receiver integer back to dB.
func RawToDB(raw int) float64 {
	if raw < VolMinRaw {
		raw = VolMinRaw
	}
	if raw > VolMaxRaw {
		raw = VolMaxRaw
	}
	return float64(raw) / 10.0
}

// DBtoPct maps a dB value to a [0, 100] percentage for slider/UI use.
// -80.0 dB → 0 %, +16.5 dB → 100 %.
func DBtoPct(db float64) float64 {
	span := VolMaxDB - VolMinDB
	pct := (db - VolMinDB) / span * 100.0
	if pct < 0 {
		return 0
	}
	if pct > 100 {
		return 100
	}
	return pct
}

// PctToDB maps a [0, 100] percentage back to dB, clamped.
func PctToDB(pct float64) float64 {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	span := VolMaxDB - VolMinDB
	return VolMinDB + (pct/100.0)*span
}

// ClampDB clamps a dB value to the valid receiver range [-80.0, +16.5].
func ClampDB(db float64) float64 {
	if db < VolMinDB {
		return VolMinDB
	}
	if db > VolMaxDB {
		return VolMaxDB
	}
	return db
}

// ---------------------------------------------------------------------------
// Tone control conversion  (-6.0 dB … +6.0 dB, raw = val × 10)
// ---------------------------------------------------------------------------

// ToneDBtoRaw converts a bass/treble dB value (−6.0 … +6.0) to raw.
func ToneDBtoRaw(db float64) int {
	if db < -6.0 {
		db = -6.0
	}
	if db > 6.0 {
		db = 6.0
	}
	return int(math.Round(db * 10))
}

// ToneRawtoDB converts a raw tone value back to dB.
func ToneRawtoDB(raw int) float64 {
	return float64(raw) / 10.0
}

// ---------------------------------------------------------------------------
// XML command builders (pure — no I/O)
// ---------------------------------------------------------------------------

// Zone names accepted by the RX-V577.
type Zone string

const (
	ZoneMain Zone = "Main_Zone"
	ZoneTwo  Zone = "Zone_2"
)

// PutXML wraps the given inner XML in the standard Yamaha PUT envelope.
func PutXML(zone Zone, inner string) string {
	return fmt.Sprintf(`<YAMAHA_AV cmd="PUT"><%s>%s</%s></YAMAHA_AV>`, zone, inner, zone)
}

// GetXML wraps the given inner XML in the standard Yamaha GET envelope.
func GetXML(zone Zone, inner string) string {
	return fmt.Sprintf(`<YAMAHA_AV cmd="GET"><%s>%s</%s></YAMAHA_AV>`, zone, inner, zone)
}

// PowerXML returns a PUT command to set receiver power (value: "On" or "Standby").
func PowerXML(zone Zone, state string) string {
	return PutXML(zone, fmt.Sprintf(`<Power_Control><Power>%s</Power></Power_Control>`, state))
}

// VolumeXML returns a PUT command to set volume via a raw integer value.
func VolumeXML(zone Zone, rawVal int) string {
	return PutXML(zone,
		fmt.Sprintf(`<Volume><Lvl><Val>%d</Val><Exp>1</Exp><Unit>dB</Unit></Lvl></Volume>`, rawVal))
}

// MuteXML returns a PUT command to set mute state ("On" or "Off").
func MuteXML(zone Zone, state string) string {
	return PutXML(zone, fmt.Sprintf(`<Volume><Mute>%s</Mute></Volume>`, state))
}

// InputXML returns a PUT command to select an input source.
func InputXML(zone Zone, inputSel string) string {
	return PutXML(zone, fmt.Sprintf(`<Input><Input_Sel>%s</Input_Sel></Input>`, inputSel))
}

// BasicStatusGetXML returns the GET command for Basic_Status.
func BasicStatusGetXML(zone Zone) string {
	return GetXML(zone, `<Basic_Status>GetParam</Basic_Status>`)
}

// BassXML returns a PUT command to set bass tone.
func BassXML(zone Zone, rawVal int) string {
	return PutXML(zone,
		fmt.Sprintf(`<Sound_Video><Tone><Bass><Val>%d</Val><Exp>1</Exp><Unit>dB</Unit></Bass></Tone></Sound_Video>`, rawVal))
}

// TrebleXML returns a PUT command to set treble tone.
func TrebleXML(zone Zone, rawVal int) string {
	return PutXML(zone,
		fmt.Sprintf(`<Sound_Video><Tone><Treble><Val>%d</Val><Exp>1</Exp><Unit>dB</Unit></Treble></Tone></Sound_Video>`, rawVal))
}

// ---------------------------------------------------------------------------
// Input source name mapping
// ---------------------------------------------------------------------------

// KnownInputs is the set of valid input selector strings for the RX-V577.
var KnownInputs = []string{
	"HDMI1", "HDMI2", "HDMI3", "HDMI4",
	"AV1", "AV2",
	"AUDIO1", "AUDIO2",
	"AirPlay",
	"SERVER",
	"NET RADIO",
	"Spotify",
	"Bluetooth",
	"USB",
	"TUNER",
	"V-AUX",
}

// IsKnownInput reports whether s is a recognised input selector.
func IsKnownInput(s string) bool {
	for _, k := range KnownInputs {
		if strings.EqualFold(k, s) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Proxy path stripping (mirrors the Director logic in main.go)
// ---------------------------------------------------------------------------

// StripProxyPrefix strips "/api/receiver" from the beginning of path,
// returning "/" when the result would be empty.
func StripProxyPrefix(path string) string {
	stripped := strings.TrimPrefix(path, "/api/receiver")
	if stripped == "" {
		return "/"
	}
	return stripped
}
