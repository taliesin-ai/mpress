package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SafeOutputPath resolves an output directory and rejects paths that can erase
// project inputs or persistent state. Missing output directories are resolved
// through their nearest existing ancestor, including any symbolic links. The
// final component is retained so removal unlinks an output symlink rather than
// deleting the directory it points to.
func (c Config) SafeOutputPath(project, output string) (string, error) {
	root, err := canonicalPath(project)
	if err != nil {
		return "", fmt.Errorf("unsafe output directory %s: resolve project: %w", output, err)
	}
	resolved, err := canonicalPath(output)
	if err != nil {
		return "", fmt.Errorf("unsafe output directory %s: %w", output, err)
	}
	if !pathBelow(root, resolved) {
		return "", fmt.Errorf("unsafe output directory %s: must stay inside the project", output)
	}
	protected := []string{
		c.ContentPath(project), c.StaticPath(project), filepath.Join(project, Filename),
		filepath.Join(project, ".git"), c.ArtifactsPath(project),
		filepath.Join(project, c.Translation.StateDir),
		filepath.Join(project, ".mpress", "cache"), filepath.Join(project, ".mpress", "backups"),
		filepath.Join(project, ".mpress", "onboarding"),
	}
	if c.Build.CustomCSS != "" {
		protected = append(protected, filepath.Join(project, c.Build.CustomCSS))
	}
	for _, path := range protected {
		canonical, err := canonicalPath(path)
		if err != nil {
			return "", fmt.Errorf("unsafe output directory %s: resolve protected path %s: %w", output, path, err)
		}
		if pathsOverlap(resolved, canonical) {
			return "", fmt.Errorf("unsafe output directory %s: overlaps protected path %s", output, path)
		}
	}
	state, err := canonicalPath(filepath.Join(project, ".mpress"))
	if err != nil {
		return "", fmt.Errorf("unsafe output directory %s: resolve project state: %w", output, err)
	}
	if pathsOverlap(resolved, state) && !temporaryOutput(state, resolved) {
		return "", fmt.Errorf("unsafe output directory %s: overlaps project state", output)
	}
	parent, err := canonicalPath(filepath.Dir(output))
	if err != nil {
		return "", fmt.Errorf("unsafe output directory %s: resolve parent: %w", output, err)
	}
	return filepath.Join(parent, filepath.Base(output)), nil
}

// RemoveOutput applies the shared build/clean policy immediately before
// removal. Root.RemoveAll also confines deletion if a parent changes to an
// external symlink after validation; it never removes the project root itself.
func (c Config) RemoveOutput(project, output string) error {
	resolved, err := c.SafeOutputPath(project, output)
	if err != nil {
		return err
	}
	root, err := canonicalPath(project)
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(root, resolved)
	if err != nil || relative == "." || !filepath.IsLocal(relative) {
		return fmt.Errorf("unsafe output directory %s", output)
	}
	directory, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.RemoveAll(relative)
}

func canonicalPath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	var missing []string
	for {
		_, err := os.Lstat(absolute)
		if err == nil {
			resolved, err := filepath.EvalSymlinks(absolute)
			if err != nil {
				return "", err
			}
			for i := len(missing) - 1; i >= 0; i-- {
				resolved = filepath.Join(resolved, missing[i])
			}
			return resolved, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(absolute)
		if parent == absolute {
			return "", err
		}
		missing = append(missing, filepath.Base(absolute))
		absolute = parent
	}
}

func pathBelow(parent, child string) bool {
	relative, err := filepath.Rel(parent, child)
	return err == nil && relative != "." && filepath.IsLocal(relative)
}

func pathsOverlap(first, second string) bool {
	forward, err := filepath.Rel(first, second)
	if err == nil && filepath.IsLocal(forward) {
		return true
	}
	backward, err := filepath.Rel(second, first)
	if err == nil && filepath.IsLocal(backward) {
		return true
	}
	// EvalSymlinks need not normalize spelling on a case-insensitive volume.
	// Compare existing directory identities so a case alias cannot hide overlap.
	return physicalAncestor(first, second) || physicalAncestor(second, first)
}

func physicalAncestor(parent, child string) bool {
	info, err := os.Stat(parent)
	if err != nil {
		return false
	}
	for {
		other, err := os.Stat(child)
		if err == nil && os.SameFile(info, other) {
			return true
		}
		next := filepath.Dir(child)
		if next == child {
			return false
		}
		child = next
	}
}

// Only the existing export and live-preview workspaces may contain output
// inside .mpress. Never remove .mpress or a workspace's containing directory.
func temporaryOutput(state, output string) bool {
	relative, err := filepath.Rel(state, output)
	if err != nil || !filepath.IsLocal(relative) {
		return false
	}
	parts := strings.Split(relative, string(filepath.Separator))
	return len(parts) >= 2 && (parts[0] == "live-preview" || (strings.HasPrefix(parts[0], "export-") && len(parts[0]) > len("export-")))
}
