package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAccessUser(t *testing.T) {
	rw := httptest.NewRecorder()
	request := http.Request{
		Method: "POST",
		Body:   ioutil.NopCloser(strings.NewReader(`{ "name":"Osman Beltran", "email":"osmancadc@hotmail.com", "password":"Abc123456@@!" }`)),
	}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "AccessUser001 - Correct execution",
			args: args{
				w: rw,
				r: &request,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AccessUser(tt.args.w, tt.args.r)
		})
	}
}

func TestInsertUser(t *testing.T) {
	type args struct {
		user User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "InsertUser001 - Correct Insert",
			args: args{
				user: User{
					Name:     "Test user",
					Email:    "test@test.com",
					Password: "abc123",
				},
			},
			wantErr: false,
		},
		{
			name: "InsertUser001 - Empty User Data",
			args: args{
				user: User{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifyUser(t *testing.T) {
	type args struct {
		verificationType int
		parameters       []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "VerifyUser001 - (Register) Existing email",
			args: args{
				verificationType: 1,
				parameters:       []string{"osmancadc@hotmail.com"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "VerifyUser001 - (Register) Non existing email",
			args: args{
				verificationType: 1,
				parameters:       []string{"some@email.us"},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "VerifyUser003 - (Login) Correct email and password",
			args: args{
				verificationType: 2,
				parameters:       []string{"osmancadc@hotmail.com", "Abc123456@@!"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "VerifyUser003 - (Login) Correct email wrong password",
			args: args{
				verificationType: 2,
				parameters:       []string{"osmancadc@hotmail.com", "wrongpass"},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyUser(tt.args.verificationType, tt.args.parameters...)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
