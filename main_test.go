package main

import "testing"

func TestLoadTursoURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		env        map[string]string
		wantMain   string
		wantAux    string
		wantErrMsg string
	}{
		{
			name:       "missing main url",
			env:        map[string]string{},
			wantErrMsg: "URL_LIBSQL_TURSO não definida",
		},
		{
			name: "missing aux url",
			env: map[string]string{
				"URL_LIBSQL_TURSO": "libsql://main",
			},
			wantErrMsg: "URL_LIBSQL_TURSO_AUX não definida",
		},
		{
			name: "both urls provided",
			env: map[string]string{
				"URL_LIBSQL_TURSO":     "libsql://main",
				"URL_LIBSQL_TURSO_AUX": "libsql://aux",
			},
			wantMain: "libsql://main",
			wantAux:  "libsql://aux",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mainURL, auxURL, err := loadTursoURLs(func(key string) string {
				return tt.env[key]
			})

			if tt.wantErrMsg != "" {
				if err == nil || err.Error() != tt.wantErrMsg {
					t.Fatalf("expected error %q, got %v", tt.wantErrMsg, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if mainURL != tt.wantMain {
				t.Fatalf("expected main url %q, got %q", tt.wantMain, mainURL)
			}
			if auxURL != tt.wantAux {
				t.Fatalf("expected aux url %q, got %q", tt.wantAux, auxURL)
			}
		})
	}
}
