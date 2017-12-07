package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/pkg/errors"
	"github.com/tj/go/env"
	"github.com/tj/go/http/request"
	"github.com/tj/go/http/response"
)

// max age (one day).
var maxage = 60 * 60 * 24

// boot scripts.
var boot = []string{
	"tar -zxf go.tar.gz -C /tmp",
	"tar -xf git.tar -C /tmp",
}

// init
func init() {
	log.SetHandler(json.Default)
	os.Setenv("GOPATH", "/tmp")
	os.Setenv("GOROOT", "/tmp/go")
	os.Setenv("PATH", os.Getenv("PATH")+":/tmp/go/bin:/tmp/usr/bin")
	os.Setenv("PATH", os.Getenv("PATH")+":/tmp")
	os.Setenv("GIT_TEMPLATE_DIR", "/tmp/usr/share/git-core/templates")
	os.Setenv("GIT_EXEC_PATH", "/tmp/usr/libexec/git-core")
	os.Setenv("LD_LIBRARY_PATH", "/tmp/usr/lib64")
	if err := commands(boot); err != nil {
		log.Fatalf("error: %s", err)
	}
}

// commands utility.
func commands(cmds []string) error {
	for _, c := range cmds {
		fmt.Printf("RUN %s:\n", c)
		cmd := exec.Command("sh", "-c", c)
		if s, err := output(cmd); err == nil {
			fmt.Printf("%s\n", s)
		} else {
			return errors.Wrap(err, c)
		}
	}
	return nil
}

// main
func main() {
	addr := ":" + env.GetDefault("PORT", "3000")
	h := http.HandlerFunc(getBinary)
	err := http.ListenAndServe(addr, h)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// GET /:pkg
//
// Install a package with optional "os" and "arch"
// parameters, defaulting to "darwin" and "amd64".
//
func getBinary(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.URL.Path == "/_health" {
		response.OK(w)
		return
	}

	pkg := r.URL.Path[1:]

	if pkg == "" {
		response.BadRequest(w)
		return
	}

	goos := request.ParamDefault(r, "os", "darwin")
	arch := request.ParamDefault(r, "arch", "amd64")

	ctx := log.WithFields(log.Fields{
		"pkg":  pkg,
		"os":   goos,
		"arch": arch,
	})

	// fetch
	ctx.Info("fetching")
	if err := get(pkg); err != nil {
		log.WithError(err).Error("fetching")
		response.InternalServerError(w)
		return
	}

	dir := filepath.Join(os.Getenv("GOPATH"), "src", pkg)
	dst := "/tmp/out"

	// build
	ctx.WithField("dir", dir).Info("building")
	if err := build(dir, dst, goos, arch); err != nil {
		log.WithError(err).Error("building")
		response.InternalServerError(w)
		return
	}

	d := time.Since(start)
	ctx.WithField("duration", int(d/time.Millisecond)).Info("built")

	// respond
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxage))
	w.Header().Set("Content-Type", "application/octet-stream")

	// copy
	f, err := os.Open(dst)
	if err != nil {
		ctx.WithError(err).Error("opening")
		response.InternalServerError(w)
		return
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		ctx.WithError(err).Error("writing")
	}
}

// get performs a `go get`.
func get(pkg string) error {
	cmd := exec.Command("go", "get", pkg)
	_, err := output(cmd)
	return err
}

// build performs a `go build`.
func build(dir, dst, goos, arch string) error {
	cmd := exec.Command("go", "build", "-o", dst)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS="+goos)
	cmd.Env = append(cmd.Env, "GOARCH="+arch)
	_, err := output(cmd)
	return err
}

// output from a command.
func output(cmd *exec.Cmd) (string, error) {
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, strings.TrimSpace(string(b)))
	}

	return string(b), nil
}
