package mssql

import "testing"

func TestUsesAzureADAuth(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want bool
	}{
		{
			name: "fedauth",
			dsn:  "Server=s.database.windows.net;Database=db;fedauth=ActiveDirectoryDefault;",
			want: true,
		},
		{
			name: "csharp authentication",
			dsn:  "Server=s.database.windows.net;Database=db;Authentication=Active Directory Interactive;",
			want: true,
		},
		{
			name: "sspi",
			dsn:  "Server=cmdp_db_uat;Database=CMDP;Integrated Security=SSPI;",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := usesAzureADAuth(tt.dsn); got != tt.want {
				t.Fatalf("usesAzureADAuth() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestNormalizeAzureADDSNConvertsCSharpInteractiveToDefaultFedAuth(t *testing.T) {
	got := normalizeAzureADDSN("Server=s.database.windows.net;Authentication=Active Directory Interactive;Database=db;")
	want := "Server=s.database.windows.net;Database=db;fedauth=ActiveDirectoryDefault"
	if got != want {
		t.Fatalf("normalized dsn = %q, want %q", got, want)
	}
}

func TestNormalizeAzureADDSNKeepsExistingFedAuth(t *testing.T) {
	got := normalizeAzureADDSN("Server=s.database.windows.net;fedauth=ActiveDirectoryAzCli;Database=db;")
	want := "Server=s.database.windows.net;fedauth=ActiveDirectoryAzCli;Database=db"
	if got != want {
		t.Fatalf("normalized dsn = %q, want %q", got, want)
	}
}
