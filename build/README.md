# Build container to allow cross compilation

This is heavily inspired by: https://github.com/grafana/grafana/tree/master/scripts/build/ci-build

## Problem

The main issue is, that the go SQLite package is using C-Code. Therefore cross compilation (with
static linking) requires the corresponding C-CrossCompiler.

## MacOS SDK

In order to build the image Download MacOSX10.15.sdk.tar.xz from
https://github.com/phracker/MacOSX-SDKs/releases/tag/10.15
into this folder

## CC

```sh
##########
CCARMV6=/opt/rpi-tools/arm-bcm2708/arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc
CCARMV7=arm-linux-gnueabihf-gcc
CCARMV7_MUSL=/tmp/arm-linux-musleabihf-cross/bin/arm-linux-musleabihf-gcc
CCARM64=aarch64-linux-gnu-gcc
CCARM64_MUSL=/tmp/aarch64-linux-musl-cross/bin/aarch64-linux-musl-gcc
CCX64=/tmp/x86_64-centos6-linux-gnu/bin/x86_64-centos6-linux-gnu-gcc
CCX64_MUSL=/tmp/x86_64-linux-musl-cross/bin/x86_64-linux-musl-gcc
```
