UNAME_OS := $(shell uname -s)
UNAME_ARC := $(shell uname -m)

os_suffix :=
static_linking := -ldflags '-extldflags "-static"'
ifeq ($(UNAME_OS), Linux)
os_name := linux
else ifeq ($(UNAME_OS), Darwin)
os_name := darwin
# to get it working the below arguments are removed (no static linking)
static_linking :=
else
os_name := windows
os_suffix := .exe
endif
ifeq ($(UNAME_ARC), arm64)
arc_name := arm64
else
arc_name := amd64
endif

help:
	@grep -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" Makefile \
		| grep -v -- -- \
		| sed 'N;s/\n/###/' \
		| sed -n 's/^#: \(.*\)###\(.*\):.*/\2:###\1/p' \
		| column -t  -s '###'

#: Install go dependencies
install-go:
	go mod download

#: Install Javascript dependencies
install-yarn:
	yarn install

#: Install all dependencies
install: install-go install-yarn

#: Teardown and start a local Grafana instance
bootstrap: teardown
	docker-compose up -d grafana
	@echo "Go to http://localhost:3000/"

#: Teardown the docker resources
teardown:
	docker-compose down --remove-orphans --volumes --timeout=2

#: Build the backend for the local architecture
build-backend:
	CGO_ENABLED=1 go build \
		-o dist/gpx_sqlite-datasource_$(os_name)_$(arc_name)$(os_suffix) \
		$(static_linking) -tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-win64:
	docker build -t cross-build ./build

	docker run -t -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=windows -e GOARCH=amd64 -e  CC=x86_64-w64-mingw32-gcc \
		cross-build \
		go build -x -o dist/gpx_sqlite-datasource_windows_amd64.exe \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux64:
	docker build -t cross-build ./build

	docker run -t -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=amd64 \
		cross-build \
		go build -x -o dist/gpx_sqlite-datasource_linux_amd64 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux-arm6:
	docker build -t cross-build ./build

	docker run -t -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=arm \
		-e CC=/opt/rpi-tools/arm-bcm2708/arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc \
		cross-build \
		go build -x -o dist/gpx_sqlite-datasource_linux_arm6 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux-arm7:
	docker build -t cross-build ./build

	docker run -t -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=arm \
		-e CC=arm-linux-gnueabihf-gcc \
		cross-build \
		go build -x -o dist/gpx_sqlite-datasource_linux_arm7 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux-arm64:
	docker build -t cross-build ./build

	docker run -t -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=arm64 -e CC=aarch64-linux-gnu-gcc \
		cross-build \
		go build -x -o dist/gpx_sqlite-datasource_linux_arm64 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

#: Build the frontend
build-frontend:
	yarn build

#: Package up the build artifacts and zip them in a file
package-and-zip:
	chmod +x ./dist/gpx_*
	cp -R dist dist_old

	mv dist/gpx_sqlite-datasource_linux_arm6 dist/gpx_sqlite-datasource_linux_arm
	rm dist/gpx_sqlite-datasource_linux_arm7
	yarn sign
	mv dist frser-sqlite-datasource
	zip frser-sqlite-datasource-$$(cat package.json | jq .version -r).zip ./frser-sqlite-datasource -r
	rm -rf frser-sqlite-datasource
	mv dist_old dist

#: Package up the build artifacts for an ARM 7 architecture and zip them in a file
package-and-zip-arm7:
	chmod +x ./dist/gpx_*
	cp -R dist dist_old

	rm dist/gpx_*
	cp dist_old/gpx_sqlite-datasource_linux_arm7 dist/gpx_sqlite-datasource_linux_arm
	yarn sign
	mv dist frser-sqlite-datasource
	zip frser-sqlite-datasource-arm7-$$(cat package.json | jq .version -r).zip ./frser-sqlite-datasource -r
	rm -rf frser-sqlite-datasource
	mv dist_old dist

#: Build the frontend and backend
build: build-frontend build-backend

#: Run the end-to-end tests with Selenium
selenium-test:
	@echo
	@echo "Make sure the plugin is built and signed for the architecture of the docker tests"
	@docker-compose rm --force --stop -v grafana
	GRAFANA_VERSION=7.3.3 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo
	GRAFANA_VERSION=8.1.0 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo

#: Run the frontend tests
frontend-test:
	yarn test

#: Run the backend tests
backend-test:
	@echo
	go test --tags="sqlite_omit_load_extension sqlite_json" ./pkg/...
	@echo

#: Sign the build artifacts with the private Grafana organization key
sign:
	yarn sign

#: Build all artifacts for the local architecture and sign them with the private Grafana organization key
build-and-sign: build sign

#: Run all tests (frontend, backend, end-to-end)
test: backend-test build-and-sign selenium-test
	docker-compose down --remove-orphans --volumes --timeout=2
