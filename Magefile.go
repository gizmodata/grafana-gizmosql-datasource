//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Build builds the backend plugin for the current platform
func Build() error {
	return sh.RunV("go", "build", "-o", "dist/gpx_gizmosql_datasource", "./pkg")
}

// BuildAll builds the backend plugin for all supported platforms
func BuildAll() error {
	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	for _, p := range platforms {
		env := map[string]string{
			"GOOS":   p.os,
			"GOARCH": p.arch,
		}

		suffix := ""
		if p.os == "windows" {
			suffix = ".exe"
		}

		output := "dist/gpx_gizmosql_datasource_" + p.os + "_" + p.arch + suffix
		if err := sh.RunWithV(env, "go", "build", "-o", output, "./pkg"); err != nil {
			return err
		}
	}

	return nil
}

// Test runs the Go tests
func Test() error {
	return sh.RunV("go", "test", "-v", "./...")
}

// Clean removes built artifacts
func Clean() error {
	return sh.Rm("dist")
}
