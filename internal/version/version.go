package version

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/leaanthony/mpress/internal/config"
	"github.com/leaanthony/mpress/internal/icons"
	"github.com/leaanthony/mpress/internal/projectfs"
	"github.com/leaanthony/mpress/internal/routes"
)

type Manifest struct {
	SchemaVersion int               `json:"schemaVersion"`
	Version       string            `json:"version"`
	CreatedAt     string            `json:"createdAt"`
	Files         map[string]string `json:"files"`
}

var labelRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

var (
	doubleQuotedRootURLAttributeRE = regexp.MustCompile(`(?i)\b(href|src|action|data-[a-z0-9-]+)="(/[^"]*)"`)
	singleQuotedRootURLAttributeRE = regexp.MustCompile(`(?i)\b(href|src|action|data-[a-z0-9-]+)='(/[^']*)'`)
)

const (
	versionMenuStart = `<!--mpress-version-menu:start-->`
	versionMenuEnd   = `<!--mpress-version-menu:end-->`
)

func Capture(project, label string, force bool) error {
	return capture(project, label, force)
}
func List(project string) ([]string, error) {
	files, _, err := openVersionStore(project)
	if errors.Is(err, errMissingStore) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer files.Close()
	entries, err := files.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() && validateLabel(e.Name()) == nil {
			complete, err := hasSnapshotManifest(files, e.Name())
			if err != nil {
				return nil, err
			}
			if complete {
				out = append(out, e.Name())
			}
		}
	}
	sort.Strings(out)
	return out, nil
}

// Listing requires a manifest file, not a full checksum verification. Keep
// incomplete snapshots hidden, while diagnosing unsafe or unreadable inputs.
func hasSnapshotManifest(files *projectfs.FS, path string) (bool, error) {
	snapshot, err := files.Sub(path)
	if err != nil {
		return false, err
	}
	defer snapshot.Close()
	info, err := snapshot.Stat("mpress-version.json")
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}

func Verify(project, label string) error {
	if err := validateLabel(label); err != nil {
		return err
	}
	files, cfg, err := openVersionStore(project)
	if err != nil {
		return err
	}
	defer files.Close()
	path := filepath.Join(cfg.ArtifactsPath(project), label)
	if err := cfg.SafeVersionPath(project, path); err != nil {
		return err
	}
	snapshot, err := files.Sub(label)
	if err != nil {
		return err
	}
	defer snapshot.Close()
	return verifySnapshot(snapshot, label)
}

// Verification uses one pinned snapshot for its manifest, checksums and
// enumeration; callers must keep that boundary open throughout the operation.
func verifySnapshot(snapshot *projectfs.FS, label string) error {
	data, err := snapshot.ReadFile("mpress-version.json")
	if err != nil {
		return err
	}
	var m Manifest
	if err = json.Unmarshal(data, &m); err != nil {
		return err
	}
	if m.SchemaVersion != 1 {
		return fmt.Errorf("unsupported version manifest schema %d", m.SchemaVersion)
	}
	if m.Version != label {
		return fmt.Errorf("manifest version %q does not match %q", m.Version, label)
	}
	for rel, want := range m.Files {
		path, pathErr := manifestFilePath(".", rel)
		if pathErr != nil {
			return pathErr
		}
		data, err = snapshot.ReadFile(path)
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != want {
			return fmt.Errorf("checksum mismatch: %s", rel)
		}
	}
	if err = snapshot.WalkDir(".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		rel, relErr := filepath.Rel(".", path)
		if relErr != nil {
			return relErr
		}
		rel = filepath.ToSlash(rel)
		if rel == "mpress-version.json" {
			return nil
		}
		if _, ok := m.Files[rel]; !ok {
			return fmt.Errorf("unexpected file: %s", rel)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func manifestFilePath(root, rel string) (string, error) {
	if routes.Output(rel) != nil {
		return "", fmt.Errorf("invalid manifest path: %s", rel)
	}
	return filepath.Join(root, filepath.FromSlash(rel)), nil
}
func Remove(project, label string) error {
	if err := validateLabel(label); err != nil {
		return err
	}
	files, cfg, err := openVersionStore(project)
	if errors.Is(err, errMissingStore) {
		return nil
	}
	if err != nil {
		return err
	}
	defer files.Close()
	path := filepath.Join(cfg.ArtifactsPath(project), label)
	if err := cfg.SafeVersionPath(project, path); err != nil {
		return err
	}
	lock, err := files.Lock(storeLockFile)
	if err != nil {
		return err
	}
	defer lock.Close()
	return files.RemoveAll(label)
}
func Mount(project, output string) (int, error) {
	cfg, err := config.Load(project)
	if err != nil || !cfg.Version.Enabled {
		return 0, err
	}
	labels, err := List(project)
	if err != nil {
		return 0, err
	}
	mounted := 0
	for _, label := range labels {
		if err = Verify(project, label); err != nil {
			return mounted, err
		}
		src := filepath.Join(cfg.ArtifactsPath(project), label)
		dest := filepath.Join(output, "versions", label)
		if err = copyTree(src, dest, func(rel string) bool { return filepath.Base(rel) == "mpress-version.json" }); err != nil {
			return mounted, err
		}
		if err = rewriteMountedVersion(dest, label, cfg.Version.Current, labels); err != nil {
			return mounted, err
		}
		mounted++
	}
	data, _ := json.MarshalIndent(map[string]any{"current": cfg.Version.Current, "versions": labels}, "", "  ")
	if err = os.MkdirAll(filepath.Join(output, "versions"), 0o755); err != nil {
		return mounted, err
	}
	err = os.WriteFile(filepath.Join(output, "versions", "versions.json"), data, 0o644)
	return mounted, err
}

func rewriteMountedVersion(root, label, currentLabel string, labels []string) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".html") {
			return walkErr
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		rewritten := rewriteRootURLs(string(data), label)
		rewritten = replaceVersionMenu(rewritten, versionRoute(filepath.ToSlash(rel)), label, currentLabel, labels)
		info, infoErr := entry.Info()
		if infoErr != nil {
			return infoErr
		}
		return os.WriteFile(path, []byte(rewritten), info.Mode().Perm())
	})
}

