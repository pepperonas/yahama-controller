// yamaha_test.go — unit tests for the pure Yamaha-protocol helpers in yamaha.go
package main

import (
	"math"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Volume conversion: DBtoRaw / RawToDB
// ---------------------------------------------------------------------------

func TestDBtoRaw_CommonValues(t *testing.T) {
	cases := []struct {
		db  float64
		raw int
	}{
		{-80.0, -800},
		{-20.0, -200},
		{-10.0, -100},
		{0.0, 0},
		{+16.5, +165},
	}
	for _, tc := range cases {
		got := DBtoRaw(tc.db)
		if got != tc.raw {
			t.Errorf("DBtoRaw(%.1f) = %d, want %d", tc.db, got, tc.raw)
		}
	}
}

func TestDBtoRaw_Rounding(t *testing.T) {
	// -20.05 dB → -200.5 raw → rounds to -201
	// -20.04 dB → -200.4 raw → rounds to -200
	if got := DBtoRaw(-20.05); got != -201 {
		t.Errorf("DBtoRaw(-20.05) = %d, want -201", got)
	}
	if got := DBtoRaw(-20.04); got != -200 {
		t.Errorf("DBtoRaw(-20.04) = %d, want -200", got)
	}
}

func TestDBtoRaw_Clamping(t *testing.T) {
	if got := DBtoRaw(-999.9); got != VolMinRaw {
		t.Errorf("DBtoRaw below min → %d, want %d", got, VolMinRaw)
	}
	if got := DBtoRaw(999.9); got != VolMaxRaw {
		t.Errorf("DBtoRaw above max → %d, want %d", got, VolMaxRaw)
	}
}

func TestRawToDB_CommonValues(t *testing.T) {
	cases := []struct {
		raw int
		db  float64
	}{
		{-800, -80.0},
		{-200, -20.0},
		{0, 0.0},
		{165, 16.5},
		{-255, -25.5},
	}
	for _, tc := range cases {
		got := RawToDB(tc.raw)
		if math.Abs(got-tc.db) > 0.001 {
			t.Errorf("RawToDB(%d) = %.4f, want %.4f", tc.raw, got, tc.db)
		}
	}
}

func TestRawToDB_Clamping(t *testing.T) {
	if got := RawToDB(-9999); got != VolMinDB {
		t.Errorf("RawToDB below min → %.1f, want %.1f", got, VolMinDB)
	}
	if got := RawToDB(9999); got != VolMaxDB {
		t.Errorf("RawToDB above max → %.1f, want %.1f", got, VolMaxDB)
	}
}

// DBtoRaw / RawToDB must be inverse operations (round-trip within ±0.1 dB)
func TestVolume_RoundTrip(t *testing.T) {
	for _, db := range []float64{-80.0, -50.0, -20.0, -10.5, 0.0, 5.0, 16.5} {
		raw := DBtoRaw(db)
		back := RawToDB(raw)
		if math.Abs(back-db) > 0.1 {
			t.Errorf("round-trip %.1f → %d → %.4f: drift > 0.1", db, raw, back)
		}
	}
}

// ---------------------------------------------------------------------------
// DBtoPct / PctToDB
// ---------------------------------------------------------------------------

func TestDBtoPct_Endpoints(t *testing.T) {
	if got := DBtoPct(VolMinDB); math.Abs(got-0) > 0.001 {
		t.Errorf("DBtoPct(min) = %.4f, want 0", got)
	}
	if got := DBtoPct(VolMaxDB); math.Abs(got-100) > 0.001 {
		t.Errorf("DBtoPct(max) = %.4f, want 100", got)
	}
}

func TestDBtoPct_Midpoint(t *testing.T) {
	mid := (VolMinDB + VolMaxDB) / 2
	got := DBtoPct(mid)
	if math.Abs(got-50.0) > 0.001 {
		t.Errorf("DBtoPct(midpoint %.2f) = %.4f, want 50", mid, got)
	}
}

func TestDBtoPct_Clamping(t *testing.T) {
	if got := DBtoPct(-999); got != 0 {
		t.Errorf("DBtoPct below min → %.1f, want 0", got)
	}
	if got := DBtoPct(999); got != 100 {
		t.Errorf("DBtoPct above max → %.1f, want 100", got)
	}
}

func TestPctToDB_Endpoints(t *testing.T) {
	if got := PctToDB(0); math.Abs(got-VolMinDB) > 0.001 {
		t.Errorf("PctToDB(0) = %.4f, want %.1f", got, VolMinDB)
	}
	if got := PctToDB(100); math.Abs(got-VolMaxDB) > 0.001 {
		t.Errorf("PctToDB(100) = %.4f, want %.1f", got, VolMaxDB)
	}
}

func TestPctToDB_Clamping(t *testing.T) {
	if got := PctToDB(-10); math.Abs(got-VolMinDB) > 0.001 {
		t.Errorf("PctToDB(-10) → %.4f, want %.1f", got, VolMinDB)
	}
	if got := PctToDB(110); math.Abs(got-VolMaxDB) > 0.001 {
		t.Errorf("PctToDB(110) → %.4f, want %.1f", got, VolMaxDB)
	}
}

func TestPctDB_RoundTrip(t *testing.T) {
	for _, pct := range []float64{0, 25, 50, 75, 100} {
		db := PctToDB(pct)
		back := DBtoPct(db)
		if math.Abs(back-pct) > 0.001 {
			t.Errorf("round-trip %.0f%% → %.2f dB → %.4f%%: drift", pct, db, back)
		}
	}
}

// ---------------------------------------------------------------------------
// ClampDB
// ---------------------------------------------------------------------------

func TestClampDB(t *testing.T) {
	cases := []struct{ in, want float64 }{
		{-100, VolMinDB},
		{VolMinDB, VolMinDB},
		{0, 0},
		{VolMaxDB, VolMaxDB},
		{100, VolMaxDB},
	}
	for _, tc := range cases {
		got := ClampDB(tc.in)
		if got != tc.want {
			t.Errorf("ClampDB(%.1f) = %.1f, want %.1f", tc.in, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Tone control
// ---------------------------------------------------------------------------

func TestToneDBtoRaw(t *testing.T) {
	cases := []struct {
		db  float64
		raw int
	}{
		{0, 0},
		{6.0, 60},
		{-6.0, -60},
		{3.5, 35},
		{-3.5, -35},
	}
	for _, tc := range cases {
		got := ToneDBtoRaw(tc.db)
		if got != tc.raw {
			t.Errorf("ToneDBtoRaw(%.1f) = %d, want %d", tc.db, got, tc.raw)
		}
	}
}

func TestToneDBtoRaw_Clamping(t *testing.T) {
	if got := ToneDBtoRaw(99); got != 60 {
		t.Errorf("ToneDBtoRaw over max → %d, want 60", got)
	}
	if got := ToneDBtoRaw(-99); got != -60 {
		t.Errorf("ToneDBtoRaw under min → %d, want -60", got)
	}
}

func TestToneRawtoDB(t *testing.T) {
	if got := ToneRawtoDB(35); math.Abs(got-3.5) > 0.001 {
		t.Errorf("ToneRawtoDB(35) = %.4f, want 3.5", got)
	}
	if got := ToneRawtoDB(-60); math.Abs(got-(-6.0)) > 0.001 {
		t.Errorf("ToneRawtoDB(-60) = %.4f, want -6.0", got)
	}
}

// ---------------------------------------------------------------------------
// XML command builders
// ---------------------------------------------------------------------------

func TestPutXML_Structure(t *testing.T) {
	got := PutXML(ZoneMain, "<Power_Control><Power>On</Power></Power_Control>")
	if !strings.HasPrefix(got, `<YAMAHA_AV cmd="PUT">`) {
		t.Errorf("PutXML missing PUT envelope prefix: %s", got)
	}
	if !strings.Contains(got, "<Main_Zone>") {
		t.Errorf("PutXML missing zone tag: %s", got)
	}
	if !strings.HasSuffix(got, "</YAMAHA_AV>") {
		t.Errorf("PutXML missing closing tag: %s", got)
	}
}

func TestGetXML_Structure(t *testing.T) {
	got := GetXML(ZoneMain, "<Basic_Status>GetParam</Basic_Status>")
	if !strings.HasPrefix(got, `<YAMAHA_AV cmd="GET">`) {
		t.Errorf("GetXML missing GET envelope prefix: %s", got)
	}
}

func TestPowerXML_On(t *testing.T) {
	got := PowerXML(ZoneMain, "On")
	if !strings.Contains(got, "<Power>On</Power>") {
		t.Errorf("PowerXML On missing power tag: %s", got)
	}
	if !strings.Contains(got, `cmd="PUT"`) {
		t.Errorf("PowerXML must be a PUT: %s", got)
	}
}

func TestPowerXML_Standby(t *testing.T) {
	got := PowerXML(ZoneMain, "Standby")
	if !strings.Contains(got, "<Power>Standby</Power>") {
		t.Errorf("PowerXML Standby: %s", got)
	}
}

func TestVolumeXML(t *testing.T) {
	got := VolumeXML(ZoneMain, -200)
	if !strings.Contains(got, "<Val>-200</Val>") {
		t.Errorf("VolumeXML missing Val: %s", got)
	}
	if !strings.Contains(got, "<Exp>1</Exp>") {
		t.Errorf("VolumeXML missing Exp: %s", got)
	}
	if !strings.Contains(got, "<Unit>dB</Unit>") {
		t.Errorf("VolumeXML missing Unit: %s", got)
	}
}

func TestMuteXML(t *testing.T) {
	onXML := MuteXML(ZoneMain, "On")
	offXML := MuteXML(ZoneMain, "Off")
	if !strings.Contains(onXML, "<Mute>On</Mute>") {
		t.Errorf("MuteXML On: %s", onXML)
	}
	if !strings.Contains(offXML, "<Mute>Off</Mute>") {
		t.Errorf("MuteXML Off: %s", offXML)
	}
}

func TestInputXML(t *testing.T) {
	got := InputXML(ZoneMain, "HDMI1")
	if !strings.Contains(got, "<Input_Sel>HDMI1</Input_Sel>") {
		t.Errorf("InputXML missing Input_Sel: %s", got)
	}
}

func TestBassXML(t *testing.T) {
	got := BassXML(ZoneMain, 30)
	if !strings.Contains(got, "<Bass>") {
		t.Errorf("BassXML missing Bass tag: %s", got)
	}
	if !strings.Contains(got, "<Val>30</Val>") {
		t.Errorf("BassXML missing Val: %s", got)
	}
}

func TestTrebleXML(t *testing.T) {
	got := TrebleXML(ZoneMain, -20)
	if !strings.Contains(got, "<Treble>") {
		t.Errorf("TrebleXML missing Treble tag: %s", got)
	}
	if !strings.Contains(got, "<Val>-20</Val>") {
		t.Errorf("TrebleXML missing Val: %s", got)
	}
}

func TestBasicStatusGetXML(t *testing.T) {
	got := BasicStatusGetXML(ZoneMain)
	if !strings.Contains(got, "Basic_Status") {
		t.Errorf("BasicStatusGetXML: %s", got)
	}
	if !strings.Contains(got, `cmd="GET"`) {
		t.Errorf("BasicStatusGetXML must be GET: %s", got)
	}
}

func TestXML_Zone2(t *testing.T) {
	got := PowerXML(ZoneTwo, "On")
	if !strings.Contains(got, "<Zone_2>") {
		t.Errorf("Zone2 XML missing Zone_2 tag: %s", got)
	}
}

// ---------------------------------------------------------------------------
// VolumeXML round-trip: dB → raw → XML contains expected Val
// ---------------------------------------------------------------------------

func TestVolumeXML_RoundTrip(t *testing.T) {
	cases := []struct {
		db  float64
		val string
	}{
		{-20.0, "-200"},
		{0.0, "0"},
		{-80.0, "-800"},
		{16.5, "165"},
	}
	for _, tc := range cases {
		raw := DBtoRaw(tc.db)
		xml := VolumeXML(ZoneMain, raw)
		if !strings.Contains(xml, "<Val>"+tc.val+"</Val>") {
			t.Errorf("VolumeXML round-trip %.1f dB: want Val=%s in %s", tc.db, tc.val, xml)
		}
	}
}

// ---------------------------------------------------------------------------
// Input source validation
// ---------------------------------------------------------------------------

func TestIsKnownInput_Valid(t *testing.T) {
	valid := []string{"HDMI1", "HDMI2", "HDMI3", "HDMI4", "AV1", "AV2", "AUDIO1", "AUDIO2",
		"AirPlay", "SERVER", "NET RADIO", "Spotify", "Bluetooth", "USB", "TUNER", "V-AUX"}
	for _, inp := range valid {
		if !IsKnownInput(inp) {
			t.Errorf("IsKnownInput(%q) = false, want true", inp)
		}
	}
}

func TestIsKnownInput_CaseInsensitive(t *testing.T) {
	cases := []string{"hdmi1", "HDMI1", "Hdmi1", "airplay", "AIRPLAY", "spotify", "SPOTIFY"}
	for _, s := range cases {
		if !IsKnownInput(s) {
			t.Errorf("IsKnownInput(%q) should be case-insensitive true, got false", s)
		}
	}
}

func TestIsKnownInput_Invalid(t *testing.T) {
	invalid := []string{"", "HDMI5", "OPTICAL", "COAX", "PHONO", "foobar", "HDMI 1"}
	for _, inp := range invalid {
		if IsKnownInput(inp) {
			t.Errorf("IsKnownInput(%q) = true, want false", inp)
		}
	}
}

// ---------------------------------------------------------------------------
// Proxy path stripping
// ---------------------------------------------------------------------------

func TestStripProxyPrefix(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/api/receiver", "/"},
		{"/api/receiver/", "/"},
		{"/api/receiver/YamahaRemoteControl/ctrl", "/YamahaRemoteControl/ctrl"},
		{"/api/receiver/YamahaRemoteControl", "/YamahaRemoteControl"},
		// path that doesn't start with the prefix is returned unchanged
		{"/other/path", "/other/path"},
	}
	for _, tc := range cases {
		got := StripProxyPrefix(tc.path)
		if got != tc.want {
			t.Errorf("StripProxyPrefix(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}
