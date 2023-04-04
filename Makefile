EXEC = ./bin/vmod
SRCS := ./cmd/vmod/main.go 
TOFILE := 1
VERB := info
ifeq (1, ${TOFILE})
	REDIRECT := --tofile
else
	REDIRECT :=
endif

echo:
	@echo ${SRCS}

b build:
	go build -gcflags "-N -l" -o ${EXEC} ${SRCS}

release rel:
	go build -o ${EXEC} ${SRCS}

r run: 
	${EXEC} chain -c test/config.json -f test/filelist.f -o out --verbose ${VERB} ${REDIRECT} test/lib1.v


.PHONY: b r rel