func rewriteRootURLs(markup, label string) string {
	prefix := "/versions/" + label + "/"
	rewrite := func(expression *regexp.Regexp, quote string, input string) string {
		return expression.ReplaceAllStringFunc(input, func(match string) string {
			parts := expression.FindStringSubmatch(match)
			path := parts[2]
			if strings.HasPrefix(path, "//") || strings.HasPrefix(path, "/versions/") {
				return match
			}
			return parts[1] + "=" + quote + prefix + strings.TrimPrefix(path, "/") + quote
		})
	}
	markup = rewrite(doubleQuotedRootURLAttributeRE, `"`, markup)
	return rewrite(singleQuotedRootURLAttributeRE, `'`, markup)
}

func versionRoute(rel string) string {
	if rel == "index.html" {
		return "/"
	}
	if strings.HasSuffix(rel, "/index.html") {
		return "/" + strings.TrimSuffix(rel, "index.html")
	}
	return "/" + rel
}

func replaceVersionMenu(markup, route, activeLabel, currentLabel string, labels []string) string {
	start := strings.Index(markup, versionMenuStart)
	end := strings.Index(markup, versionMenuEnd)
	if start < 0 || end < start {
		return markup
	}
	end += len(versionMenuEnd)
	menu := versionMenuMarkup(route, activeLabel, currentLabel, labels)
	return markup[:start] + versionMenuStart + menu + versionMenuEnd + markup[end:]
}

func versionMenuMarkup(route, activeLabel, currentLabel string, labels []string) string {
	if len(labels) == 0 {
		return ""
	}
	if strings.TrimSpace(currentLabel) == "" {
		currentLabel = "Current"
	}
	var out strings.Builder
	out.WriteString(`<div class="header-group utility-select utility-menu version-select"><button class="utility-menu-trigger" type="button" popovertarget="mpress-version-menu" aria-label="Select version" aria-haspopup="menu" aria-expanded="false">`)
	out.WriteString(icons.Lucide("git-branch", 19))
	out.WriteString(`<span class="utility-menu-current">` + html.EscapeString(activeLabel) + `</span>`)
	out.WriteString(icons.Lucide("chevron-down", 13))
	out.WriteString(`</button><menu id="mpress-version-menu" class="utility-menu-panel utility-version-menu" popover data-utility-menu>`)
	writeVersionLink := func(label, url string, current bool) {
		out.WriteString(`<li><a href="` + html.EscapeString(url) + `" role="menuitem"`)
		if current {
			out.WriteString(` aria-current="true"`)
		}
		out.WriteString(`><span class="utility-version-label">` + html.EscapeString(label) + `</span>`)
		if current {
			out.WriteString(`<small>Current documentation</small>`)
		}
		out.WriteString(`</a></li>`)
	}
	writeVersionLink(currentLabel, route, false)
	for _, label := range labels {
		writeVersionLink(label, "/versions/"+label+route, label == activeLabel)
	}
	out.WriteString(`</menu></div>`)
	return out.String()
}
func copyTree(src, dst string, skip func(string) bool) error {
	return filepath.WalkDir(src, func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		if rel == "." {
			return os.MkdirAll(dst, 0o755)
		}
		if skip != nil && skip(rel) {
			if e.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		if e.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, in)
		closeErr := out.Close()
		if err != nil {
			return err
		}
		return closeErr
	})
}
