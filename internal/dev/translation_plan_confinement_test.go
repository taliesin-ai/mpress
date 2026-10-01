package dev

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/leaanthony/mpress/internal/translate"
)

func TestAuthoringTranslationPlanRejectsExternalFiles(t *testing.T) {
	for _, boundary := range []string{"source", "target", "state-file", "state-parent", "state-discovery", "style-guide", "glossary"} {
		t.Run(boundary, func(t *testing.T) {
			server := translationReviewFixture(t)
			var sentinel string
			var original []byte
			if boundary == "style-guide" || boundary == "glossary" {
				name, data := "style.txt", "Use the configured writing style.\n"
				if boundary == "glossary" {
					name, data = "glossary.yaml", "terms:\n  - source: compiler\n    translations:\n      fr: compilateur\n"
					server.cfg.Translation.Glossary = name
				} else {
					server.cfg.Translation.StyleGuide = name
				}
				sentinel = filepath.Join(t.TempDir(), name)
				original = []byte(data)
				if err := os.WriteFile(sentinel, original, 0600); err != nil {
					t.Fatal(err)
				}
				devSymlink(t, sentinel, filepath.Join(server.project, name))
			} else {
				sentinel, original = installTranslationReviewBoundary(t, server, boundary)
			}
			host := httptest.NewServer(server.Handler())
			defer host.Close()
			response := requestJSON(t, http.MethodGet, host.URL+"/__mpress/api/translations?lang=fr&file=index.md", nil, "translation-boundary")
			defer response.Body.Close()
			if response.StatusCode < 400 {
				t.Errorf("translation plan read external %s with HTTP %d", boundary, response.StatusCode)
			}
			assertDevFileUnchanged(t, filepath.Dir(sentinel), filepath.Base(sentinel), string(original))
		})
	}
}

type countedBoundaryProvider struct{ calls int }

func (p *countedBoundaryProvider) Name() string  { return "test" }
func (p *countedBoundaryProvider) Model() string { return "test" }
func (p *countedBoundaryProvider) Translate(ctx context.Context, request translate.TranslationRequest) (map[string]string, error) {
	p.calls++
	return pageUpdateProvider{}.Translate(ctx, request)
}

type rearrangingBoundaryProvider struct {
	hook  func() error
	calls int
}

func (p *rearrangingBoundaryProvider) Name() string  { return "test" }
func (p *rearrangingBoundaryProvider) Model() string { return "test" }
func (p *rearrangingBoundaryProvider) Translate(ctx context.Context, request translate.TranslationRequest) (map[string]string, error) {
	if p.calls == 0 {
		if err := p.hook(); err != nil {
			return nil, err
		}
	}
	p.calls++
	return pageUpdateProvider{}.Translate(ctx, request)
}

func TestTranslationRunRejectsParentsRearrangedDuringProvider(t *testing.T) {
	for _, boundary := range []string{"target", "state-file"} {
		t.Run(boundary, func(t *testing.T) {
			server := translationReviewFixture(t)
			path := translationReviewPath(server, boundary)
			parent := filepath.Dir(path)
			outside := t.TempDir()
			sentinel := filepath.Join(outside, filepath.Base(path))
			if err := os.WriteFile(sentinel, []byte("outside must survive"), 0600); err != nil {
				t.Fatal(err)
			}
			// Probe native support before the worker callback; later fixture
			// creation failures are returned and must fail this case.
			probe := filepath.Join(outside, "probe")
			devSymlink(t, "missing", probe)
			if err := os.Remove(probe); err != nil {
				t.Fatal(err)
			}
			provider := &rearrangingBoundaryProvider{hook: func() error {
				if err := os.Rename(parent, parent+".owned"); err != nil {
					return err
				}
				return os.Symlink(outside, parent)
			}}
			_, err := translate.NewEngine(server.project, server.cfg, provider).RunRoot(context.Background(), server.files, translate.Options{Language: "fr", File: "index.md", Scope: "all", Force: true})
			if err == nil || provider.calls != 1 {
				t.Errorf("provider rearrangement escaped boundary: calls=%d error=%v", provider.calls, err)
			}
			assertDevFileUnchanged(t, outside, filepath.Base(sentinel), "outside must survive")
		})
	}
}

func TestTranslationRunRejectsExternalFilesBeforeProvider(t *testing.T) {
	for _, boundary := range []string{"source", "target", "state-file", "state-parent", "state-discovery"} {
		t.Run(boundary, func(t *testing.T) {
			server := translationReviewFixture(t)
			sentinel, original := installTranslationReviewBoundary(t, server, boundary)
			provider := &countedBoundaryProvider{}
			_, err := translate.NewEngine(server.project, server.cfg, provider).RunRoot(context.Background(), server.files, translate.Options{Language: "fr", File: "index.md", Scope: "all", Force: true})
			if err == nil || provider.calls != 0 {
				t.Errorf("unsafe run reached provider: calls=%d error=%v", provider.calls, err)
			}
			assertDevFileUnchanged(t, filepath.Dir(sentinel), filepath.Base(sentinel), string(original))
		})
	}
}

func TestTranslationRunRetainsInternalAliases(t *testing.T) {
	for _, boundary := range []string{"source", "target", "state-file", "state-parent"} {
		t.Run(boundary, func(t *testing.T) {
			server := translationReviewFixture(t)
			path := translationReviewPath(server, boundary)
			if boundary == "state-parent" {
				path = filepath.Dir(path)
			}
			saved := path + ".saved"
			if err := os.Rename(path, saved); err != nil {
				t.Fatal(err)
			}
			devSymlink(t, filepath.Base(saved), path)
			provider := &countedBoundaryProvider{}
			report, err := translate.NewEngine(server.project, server.cfg, provider).RunRoot(context.Background(), server.files, translate.Options{Language: "fr", File: "index.md", Scope: "all", Force: true})
			if err != nil || report.Written != 1 || provider.calls == 0 {
				t.Fatalf("safe run alias failed: %v %+v calls=%d", err, report, provider.calls)
			}
			response := translationReviewRequest(t, server)
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				t.Fatalf("written translation cannot be reviewed: %d", response.StatusCode)
			}
			assertReviewedTranslationState(t, translationReviewPath(server, "state-file"))
		})
	}
}
