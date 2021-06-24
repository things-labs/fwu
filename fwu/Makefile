# 固件版本
version = v0.0.2

# 设置固件名称
firmwareName = fwu

# 编译平台
platform = CGO_ENABLED=0
ifeq (${MAKECMDGOALS},arm5)
	platform += GOOS=linux GOARCH=arm GOARM=5
	execveFile := ${firmwareName}-linux-arm5-${version}
else ifeq (${MAKECMDGOALS},arm7)
	platform += GOOS=linux GOARCH=arm GOARM=7
	execveFile := ${firmwareName}-linux-arm7-${version}
else ifeq (${MAKECMDGOALS},mipsle)
	platform += GOOS=linux GOARCH=mipsle GOARM=5 \
      STAGING_DIR=/opt/toolchain/openwrt18.06/staging_dir \
      CC=/opt/toolchain/openwrt18.06/staging_dir/gcc-mipsel-linux-7.3.0/bin/mipsel-openwrt-linux-gcc
	execveFile := ${firmwareName}-linux-mipsle-${version}
else ifeq (${MAKECMDGOALS},windows)
	platform += GOOS=windows GOARCH=amd64
	execveFile := ${firmwareName}-windows-amd64-${version}
else # default system
	execveFile := ${firmwareName}-${version}
endif

# 编译tags,多个采用','分开
opts = -trimpath -tags=noprotopack,noyamlpack,nomsgpack
# 编译flags
path = main
flags = -ldflags "-X '${path}.BuildTime=`date "+%F %T %z"`' \
 				  -X '${path}.Version=${version}' \
 				  -s -w"

target:
	@echo "----> executable building..."
	${platform} go build ${opts} ${flags} -o ${execveFile} .
	@upx --best --lzma ${execveFile}
	@bzip2 -c ${execveFile} > ${execveFile}.bz2
	@echo "----> executable build successful"

system: target
arm5: target
arm7: target
mipsle: target
windows: target

goreleaser:
	@goreleaser --snapshot --skip-publish --rm-dist

help:
	@echo " ------------- How to build ------------- "
	@echo " make         -- same as 'make system'"
	@echo " make system  -- build target for system"
	@echo " make arm5    -- build target for arm5"
	@echo " make mipsle  -- build target for mipsle"
	@echo " make windows -- build target for windows"
	@echo " make clean   -- clean build files"
	@echo " make install -- install executable files"
	@echo " ------------- How to build ------------- "

clean:
	@echo "----> cleanning"
	@go clean
	@echo "----> clean successful"

.PHONY: system arm5 mipsle install clean



