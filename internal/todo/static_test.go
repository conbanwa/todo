package todo

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupStaticTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Setup static file serving
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}

// TestStaticFileServing tests that HTML files are served correctly
func TestStaticFileServing(t *testing.T) {
	// Create temporary static directory and index.html
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	os.MkdirAll(staticDir, 0755)

	indexPath := filepath.Join(staticDir, "index.html")
	testContent := "<html><body>Test</body></html>"
	os.WriteFile(indexPath, []byte(testContent), 0644)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Static("/static", staticDir)
	r.GET("/", func(c *gin.Context) {
		c.File(indexPath)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	// Test root route
	resp, err := http.Get(s.URL + "/")
	if err != nil {
		t.Fatalf("failed to fetch root: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != testContent {
		t.Errorf("expected content %q, got %q", testContent, string(body))
	}

	// Test static file route
	resp2, err := http.Get(s.URL + "/static/index.html")
	if err != nil {
		t.Fatalf("failed to fetch static file: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for static file, got %d", resp2.StatusCode)
	}

	contentType := resp2.Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" && contentType != "text/html" {
		t.Errorf("expected Content-Type 'text/html', got %q", contentType)
	}
}

// TestStaticFileNotFound tests 404 for non-existent files
func TestStaticFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	os.MkdirAll(staticDir, 0755)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Static("/static", staticDir)

	s := httptest.NewServer(r)
	defer s.Close()

	// Request non-existent file
	resp, err := http.Get(s.URL + "/static/nonexistent.html")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

// TestRootRoute tests that root path serves index.html
func TestRootRoute(t *testing.T) {
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	os.MkdirAll(staticDir, 0755)

	indexPath := filepath.Join(staticDir, "index.html")
	testHTML := `<!DOCTYPE html><html><head><title>Todo App</title></head><body><h1>Todo List</h1></body></html>`
	os.WriteFile(indexPath, []byte(testHTML), 0644)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.File(indexPath)
	})

	s := httptest.NewServer(r)
	defer s.Close()

	resp, err := http.Get(s.URL + "/")
	if err != nil {
		t.Fatalf("failed to fetch root: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != testHTML {
		t.Errorf("expected HTML content, got different content")
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/html; charset=utf-8" && contentType != "text/html" {
		t.Errorf("expected Content-Type 'text/html', got %q", contentType)
	}
}
