#!/bin/bash

BUILD_DIR=$(dirname "$0")/build
mkdir -p $BUILD_DIR
cd $BUILD_DIR

export GO111MODULE=on
echo "Setting GO111MODULE to" $GO111MODULE

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

LDFLAGS="-s -w"
GCFLAGS=""

# AMD64
OSES=(linux darwin windows freebsd)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env CGO_ENABLED=0 GOOS=$os GOARCH=amd64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o proxy_${os}_amd64${suffix} feather-proxy/cmd
#	if $UPX; then upx -9 proxy_${os}_amd64${suffix};fi
done

# 386
OSES=(linux windows)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env CGO_ENABLED=0 GOOS=$os GOARCH=386 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o proxy_${os}_386${suffix} feather-proxy/cmd
#	if $UPX; then upx -9 proxy_${os}_386${suffix};fi
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o proxy_linux_arm$v feather-proxy/cmd
#  if $UPX; then upx -9 proxy_linux_arm$v;fi
done

# ARM64
OSES=(linux darwin windows)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env CGO_ENABLED=0 GOOS=$os GOARCH=arm64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o proxy_${os}_arm64${suffix} feather-proxy/cmd
#	if $UPX; then upx -9 proxy_${os}_arm64${suffix};fi
done

#MIPS32LE
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o proxy_linux_mipsle feather-proxy/cmd
env CGO_ENABLED=0 GOOS=linux GOARCH=mips GOMIPS=softfloat go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o proxy_linux_mips feather-proxy/cmd
if $UPX; then upx -9 *;fi
exit 0