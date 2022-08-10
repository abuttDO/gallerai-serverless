package main

import (
	"testing"
)

func Test_getImagesByUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    []GalleryImage
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				username: "test",
			},
			want:    []GalleryImage{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getImagesByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("getImagesByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("getImagesByUsername() got = %v, want %v", got, tt.want)
			}
		})
	}
}
