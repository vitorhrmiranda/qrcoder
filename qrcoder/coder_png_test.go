package qrcoder

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
)

var pngResult = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 1, 244, 0, 0, 1, 244, 1, 3, 0, 0, 0, 241, 24, 77, 201, 0, 0, 0, 6, 80, 76, 84, 69, 255, 255, 255, 0, 0, 0, 85, 194, 211, 126, 0, 0, 1, 144, 73, 68, 65, 84, 120, 218, 236, 220, 49, 114, 163, 48, 20, 128, 97, 101, 82, 164, 228, 8, 62, 74, 142, 70, 142, 230, 163, 248, 8, 148, 91, 100, 242, 118, 64, 2, 132, 215, 94, 192, 75, 181, 124, 127, 103, 123, 62, 151, 140, 244, 36, 59, 73, 146, 36, 73, 146, 36, 105, 91, 31, 177, 232, 187, 127, 175, 141, 136, 91, 106, 226, 59, 189, 45, 63, 13, 158, 231, 249, 35, 253, 79, 253, 98, 242, 17, 93, 74, 189, 159, 63, 125, 231, 121, 158, 255, 195, 71, 1, 31, 227, 35, 231, 154, 46, 179, 207, 95, 247, 139, 231, 121, 126, 139, 31, 23, 51, 60, 207, 243, 59, 252, 176, 126, 249, 135, 231, 15, 207, 243, 252, 30, 127, 183, 127, 202, 53, 227, 243, 107, 247, 254, 139, 231, 249, 179, 248, 187, 249, 237, 91, 196, 87, 250, 124, 125, 254, 203, 243, 252, 105, 252, 163, 242, 249, 207, 235, 241, 60, 127, 42, 223, 150, 253, 79, 244, 203, 153, 126, 255, 147, 215, 47, 165, 203, 176, 255, 225, 121, 158, 63, 212, 151, 249, 75, 191, 228, 233, 134, 61, 210, 79, 253, 252, 154, 230, 47, 221, 243, 251, 115, 60, 207, 159, 220, 71, 92, 83, 51, 238, 146, 250, 47, 155, 54, 83, 253, 203, 105, 163, 197, 243, 60, 255, 96, 254, 114, 45, 111, 76, 35, 151, 46, 189, 207, 164, 121, 126, 254, 195, 243, 252, 89, 253, 80, 253, 200, 25, 250, 156, 215, 47, 49, 93, 166, 229, 121, 158, 63, 222, 215, 45, 238, 223, 166, 250, 110, 75, 240, 60, 207, 175, 157, 63, 103, 31, 121, 152, 187, 58, 127, 225, 121, 254, 188, 254, 209, 253, 183, 106, 49, 19, 229, 199, 200, 81, 223, 141, 227, 121, 158, 127, 254, 251, 159, 91, 222, 240, 228, 249, 237, 229, 111, 235, 23, 158, 231, 249, 227, 252, 87, 25, 249, 78, 87, 110, 215, 239, 175, 240, 60, 207, 47, 230, 183, 229, 252, 168, 93, 253, 255, 37, 158, 231, 207, 235, 239, 246, 79, 109, 57, 140, 230, 121, 158, 223, 251, 255, 143, 121, 254, 82, 157, 255, 220, 134, 251, 43, 219, 230, 191, 60, 207, 243, 27, 189, 36, 73, 146, 36, 73, 210, 127, 218, 239, 0, 0, 0, 255, 255, 167, 233, 155, 124, 68, 70, 18, 39, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}

func TestNewCoderPNG(t *testing.T) {
	type args struct {
		content string
		size    int
	}
	tests := []struct {
		name string
		args args
		want *CoderPNG
	}{
		{
			name: "when the size is within the allowed",
			args: args{
				content: faker.WORD,
				size:    500,
			},
			want: &CoderPNG{
				content: faker.WORD,
				size:    500,
			},
		},
		{
			name: "when the size is larger than allowed",
			args: args{
				content: faker.WORD,
				size:    2049,
			},
			want: &CoderPNG{
				content: faker.WORD,
				size:    MAX_SIZE,
			},
		},
		{
			name: "when the size is smaller than allowed",
			args: args{
				content: faker.WORD,
				size:    1,
			},
			want: &CoderPNG{
				content: faker.WORD,
				size:    MIN_SIZE,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCoderPNG(tt.args.content, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCoderPNG() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoderPNG_Encode(t *testing.T) {
	type fields struct {
		content string
		size    int
	}
	tests := []struct {
		name    string
		fields  fields
		wantRaw []byte
		wantErr bool
	}{
		{
			name:    "when genereted qrcode with success",
			fields:  fields{content: "Foo Bar", size: 500},
			wantRaw: pngResult,
			wantErr: false,
		},
		{
			name:    "when has error",
			fields:  fields{content: "", size: 500},
			wantRaw: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := NewCoderPNG(tt.fields.content, tt.fields.size)
			gotRaw, err := cs.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("CoderPNG.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRaw, tt.wantRaw) {
				ioutil.WriteFile("./error.png", gotRaw, 0644)
				t.Errorf("CoderPNG.Encode()")
			}
		})
	}
}
