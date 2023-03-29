package options

import (
	"os"
	"reflect"
	"testing"
)

func TestOptions_LoadFromEnv(t *testing.T) {
	type fields struct {
		GWTkeyID            string
		GWTserviceAccountID string
		GWTkeyFile          string
	}
	tests := []struct {
		name   string
		fields fields
		want   *Options
	}{
		{
			name:   "1",
			fields: fields{},
			want: &Options{
				GWTkeyID:            "GWTkeyID",
				GWTserviceAccountID: "GWTserviceAccountID",
				GWTkeyFile:          "GWTkeyFile",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("GWT_KEY_ID", "GWTkeyID")
			os.Setenv("GWT_SERVICE_ACCOUNT_ID", "GWTserviceAccountID")
			os.Setenv("GWT_KEY_FILE", "GWTkeyFile")
			o := &Options{
				GWTkeyID:            tt.fields.GWTkeyID,
				GWTserviceAccountID: tt.fields.GWTserviceAccountID,
				GWTkeyFile:          tt.fields.GWTkeyFile,
			}
			if got := o.loadFromEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.LoadFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
