package workrecord

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The contract rows — which beads Gated covers and what ValidateOnClose calls a
// violation — are pinned by the two planes that run them:
// cmd/gc/work_record_gate_test.go against the bd-argv adapter and
// internal/api/work_record_close_gate_api_test.go against the HTTP handlers.
// CommitReachableOnBranch is the one piece neither plane can pin honestly,
// because both inject a fake oracle to stay off disk — so it is pinned here,
// against a real repository.

// The git process is declared Medium ownership in test/test-resources.toml, so
// it is deliberately built inside this runnable rather than in a package-level
// helper: the resource census only recognizes an owner it can name, and a
// process a shared helper spawns belongs to no test in particular.
func TestCommitReachableOnBranch(t *testing.T) {
	runGit := func(dir string, args ...string) string {
		t.Helper()
		out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return string(out)
	}

	// A repo with one commit on main and one commit on a side branch main never
	// saw: reachability is only a real question when some commit is unreachable.
	repoDir := t.TempDir()
	runGit(repoDir, "init", "--initial-branch=main")
	runGit(repoDir, "config", "user.name", "Gas City Test")
	runGit(repoDir, "config", "user.email", "gc-test@test.local")
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(repoDir, name), []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
		runGit(repoDir, "add", name)
	}
	write("shipped.txt", "on main\n")
	runGit(repoDir, "commit", "-m", "test: land the artifact")
	onMain := strings.TrimSpace(runGit(repoDir, "rev-parse", "HEAD"))

	runGit(repoDir, "checkout", "-b", "side")
	write("side.txt", "never merged\n")
	runGit(repoDir, "commit", "-m", "test: a commit main never saw")
	offMain := strings.TrimSpace(runGit(repoDir, "rev-parse", "HEAD"))
	runGit(repoDir, "checkout", "main")

	notARepo := t.TempDir()

	tests := []struct {
		name           string
		dir            string
		commit, branch string
		want           bool
	}{
		{"commit on the branch is reachable", repoDir, onMain, "main", true},
		{"commit on another branch is not reachable", repoDir, offMain, "main", false},
		{"commit reachable on its own branch", repoDir, offMain, "side", true},
		{"unknown commit is not reachable", repoDir, "0000000000000000000000000000000000000000", "main", false},
		{"unknown branch is not reachable", repoDir, onMain, "no-such-branch", false},
		{"a directory that is not a repo is not reachable", notARepo, onMain, "main", false},
		{"an empty repo dir is not reachable", "", onMain, "main", false},
		{"an empty commit is not reachable", repoDir, "", "main", false},
		{"an empty branch is not reachable", repoDir, onMain, "", false},
		{"a flag-shaped commit is rejected", repoDir, "--all", "main", false},
		{"a flag-shaped branch is rejected", repoDir, onMain, "--all", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := CommitReachableOnBranch(tc.dir, tc.commit, tc.branch); got != tc.want {
				t.Fatalf("CommitReachableOnBranch(%q, %q, %q) = %v, want %v", tc.dir, tc.commit, tc.branch, got, tc.want)
			}
		})
	}
}
