package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTmuxSessionFinder_FindSessions(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cwd = filepath.Join(cwd, "../../")
	cwd = filepath.Clean(cwd)

	tests := []struct {
		name    string // description of this test case
		finder  TmuxSessionFinder
		want    []Session
		wantErr bool
	}{
		{
			name:   "check for self session",
			finder: TmuxSessionFinder{},
			want: []Session{
				{
					Name:           filepath.Base(cwd),
					WorkingPath:    cwd,
					RepositoryPath: cwd,
					Branch:         "",
					IsActive:       true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.finder.FindSessions()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("FindSessions() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("FindSessions() succeeded unexpectedly")
			}
			if len(got) != len(tt.want) {
				t.Errorf("FindSessions() returned %v sessions, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if got[i].Name != tt.want[i].Name || got[i].WorkingPath != tt.want[i].WorkingPath || got[i].RepositoryPath != tt.want[i].RepositoryPath || got[i].IsActive != tt.want[i].IsActive {
					t.Errorf("FindSessions() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
