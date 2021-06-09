package qrcoder

import (
	"io/ioutil"
	"testing"

	"github.com/bxcodec/faker/v3"
)

func TestCreateQRCode(t *testing.T) {
	type args struct {
		extension string
		size      int
	}
	tests := []struct {
		name         string
		args         args
		wantMimetype string
		wantErr      bool
	}{
		{
			name: "when create svg",
			args: args{
				extension: SVG,
				size:      500,
			},
			wantMimetype: MIME_SVG,
			wantErr:      false,
		},
		{
			name: "when create png",
			args: args{
				extension: PNG,
				size:      500,
			},
			wantMimetype: MIME_PNG,
			wantErr:      false,
		},
		{
			name: "when create pdf",
			args: args{
				extension: PDF,
				size:      500,
			},
			wantMimetype: MIME_PDF,
			wantErr:      false,
		},
		{
			name: "when create unknown extension",
			args: args{
				extension: "gif",
				size:      500,
			},
			wantMimetype: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, gotMimetype, err := CreateQRCode(faker.Word(), tt.args.extension, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateQRCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			ioutil.WriteFile("test."+tt.args.extension, raw, 0644)
			if gotMimetype != tt.wantMimetype {
				t.Errorf("CreateQRCode() gotMimetype = %v, want %v", gotMimetype, tt.wantMimetype)
			}
		})
	}
}
