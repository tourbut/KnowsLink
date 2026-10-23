// Automatic setup distribution (SAR-AUTOSETUP-001): the pinned installer and plugin ZIPs written by
// `scripts/package_plugin.py --release` are served read-only from InstallDir at /install/<version>/<file>.
package relay

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

var installVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)

// installFiles is the complete published set; anything else under InstallDir is never served.
var installFiles = map[string]string{
	"knowslink-setup.ps1":         "text/plain; charset=utf-8", // irm decodes the Korean installer text by this charset.
	"knowslink-claude-plugin.zip": "application/zip",
	"knowslink-codex-plugin.zip":  "application/zip",
}

func (s *Service) installFile(w http.ResponseWriter, r *http.Request) {
	version, name := r.PathValue("version"), r.PathValue("file")
	kind, ok := installFiles[name]
	if s.InstallDir == "" || !ok || !installVersion.MatchString(version) {
		http.NotFound(w, r)
		return
	}
	data, err := os.ReadFile(filepath.Join(s.InstallDir, version, name))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", kind)
	_, _ = w.Write(data)
}
