PLUGIN_VERSION ?= $(shell cat package.json | jq .version -r)

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
	# install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin 
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

#: Install Javascript dependencies
install-js-dependencies:
	npm install

#: Install all dependencies
install-dependencies: install-go-dependencies install-js-dependencies

#: Teardown and start a local Grafana instance
start: teardown
	docker-compose up -d grafana
	@echo "Go to http://localhost:3000/"

#: Teardown the docker resources
teardown:
	docker-compose down --remove-orphans --volumes --timeout=2

#: Build the backend for the local architecture
build-backend-local:
	go build -o dist/gpx_sqlite-datasource_$(os_name)_$(arc_name)$(os_suffix) ./pkg

#: Build the backend for the docker architecture
build-backend-docker:
	GOOS=linux GOARCH=$(arc_name) go build -o dist/gpx_sqlite-datasource_linux_$(arc_name) ./pkg

#: Build the backend for all supported environments
build-backend-all:
	mage BuildAllAndMore

#: Build the frontend
build-frontend:
	npm run build

#: Build the frontend and watch for changes
build-frontend-watch:
	npm run dev

#: Package up the build artifacts and zip them in a file
package-and-zip:
	chmod +x ./dist/gpx_*
	cp -R dist dist_old

	npm run sign
	mv dist frser-sqlite-datasource
	zip frser-sqlite-datasource-$(PLUGIN_VERSION).zip ./frser-sqlite-datasource -r
	rm -rf frser-sqlite-datasource
	mv dist_old dist

#: Build the frontend and backend for the local and test environment
build: build-frontend build-backend-local build-backend-docker

#: Run the end-to-end tests with Selenium after building the backend for docker
test-e2e: build-backend-docker build-frontend test-e2e-no-build

#: Run the end-to-end tests with Selenium without building the backend for docker. This can be helpful if the packages have already been built and signed
test-e2e-no-build:
	@echo
	@docker-compose rm --force --stop -v grafana
	GRAFANA_VERSION=7.3.3 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo
	GRAFANA_VERSION=8.1.0 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo

#: Run the frontend tests
test-frontend:
	npm run test:ci
	npm run lint
	npm run typecheck

#: Run the backend tests
test-backend:
	gotestsum --format testname -- -count=1 -cover ./pkg/...
	@echo
	@echo "Linting Checks:"
	@golangci-lint run ./pkg/... && echo "Linting passed!\n"

#: Run all tests (frontend, backend, end-to-end)
test: 
	# clear the dist directory in case a previous version of the plugin was signed
	rm -rf dist
	make test-backend
	make test-frontend
	make test-e2e
	make teardown
