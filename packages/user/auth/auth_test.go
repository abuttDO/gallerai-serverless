package main

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"
)

func Test_validatePassword(t *testing.T) {
	type args struct {
		password     string
		passwordHash string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_validatePassword",
			args: args{
				password:     "password",
				passwordHash: fmt.Sprintf("%x", sha256.Sum256([]byte("password"))),
			},
			want: true,
		},
		{
			name: "test_validatePassword",
			args: args{
				password:     "password",
				passwordHash: fmt.Sprintf("%x", sha256.Sum256([]byte("incorrect"))),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePassword(tt.args.password, tt.args.passwordHash); got != tt.want {
				t.Errorf("validatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUserByEmail(t *testing.T) {
	type args struct {
		search string
	}
	tests := []struct {
		name     string
		args     args
		wantUser *User
	}{
		{
			name: "test_getUserByEmail",
			args: args{
				search: "test",
			},
			wantUser: &User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUser := getUserByEmail(tt.args.search); !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("getUserByEmail() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func Test_issueSignedJwtToken(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name     string
		args     args
		dontWant string
		wantErr  bool
	}{
		{
			name: "test_issueSignedJwtToken",
			args: args{
				user: &User{
					ID:       1,
					Username: "test",
					Email:    "test@test.com",
					Password: "password",
				},
			},
			dontWant: "",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := issueSignedJwtToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("issueSignedJwtToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.dontWant {
				t.Errorf("issueSignedJwtToken() got = %v, want %v", got, tt.dontWant)
			}
		})
	}
}

func TestMainRun(t *testing.T) {
	type args struct {
		in Request
	}
	tests := []struct {
		name    string
		args    args
		want    *Response
		wantErr bool
	}{
		{
			name: "test_Main",
			args: args{
				in: Request{
					UsernameOrEmail: "test",
					Password:        "password",
				},
			},
			want:    makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, fmt.Errorf("invalid password")), fmt.Errorf("invalid password")),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Main(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Main() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Main() got = %v, want %v", got, tt.want)
			}
		})
	}
}
