package qrcoder

import (
	"io/ioutil"
	"testing"
)

func TestCoderPDF_Encode(t *testing.T) {
	type fields struct {
		content string
		size    int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "when create pdf with success",
			fields:  fields{content: "Foo Bar", size: 500},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := NewCoderPDF(tt.fields.content, tt.fields.size)
			gotRaw, err := cs.Encode()
			if (err != nil) != tt.wantErr {
				ioutil.WriteFile("./error.pdf", gotRaw, 0644)
				t.Errorf("CoderPDF.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
