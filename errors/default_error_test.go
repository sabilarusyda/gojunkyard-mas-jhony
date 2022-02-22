package errors

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGetDefaultError(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			"500",
			args{500},
			&Error{
				Status: 500,
				Code:   fmt.Sprint(500),
				Title:  http.StatusText(500),
				Detail: http.StatusText(500),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultError(tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDefaultError() = %v, want %v", got, tt.want)
			}
		})
	}
}
