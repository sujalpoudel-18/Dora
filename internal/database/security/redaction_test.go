package security

import "testing"

func TestRedactSecrets(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "postgres://admin:SuperSecret123@localhost:5432/mydb",
			want:  "postgres://admin:REDACTED@localhost:5432/mydb",
		},
		{
			input: "password=MyPassword123;host=localhost",
			want:  "password=REDACTED;host=localhost",
		},
		{
			input: "mongodb+srv://root:P%40ssw0rd!@cluster0.mongodb.net",
			want:  "mongodb+srv://root:REDACTED@cluster0.mongodb.net",
		},
	}

	for _, tt := range tests {
		got := RedactSecrets(tt.input)
		if got != tt.want {
			t.Errorf("RedactSecrets(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
