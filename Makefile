UNAME_OS := $(shell uname -s)
UNAME_ARC := $(shell uname -m)

install-go:
	go mod download

install-yarn:
	yarn install

install: install-go install-yarn

bootstrap: teardown
	docker-compose up -d grafana
	@echo "Go to http://localhost:3000/"

teardown:
	docker-compose down --remove-orphans --volumes --timeout=2

build-backend:
ifeq ($(UNAME_OS), Linux)
	CGO_ENABLED=1 go build \
		-o dist/gpx_sqlite-datasource_linux_amd64 \
		-ldflags '-extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg
else ifeq ($(UNAME_OS), Darwin)
	# to get it working the below arguments are removed (no static linking):
	# -ldflags '-extldflags "-static"'
	CGO_ENABLED=1 go build \
		-o dist/gpx_sqlite-datasource_darwin_amd64 \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg
else
	CGO_ENABLED=1 go build \
		-o dist/gpx_sqlite-datasource_windows_amd64.exe \
		-ldflags '-extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg
endif

build-backend-cross-win64:
	@docker build -t cross-build ./build

	docker run -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=windows -e GOARCH=amd64 -e  CC=x86_64-w64-mingw32-gcc \
		cross-build \
		go build -o dist/gpx_sqlite-datasource_windows_amd64.exe \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux64:
	@docker build -t cross-build ./build

	docker run -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=amd64 \
		cross-build \
		go build -o dist/gpx_sqlite-datasource_linux_amd64 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux-arm6:
	@docker build -t cross-build ./build

	docker run -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=arm \
		-e CC=/opt/rpi-tools/arm-bcm2708/arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc \
		cross-build \
		go build -o dist/gpx_sqlite-datasource_linux_arm6 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux-arm7:
	docker build -t cross-build ./build

	docker run -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=arm \
		-e CC=arm-linux-gnueabihf-gcc \
		cross-build \
		go build -o dist/gpx_sqlite-datasource_linux_arm7 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-backend-cross-linux-arm64:
	docker build -t cross-build ./build

	docker run -v "$${PWD}":/usr/src/app -w /usr/src/app \
		-e CGO_ENABLED=1 -e GOOS=linux -e GOARCH=arm64 -e CC=aarch64-linux-gnu-gcc \
		cross-build \
		go build -o dist/gpx_sqlite-datasource_linux_arm64 \
		-ldflags '-w -s -extldflags "-static"' \
		-tags osusergo,netgo,sqlite_omit_load_extension,sqlite_json \
		./pkg

build-frontend:
	yarn build

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

build: build-frontend build-backend

selenium-test:
	@echo
	@docker-compose rm --force --stop -v grafana
	GRAFANA_VERSION=7.3.3 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo
	GRAFANA_VERSION=8.1.0 docker-compose run --rm start-setup
	npx jest --runInBand --testMatch '<rootDir>/selenium/**/*.test.{js,ts}'
	@echo

frontend-test:
	yarn test

backend-test:
	@echo
	go test --tags="sqlite_omit_load_extension sqlite_json" ./pkg/...
	@echo

sign:
	yarn sign

build-and-sign: build sign

test: backend-test build-and-sign selenium-test
	docker-compose down --remove-orphans --volumes --timeout=2
