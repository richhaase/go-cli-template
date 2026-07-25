package config

import "testing"

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		wantDebug bool
	}{
		{
			name:      "defaults",
			env:       map[string]string{"MYCLI_DEBUG": ""},
			wantDebug: false,
		},
		{
			name:      "debug via environment",
			env:       map[string]string{"MYCLI_DEBUG": "true"},
			wantDebug: true,
		},
		{
			name:      "debug requires exact 'true'",
			env:       map[string]string{"MYCLI_DEBUG": "1"},
			wantDebug: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Isolate from any config file in the working directory.
			t.Chdir(t.TempDir())
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if cfg.Debug != tt.wantDebug {
				t.Errorf("cfg.Debug = %v, want %v", cfg.Debug, tt.wantDebug)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Debug {
		t.Error("DefaultConfig().Debug = true, want false")
	}
}
