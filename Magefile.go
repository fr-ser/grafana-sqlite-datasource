//go:build mage
// +build mage

package main

import (
	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// BuildAllAndMore uses the official BuildAll and runs some more steps
func BuildAllAndMore() error {
	err := sh.RunWith(
		map[string]string{"GOOS": "freebsd", "GOARCH": "amd64"},
		"go", "build", "-o", "dist/gpx_sqlite-datasource_freebsd_amd64", "./pkg",
	)

	if err != nil {
		return err
	}
	build.BuildAll()
	return nil
}
