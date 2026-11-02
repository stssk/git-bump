package versioninfo

import (
	"strings"
	"testing"
)

func TestHighest(t *testing.T) {
	tests := []struct {
		name      string
		tags      []string
		prefix    string
		prefixSet bool
		want      string
		wantFound bool
		wantErr   string
	}{
		{name: "no tags", tags: nil},
		{name: "only non-version tags", tags: []string{"release", "latest"}},
		{name: "single prefix", tags: []string{"v1.2.0", "v1.10.0", "v1.9.9"}, want: "v1.10.0", wantFound: true},
		{name: "no prefix", tags: []string{"1.0.0", "2.0.0-rc.1", "1.5.0"}, want: "2.0.0-rc.1", wantFound: true},
		{name: "release beats pre-release", tags: []string{"v1.0.0-rc.1", "v1.0.0"}, want: "v1.0.0", wantFound: true},
		{
			name:    "planted tag with other prefix is an error",
			tags:    []string{"v1.4.2", "x99.0.0"},
			wantErr: `"v", "x"`,
		},
		{
			name:    "mixed prefixed and unprefixed is an error",
			tags:    []string{"v1.4.2", "99.0.0"},
			wantErr: `"", "v"`,
		},
		{
			name: "explicit prefix ignores other prefixes",
			tags: []string{"v1.4.2", "x99.0.0", "99.0.0"}, prefix: "v", prefixSet: true,
			want: "v1.4.2", wantFound: true,
		},
		{
			name: "explicit empty prefix",
			tags: []string{"v1.4.2", "x99.0.0", "2.0.0"}, prefix: "", prefixSet: true,
			want: "2.0.0", wantFound: true,
		},
		{
			name: "explicit prefix with no matching tags",
			tags: []string{"x99.0.0"}, prefix: "v", prefixSet: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found, err := Highest(tt.tags, tt.prefix, tt.prefixSet)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("err = %v, want error containing %s", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if found != tt.wantFound {
				t.Fatalf("found = %v, want %v", found, tt.wantFound)
			}
			if found && got.String() != tt.want {
				t.Fatalf("got %s, want %s", got.String(), tt.want)
			}
		})
	}
}
