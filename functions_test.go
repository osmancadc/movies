package main

import (
	"os"
	"testing"
)

func Test_ValidateParameters(t *testing.T) {

	type args struct {
		user User
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "validateParameters001 - Valid email and password",
			args: args{
				user: User{
					Email:    "valid@email.com",
					Password: "ThisPasswordIsStrong!@",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "validateParameters002 - Invalid email",
			args: args{
				user: User{
					Email:    "invalid@email",
					Password: "ThisPasswordIsStrong!@",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "validateParameters003 - Weak password",
			args: args{
				user: User{
					Email:    "valid@email.com",
					Password: "abc123",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateParameters(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasUpper(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "HasUpper001 - String with Uppercase",
			args: args{
				s: "Uppercase",
			},
			want: true,
		},
		{
			name: "HasUpper002 - String without Uppercase",
			args: args{
				s: "lowercase",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasUpper(tt.args.s); got != tt.want {
				t.Errorf("HasUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasLower(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "HasLower001 - String with lowercase",
			args: args{
				s: "lowercase",
			},
			want: true,
		},
		{
			name: "HasLower002 - String without lowercase",
			args: args{
				s: "UPPERCASE",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasLower(tt.args.s); got != tt.want {
				t.Errorf("HasLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GenerateToken001 - Correct execution",
			args: args{
				email: "test@email.com",
			},
			wantErr: false,
		},
		{
			name:    "GenerateToken002 - Without parameters",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() = %v ", got)
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	os.Setenv("token", "eyJUeXAiOiJKV1QiLCJBbGciOiJIUzI1NiIsIkN0eSI6IiJ9.eyJlbWFpbCI6Im9zbWFuY2FkY0Bob3RtYWlsLmNvbSIsImV4cCI6MTYzMzM4NTk2MiwiaWF0IjoxNjMzMzg1OTYxfQ.eiiFgRztNGEItSyXirPvSCcvD0Fdv0iqkSaFbIAVoyw")
	os.Setenv("invalid_token", "eyJUeXAiOiJKV1iLCJBbGciOiJIUzI1NiIsIkN0eSI6IiJ9.eyJlbWFpbCI6Im9zbWFuY2FkY0Bob3RtYWlsLmNvbSIsImV4cCI6MTYzMzM4MjI4NywiaWF0IjoxNjMzMzgxMDg3fQ.R7dZIxtllFgmyoGCzlB2AdvJK3199fhG7NN8MR7pEOQ")

	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidateToken001 - Expired token",
			args: args{
				token: os.Getenv("token"),
			},
			want: false,
		},
		{
			name: "ValidateToken001 - Invalid token",
			args: args{
				token: os.Getenv("invalid_token"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateToken(tt.args.token); got != tt.want {
				t.Errorf("ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
