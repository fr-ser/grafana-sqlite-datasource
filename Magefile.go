//+build mage

package main

import (
	"fmt"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// Install dependencies (yarn and go mod)
func Install() error {
	if err := sh.Run("yarn", "install"); err != nil {
		return err
	}
	return sh.Run("go", "mod", "download")
}

// BuildRelease builds the JavaScript frontend and Go Backend code
func BuildRelease() error {
	if err := sh.Run("yarn", "build"); err != nil {
		return err
	}
	build.BuildAll()
	return nil
}

// Bootstrap starts the docker-compose file (grafana mostly)
func Bootstrap() error {
	if err := Teardown(); err != nil {
		return err
	}

	if err := sh.Run("docker-compose", "up", "-d", "grafana"); err != nil {
		return err
	}
	fmt.Println("Go to http://localhost:3000/")
	return nil
}

// Teardown starts the docker-compose file (grafana mostly)
func Teardown() error {
	if err := sh.Run(
		"docker-compose", "down", "--remove-orphans", "--volumes", "--timeout=2",
	); err != nil {
		return err
	}
	fmt.Println("Go to http://localhost:3000/")
	return nil
}

// Selenium tests the plugin via selenium
func Selenium() error {
	if err := Teardown(); err != nil {
		return err
	}

	if err := sh.Run("docker-compose", "run", "--rm", "start-setup"); err != nil {
		return err
	}

	return sh.Run("yarn", "test")
}

// Default configures the default target.
var Default = BuildRelease
