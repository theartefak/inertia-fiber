package utils

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// vite returns the HTML for the specified entrypoints.
func Vite(entrypoints []string, buildDirectory ...string) template.HTML {
	// If running in hot mode, return the HTML for the hot asset.
	if isRunningHot() {
		html := makeTagForChunk(hotAsset("@vite/client"))
		for _, v := range entrypoints {
			html += makeTagForChunk(hotAsset(v))
		}
		return html
	}

	// Otherwise, return the HTML for the specified entrypoints.
	manifest := manifest(buildDirectory...)
	html := template.HTML("")

	for _, v := range entrypoints {
		m := manifest[v]
		for _, css := range m.Css {
			html += template.HTML(makeStylesheetTag(css))
		}
		html += template.HTML(makeScriptTag(m.File))
	}

	return html
}

// makeTagForChunk returns the HTML tag for the specified URL.
func makeTagForChunk(url string) template.HTML {
	if isCssPath(url) {
		return template.HTML(makeStylesheetTag(url))
	}

	return template.HTML(makeScriptTag(url))
}

// makeStylesheetTag returns the HTML tag for the specified stylesheet URL.
func makeStylesheetTag(url string) string {
	return fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, url)
}

// makeScriptTag returns the HTML tag for the specified script URL.
func makeScriptTag(url string) string {
	return fmt.Sprintf(`<script type="module" src="%s"></script>`, url)
}

// isCssPath returns true if the path is a CSS path.
func isCssPath(path string) bool {
	return regexp.MustCompile(`\.(css|less|sass|scss|styl|stylus|pcss|postcss)$`).MatchString(path)
}

// hotAsset returns the hot asset for the specified asset.
func hotAsset(asset string) string {
	data, err := os.ReadFile(hotFile())
	if err != nil {
		panic(err)
	}

	return strings.TrimSuffix(string(data), "\n") + "/" + asset
}

// isRunningHot returns true if running in hot mode.
func isRunningHot() bool {
	filename := hotFile()
	if filename == "" {
		return false
	}

	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

// hotFile returns the hot file path.
func hotFile() string {
	return filepath.Join("public", "hot")
}
