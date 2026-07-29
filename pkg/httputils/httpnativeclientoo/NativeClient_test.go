package httpnativeclientoo

import (
	"context"
	"reflect"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
)

func TestNativeClient_SendRequest(t *testing.T) {
	type fields struct {
		baseUrl   string
		port      int
		basicAuth *httpoptions.BasicAuth
	}
	type args struct {
		ctx            context.Context
		requestOptions *httpoptions.RequestOptions
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse httputilsinterfaces.Response
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &NativeClient{
				baseUrl:   tt.fields.baseUrl,
				port:      tt.fields.port,
				basicAuth: tt.fields.basicAuth,
			}
			gotResponse, err := c.SendRequest(tt.args.ctx, tt.args.requestOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("NativeClient.SendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("NativeClient.SendRequest() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
