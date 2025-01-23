package gen

import "testing"

func TestInitVersion(t *testing.T) {
	type args struct {
		buildVersion string
		buildDate    string
		buildCommit  string
	}
	tests := []struct {
		name string
		args args
		want *Gen
	}{
		{
			name: "check default",
			args: args{
				buildVersion: "",
				buildDate:    "",
				buildCommit:  "",
			},
			want: &Gen{
				BuildVersion: "N/A",
				BuildDate:    "N/A",
				BuildCommit:  "N/A",
			},
		},
		{
			name: "check set version",
			args: args{
				buildVersion: "1",
				buildDate:    "2",
				buildCommit:  "3",
			},
			want: &Gen{
				BuildVersion: "1",
				BuildDate:    "2",
				BuildCommit:  "3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitVersion(tt.args.buildVersion, tt.args.buildDate, tt.args.buildCommit); err != nil {
				t.Errorf("InitVersion() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}
