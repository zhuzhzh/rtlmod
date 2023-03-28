EXEC = ./bin/vmod
SRCS := ./cmd/vmod/main.go 

echo:
	@echo ${SRCS}

b build:
	go build -gcflags "-N -l" -o ${EXEC} ${SRCS}

release rel:
	go build -o ${EXEC} ${SRCS}

r run: 
	${EXEC} -c tests/config.json -f tests/filelist.f -o out


.PHONY: b r rel
