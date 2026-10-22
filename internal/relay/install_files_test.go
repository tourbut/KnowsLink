// Auto-setup distribution serves only the published file set of a semantic version, with the installer's charset.
package relay

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallFilesServeOnlyPublishedSet(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "0.1.0"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"0.1.0/knowslink-setup.ps1": "# 설치기", "0.1.0/notes.txt": "x", "secret.txt": "x"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, configured := range []string{dir, ""} {
		s := &Service{InstallDir: configured}
		mux := http.NewServeMux()
		mux.HandleFunc("GET /install/{version}/{file}", s.installFile)
		for target, want := range map[string]int{
			"/install/0.1.0/knowslink-setup.ps1":           200,
			"/install/0.1.0/knowslink-claude-plugin.zip":   404,
			"/install/0.1.0/notes.txt":                     404,
			"/install/0.1/knowslink-setup.ps1":             404,
			"/install/..%2Fsecret.txt/knowslink-setup.ps1": 404,
			"/install/0.1.0/..%2F..%2Fsecret.txt":          404,
		} {
			if configured == "" && want == 200 {
				want = 404
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest("GET", target, nil))
			if w.Code != want {
				t.Fatalf("%q %s = %d, want %d", configured, target, w.Code, want)
			}
			if want == 200 && (w.Header().Get("Content-Type") != "text/plain; charset=utf-8" || w.Body.String() != "# 설치기") {
				t.Fatalf("installer served as %q %q", w.Header().Get("Content-Type"), w.Body.String())
			}
		}
	}
}
