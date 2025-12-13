package database

import (
	"testing"
)

func TestGetDatabaseName(t *testing.T) {
	tests := []struct {
		name string
		conn string
		want string
	}{
		{
			name: "URI format",
			conn: "postgres://user:pass@localhost:5432/mydb?sslmode=disable",
			want: "mydb",
		},
		{
			name: "DSN format",
			conn: "user=user password=pass host=localhost port=5432 dbname=mydb sslmode=disable",
			want: "mydb",
		},
		{
			name: "Invalid format",
			conn: "invalid-string",
			want: "", // Should return empty string, not panic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDatabaseName(tt.conn)
			if got != tt.want {
				t.Errorf("GetDatabaseName() = %v, want %v", got, tt.want)
			}
		})
	}
}
