CGO_CFLAGS = -I/usr/include/SDL3
TAGS = sdl3

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" go build -tags $(TAGS) -o raylib_test .

run: build
	./raylib_test

