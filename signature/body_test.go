package signature

import (
	"bytes"
	"crypto"
	"io"
	"testing"
)

func TestHash(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		wantHx  string
		wantErr bool
	}{
		{"non empty body", args{[]byte("{}")}, "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd", false},
		{"empty body", args{nil}, "", false},
		{"empty body", args{[]byte{}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHx, err := Hash(tt.args.body, crypto.SHA512)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHx != tt.wantHx {
				t.Errorf("Hash() = %v, want %v", gotHx, tt.wantHx)
			}
		})
	}
}

func TestHashReader(t *testing.T) {
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantHx  string
		wantErr bool
	}{
		{"non empty", args{bytes.NewBuffer([]byte("{}"))}, "27c74670adb75075fad058d5ceaf7b20c4e7786c83bae8a32f626f9782af34c9a33c2046ef60fd2a7878d378e29fec851806bbd9a67878f3a9f1cda4830763fd", false},
		{"empty", args{bytes.NewBuffer([]byte(""))}, "", false},
		{"empty", args{bytes.NewBuffer(nil)}, "", false},
		{"nil", args{nil}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHx, err := HashReader(tt.args.body, crypto.SHA512)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHx != tt.wantHx {
				t.Errorf("HashReader() = %v, want %v", gotHx, tt.wantHx)
			}
		})
	}
}
