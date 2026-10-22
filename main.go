package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"

	"go-templ-mock-1/assets"
	"go-templ-mock-1/pages"
)

// greet renders an HTMX fragment that replaces the greeting target.
func greet(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = "stranger"
	}
	templ.Handler(pages.Greeting(name)).ServeHTTP(w, r)
}

// submitForm renders an HTMX fragment confirming the contact form.
func submitForm(w http.ResponseWriter, r *http.Request) {
	templ.Handler(pages.FormResult(strings.TrimSpace(r.FormValue("name")), r.FormValue("topic"))).ServeHTTP(w, r)
}

// slowResult simulates a slow backend so the htmx indicator is visible.
func slowResult(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	templ.Handler(pages.SlowResult("Loaded at " + time.Now().Format("15:04:05"))).ServeHTTP(w, r)
}

// revenueChart renders a chart with fresh random data for HTMX to swap in.
func revenueChart(w http.ResponseWriter, r *http.Request) {
	templ.Handler(pages.RevenueChart(pages.RandomRevenue())).ServeHTTP(w, r)
}

func main() {
	mux := http.NewServeMux()
	setupAssetsRoutes(mux)
	mux.Handle("GET /", templ.Handler(pages.Home()))
	mux.HandleFunc("POST /greet", greet)
	mux.Handle("GET /forms", templ.Handler(pages.Forms()))
	mux.HandleFunc("POST /forms/submit", submitForm)
	mux.Handle("GET /feedback", templ.Handler(pages.Feedback()))
	mux.HandleFunc("GET /feedback/slow", slowResult)
	mux.Handle("GET /overlays", templ.Handler(pages.Overlays()))
	mux.Handle("GET /data", templ.Handler(pages.Data()))
	mux.Handle("GET /charts", templ.Handler(pages.Charts()))
	mux.Handle("GET /analytics", templ.Handler(pages.Analytics()))
	mux.HandleFunc("GET /analytics/revenue", revenueChart)

	ln, err := listen()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Server is running on http://localhost:%d\n", ln.Addr().(*net.TCPAddr).Port)
	if err := http.Serve(ln, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// listen binds the port. An explicit PORT binds exactly (the templ proxy in
// the Taskfile points at it, so silently moving would break hot reload);
// without PORT it walks from 8090 to the next free port like node dev
// servers do.
func listen() (net.Listener, error) {
	if env := os.Getenv("PORT"); env != "" {
		port, err := strconv.Atoi(env)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT %q", env)
		}
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			return nil, fmt.Errorf("port %d is busy - stop the other app or run with another port, e.g. 'task dev PORT=%d'", port, port+1)
		}
		return ln, nil
	}
	for port := 8090; port < 8100; port++ {
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err == nil {
			return ln, nil
		}
	}
	return nil, fmt.Errorf("no free port found from 8090 upwards")
}

func setupAssetsRoutes(mux *http.ServeMux) {
	development := isDevelopment()
	files := http.FileServer(http.FS(assets.Assets))
	if development {
		files = http.FileServer(http.Dir("./assets"))
	}
	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if development {
			w.Header().Set("Cache-Control", "no-store")
		} else if strings.HasPrefix(strings.TrimPrefix(r.URL.Path, "/"), "js/shadcn-templ-") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))
}

func isDevelopment() bool {
	// GO_ENV=development is a deprecated alias, removed after this minor version.
	return os.Getenv("SHADCN_TEMPL_DEV") == "true" || os.Getenv("TEMPL_DEV_MODE") == "true" || os.Getenv("GO_ENV") == "development"
}
