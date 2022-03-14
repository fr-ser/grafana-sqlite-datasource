UNAME_OS := $(shell uname -s)
UNAME_ARC := $(shell uname -m)

os_suffix := .exe
arc_name := amd64
os_name := windows

ifeq ($(UNAME_OS), Linux)
os_name := linux
os_suffix :=
else ifeq ($(UNAME_OS), Darwin)
os_name := darwin
os_suffix :=
endif

ifeq ($(UNAME_ARC), arm64)
arc_name := arm64
endif

help:
	@grep -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" Makefile \
		| grep -v -- -- \
		| sed 'N;s/\n/###/' \
		| sed -n 's/^#: \(.*\)###\(.*\):.*/\2:###\1/p' \
		| column -t  -s '###'

#: Add git hooks of the project
add-git-hook:
	ln -s ../../githooks/pre-push .git/hooks/pre-push

#: Install go dependencies
install-go-dependencies:
	go mod download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

#: Install Javascript dependencies
install-js-dependencies:
	yarn install

#: Install all dependencies
install-dependencies: install-go-dependencies install-js-dependencies

#: Teardown and start a local Grafana instance
bootstrap: teardown
	docker-compose up -d grafana
	@echo "Go to http://localhost:3000/"

#: Teardown the docker resources
teardown:
	docker-compose down --remove-orphans --volumes --timeout=2

#: Build the backend for the local architecture
build-backend-local:
	go build -o dist/gpx_sqlite-datasource_$(os_name)_$(arc_name)$(os_suffix) ./pkg

#: Build the backend for the docker architecture
build-backend-docker: build-backend-cross-linux-amd64

#: Build the backend for all supported environments
build-backend-all: build-backend-cross-win-amd64 build-backend-cross-linux-amd64 build-backend-cross-linux-arm build-backend-cross-linux-arm64 build-backend-cross-freebsd-amd64 build-backend-cross-darwin-amd64 build-backend-cross-darwin-arm64

build-backend-cross-win-amd64:
	GOOS=windows GOARCH=amd64 go build -o dist/gpx_sqlite-datasource_windows_amd64.exe ./pkg

build-backend-cross-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/gpx_sqlite-datasource_linux_amd64 ./pkg

build-backend-cross-linux-arm:
	GOOS=linux GOARCH=arm go build -o dist/gpx_sqlite-datasource_linux_arm ./pkg

build-backend-cross-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o dist/gpx_sqlite-datasource_linux_arm64 ./pkg

build-backend-cross-freebsd-amd64:
	GOOS=freebsd GOARCH=amd64 go build -o dist/gpx_sqlite-datasource_freebsd_amd64 ./pkg

build-backend-cross-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o dist/gpx_sqlite-datasource_darwin_amd64 ./pkg

build-backend-cross-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o dist/gpx_sqlite-datasource_darwin_arm64 ./pkg

#: Build the frontend
build-frontend:
	yarn build --skipTest --skipLint

#: Package up the build artifacts and zip them in a file
package-and-zip:
	chmod +x ./dist/gpx_*
	cp -R dist dist_old

	yarn sign
	mv dist frser-sqlite-datasource
	zip frser-sqlite-datasource-$$(cat package.json | jq .version -r).zip ./frser-sqlite-datasource -r
	rm -rf frser-sqlite-datasource
	mv dist_old dist

#: Build the frontend and backend for the local and test environment
build: build-frontend build-backend-local build-backend-docker

#: Run the end-to-end tests with Selenium after building the backend for docker
selenium-test: build-backend-docker selenium-test-no-build

#: Run the end-to-end tests with Selenium without building the backend for docker. This can be helpful if the packages have already been built and signed
selenium-test-no-build:
	@echo
	@docker-compose rm --force --stop -v grafana
	GRAFANA_VERSION=7.3.3 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo
	GRAFANA_VERSION=8.1.0 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo

#: Run the frontend tests without linting and building the code
frontend-test-fast:
	yarn test

#: Run the frontend tests
frontend-test:
	# build is run as this is the only way to include linting
	yarn build

#: Run the backend tests
backend-test:
	gotestsum --format testname -- -count=1 -cover ./pkg/...
	@echo
	@echo "Linting Checks:"
	@golangci-lint run ./pkg/... && echo "Linting passed!\n"

#: Sign the build artifacts with the private Grafana organization key
sign:
	yarn sign

#: Build all artifacts for the local architecture and sign them with the private Grafana organization key
build-and-sign: build sign

#: Run all tests (frontend, backend, end-to-end)
test: 
	# clear the dist directory in case a previous version of the plugin was signed
	rm -rf dist
	make backend-test
	make frontend-test
	make selenium-test
	make teardown
