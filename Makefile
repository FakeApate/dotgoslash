BIN    = dotgoslash
SRC    = main.go scanner.go payload.go
OUT    = dist

define build
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) $(3) go build -o $(OUT)/$(BIN)-$(4) $(SRC)
endef

.PHONY: all clean

all: \
	$(OUT)/$(BIN)-linux-amd64 \
	$(OUT)/$(BIN)-linux-arm64 \
	$(OUT)/$(BIN)-linux-armv6 \
	$(OUT)/$(BIN)-linux-armv7 \
	$(OUT)/$(BIN)-darwin-amd64 \
	$(OUT)/$(BIN)-darwin-arm64 \
	$(OUT)/$(BIN)-windows-amd64.exe \
	$(OUT)/$(BIN)-windows-arm64.exe

$(OUT)/$(BIN)-linux-amd64:
	$(call build,linux,amd64,,linux-amd64)

$(OUT)/$(BIN)-linux-arm64:
	$(call build,linux,arm64,,linux-arm64)

# Pi 1, Zero (ARMv6)
$(OUT)/$(BIN)-linux-armv6:
	$(call build,linux,arm,GOARM=6,linux-armv6)

# Pi 2, Pi 3 32-bit OS (ARMv7)
$(OUT)/$(BIN)-linux-armv7:
	$(call build,linux,arm,GOARM=7,linux-armv7)

$(OUT)/$(BIN)-darwin-amd64:
	$(call build,darwin,amd64,,darwin-amd64)

$(OUT)/$(BIN)-darwin-arm64:
	$(call build,darwin,arm64,,darwin-arm64)

$(OUT)/$(BIN)-windows-amd64.exe:
	$(call build,windows,amd64,,windows-amd64.exe)

$(OUT)/$(BIN)-windows-arm64.exe:
	$(call build,windows,arm64,,windows-arm64.exe)

clean:
	rm -rf $(OUT)
