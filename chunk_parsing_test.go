package aiff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecoder_parseBascChunk(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		hasInfo bool
		info    AppleMetadata
	}{
		{"no apple metadata", "fixtures/kick.aif", false, AppleMetadata{}},
		{"full data", "fixtures/ring.aif", true, AppleMetadata{Beats: 3, Note: 48, Scale: 2, Numerator: 4, Denominator: 4, IsLooping: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, _ := filepath.Abs(tt.input)
			t.Log(path)
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			d := NewDecoder(f)
			if err := d.Drain(); err != nil {
				t.Fatalf("draining %s failed - %s\n", path, err)
			}
			if tt.hasInfo != d.HasAppleInfo {
				t.Fatalf("%s was expected to have Apple info set to %T but was %T", path, tt.hasInfo, d.HasAppleInfo)
			}
			if d.HasAppleInfo {
				if tt.info.Beats != d.AppleInfo.Beats {
					t.Fatalf("expected to have %d beats but got %d", tt.info.Beats, d.AppleInfo.Beats)
				}
				if tt.info.Note != d.AppleInfo.Note {
					t.Fatalf("expected to have root note set to %d but got %d", tt.info.Note, d.AppleInfo.Note)
				}
				if tt.info.Scale != d.AppleInfo.Scale {
					t.Fatalf("expected to have its scale set to %d but got %d", tt.info.Scale, d.AppleInfo.Scale)
				}
				if tt.info.Numerator != d.AppleInfo.Numerator {
					t.Fatalf("expected to have its Numerator set to %d but got %d", tt.info.Numerator, d.AppleInfo.Numerator)
				}
				if tt.info.Denominator != d.AppleInfo.Denominator {
					t.Fatalf("expected to have its denominator set to %d but got %d", tt.info.Denominator, d.AppleInfo.Denominator)
				}
				if tt.info.IsLooping != d.AppleInfo.IsLooping {
					t.Fatalf("expected to have its looping set to %T but got %T", tt.info.IsLooping, d.AppleInfo.IsLooping)
				}
			}
		})
	}
}
