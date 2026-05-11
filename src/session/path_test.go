package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathSessionFinder_FindSessions(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cwd = filepath.Join(cwd, "../../")
	cwd = filepath.Clean(cwd)

	tests := []struct {
		name    string // description of this test case
		finder  PathSessionFinder
		want    []Session
		wantErr bool
	}{
		{
			name: "include extra project",
			finder: PathSessionFinder{
				IncludePaths: []string{cwd + "/test/extra-project"},
			},
			want: []Session{
				{
					Name:           "extra-project",
					WorkingPath:    cwd + "/test/extra-project",
					RepositoryPath: cwd + "/test/extra-project",
					Branch:         "",
					IsActive:       false,
				},
			},
			wantErr: false,
		},
		{
			name: "one depth basic",
			finder: PathSessionFinder{
				SearchPaths: []string{cwd + "/test/one-depth-basic/*"},
			},
			want: []Session{
				{
					Name:           "deep-project-one",
					WorkingPath:    cwd + "/test/one-depth-basic/skip/deep-project-one",
					RepositoryPath: cwd + "/test/one-depth-basic/skip/deep-project-one",
					Branch:         "",
					IsActive:       false,
				},
				{
					Name:           "deep-project-two",
					WorkingPath:    cwd + "/test/one-depth-basic/skip/deep-project-two",
					RepositoryPath: cwd + "/test/one-depth-basic/skip/deep-project-two",
					Branch:         "",
					IsActive:       false,
				},
			},
			wantErr: false,
		},
		{
			name: "zero-depth-basic",
			finder: PathSessionFinder{
				SearchPaths: []string{cwd + "/test/zero-depth-basic"},
			},
			want: []Session{
				{
					Name:           "project-one",
					WorkingPath:    cwd + "/test/zero-depth-basic/project-one",
					RepositoryPath: cwd + "/test/zero-depth-basic/project-one",
					Branch:         "",
					IsActive:       false,
				},
				{
					Name:           "project-three",
					WorkingPath:    cwd + "/test/zero-depth-basic/project-three",
					RepositoryPath: cwd + "/test/zero-depth-basic/project-three",
					Branch:         "",
					IsActive:       false,
				},
				{
					Name:           "project-two",
					WorkingPath:    cwd + "/test/zero-depth-basic/project-two",
					RepositoryPath: cwd + "/test/zero-depth-basic/project-two",
					Branch:         "",
					IsActive:       false,
				},
			},
			wantErr: false,
		},
		{
			name: "zero-depth-worktree",
			finder: PathSessionFinder{
				SearchPaths: []string{cwd + "/test/zero-depth-worktree"},
			},
			want: []Session{
				{
					Name:           "project-one[main]",
					WorkingPath:    cwd + "/test/zero-depth-worktree/project-one/main",
					RepositoryPath: cwd + "/test/zero-depth-worktree/project-one",
					Branch:         "",
					IsActive:       false,
				},
			},
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
