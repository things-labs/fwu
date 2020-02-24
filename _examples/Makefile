
# 固件版本
version = v0.1.2
# 设置固件名称
firmwareName = anytool

ifeq (${MAKECMDGOALS},arm5)
	firmware := ${firmwareName}-linux-arm5-${version}
else ifeq (${MAKECMDGOALS},mipsle)
	firmware := ${firmwareName}-linux-mipsle-${version}
else ifeq (${MAKECMDGOALS},windows)
	firmware := ${firmwareName}-windows-amd64-${version}
else # default system
	firmware := ${firmwareName}-system-${version}
endif

# 路径相关
BINDIR=${CURDIR}/bin

system:
	@echo "----> system executable building..."
	@mkdir -p ${BINDIR}
	@CGO_ENABLED=0 go build -ldflags "-s -w" -o ${BINDIR}/${firmware} .
	@upx -9 ${BINDIR}/${firmware}
	@bzip2 -c ${BINDIR}/${firmware} > ${BINDIR}/${firmware}.bz2
	@echo "----> system executable build successful"

arm5:
	@echo "----> linux arm5 executable building..."
	@mkdir -p ${BINDIR}
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 CC=arm-linux-gcc \
	go build -ldflags "-s -w" -o ${BINDIR}/${firmware} .
	@upx -9 ${BINDIR}/${firmware}
	@bzip2 -c ${BINDIR}/${firmware} > ${BINDIR}/${firmware}.bz2
	@echo "----> linux arm5 executable build successful"

mipsle:
	@echo "----> linux mipsle executable building..."
	@mkdir -p ${BINDIR}
	@CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOARM=5 \
	STAGING_DIR=/opt/toolchain/openwrt18.06/staging_dir \
	CC=/opt/toolchain/openwrt18.06/staging_dir/gcc-mipsel-linux-7.3.0/bin/mipsel-openwrt-linux-gcc \
	go build -ldflags "-s -w" -o ${BINDIR}/${firmware} .
	@upx -9 ${BINDIR}/${firmware}
	@bzip2 -c ${BINDIR}/${firmware} > ${BINDIR}/${firmware}.bz2
	@echo "----> linux mipsle executable build successful"

windows:
	@echo "----> windows amd64 executable building..."
	@mkdir -p ${BINDIR}
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${BINDIR}/${firmware} .
	@upx -9 ${BINDIR}/${firmware}
	@bzip2 -c ${BINDIR}/${firmware} > ${BINDIR}/${firmware}.bz2
	@echo "----> windows amd64 executable build successful"

install:
	@go install

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
	@rm -rf ${BINDIR}/*
	@echo "----> clean successful"

.PHONY: system arm5 mipsle install clean



