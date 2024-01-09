package tfmodule

import (
	"testing"
)

func Test_truncatePatchVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "regular version",
			args:    args{version: "1.2.3"},
			want:    "1.2",
			wantErr: false,
		},
		{
			name:    "minor only version",
			args:    args{version: "1.2"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "major only version",
			args:    args{version: "1"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "regular big",
			args:    args{version: "12.34.56"},
			want:    "12.34",
			wantErr: false,
		},
		{
			name:    "bad version",
			args:    args{version: "02.34.56"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := truncatePatchVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("truncatePatchVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("truncatePatchVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
