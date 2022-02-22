package valkyrie

import (
	"testing"
)

func TestValkyrie_GetProjectID(t *testing.T) {
	type args struct {
		identifier string
	}
	type want struct {
		pid int64
		ok  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "without project",
		},
		{
			name: "with invalid project",
			args: args{
				identifier: "supermantap",
			},
		},
		{
			name: "with valid project",
			args: args{
				identifier: "supersoccertv_analytic",
			},
			want: want{
				pid: 3,
				ok:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Valkyrie{getTestCore()}
			gotPid, gotOk := c.GetProjectID(tt.args.identifier)
			if gotPid != tt.want.pid {
				t.Errorf("Valkyrie.GetProjectID() gotPid = %v, want %v", gotPid, tt.want.pid)
			}
			if gotOk != tt.want.ok {
				t.Errorf("Valkyrie.GetProjectID() gotOk = %v, want %v", gotOk, tt.want.ok)
			}
		})
	}
}
