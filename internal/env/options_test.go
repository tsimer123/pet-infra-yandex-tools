package env

import (
	"os"
	"reflect"
	"testing"
)

func TestOptions_LoadFromEnv(t *testing.T) {
	type fields struct {
		YaJWTkeyID            string
		YaJWTserviceAccountID string
		YaJWTkeyFile          string
		GithubToken           string
		GithubOwner           string
		GithubRepo            string
		GithubSecretName      string
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
				YaJWTkeyID:            "YaJWTkeyID",
				YaJWTserviceAccountID: "YaJWTserviceAccountID",
				YaJWTkeyFile:          "YaJWTkeyFile",
				GithubToken:           "GithubToken",
				GithubOwner:           "GithubOwner",
				GithubRepo:            "GithubRepo",
				GithubSecretName:      "GithubSecretName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("YA_JWT_KEY_ID", "YaJWTkeyID")
			os.Setenv("YA_JWT_SERVICE_ACCOUNT_ID", "YaJWTserviceAccountID")
			os.Setenv("YA_JWT_KEY_FILE", "YaJWTkeyFile")
			os.Setenv("GITHUB_TOKEN", "GithubToken")
			os.Setenv("GITHUB_OWNER", "GithubOwner")
			os.Setenv("GITHUB_REPO", "GithubRepo")
			os.Setenv("GITHUB_SECRET_NAME", "GithubSecretName")
			o := &Options{
				YaJWTkeyID:            tt.fields.YaJWTkeyID,
				YaJWTserviceAccountID: tt.fields.YaJWTserviceAccountID,
				YaJWTkeyFile:          tt.fields.YaJWTkeyFile,
				GithubToken:           tt.fields.GithubToken,
				GithubOwner:           tt.fields.GithubOwner,
				GithubRepo:            tt.fields.GithubRepo,
				GithubSecretName:      tt.fields.GithubSecretName,
			}
			if got := o.loadFromEnv(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options.LoadFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
