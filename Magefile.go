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
	build.BuildAll()

	// freebsd is not included by default so we add it here
	err := sh.RunWith(
		map[string]string{"GOOS": "freebsd", "GOARCH": "amd64"},
		"go", "build", "-o", "dist/gpx_sqlite-datasource_freebsd_amd64", "./pkg",
	)

	if err != nil {
		return err
	}

	// Most 32-bit Raspberry Pi models have an ARMv6 architecture. (Pi Zero and 1 models)
	// All other models (2 Mod. B v1.2, 3 and 4) have an 64Bit ARMv8 architecture.
	// Only the Raspberry Pi 2 Mod. B has an ARMv7 architecture
	// To enable seamless installation for the raspberry pi zero we build for ARMv6 here.
	err = sh.RunWith(
		map[string]string{"GOOS": "linux", "GOARCH": "arm", "GOARM": "6"},
		"go", "build", "-o", "dist/gpx_sqlite-datasource_linux_arm", "./pkg",
	)

	if err != nil {
		return err
	}
	return nil
}
