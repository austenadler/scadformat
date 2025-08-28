scadformat: test cmd/*.go internal/parser/*.go internal/formatter/*.go internal/logutil/*.go
	go build cmd/scadformat.go

internal/parser/*.go: OpenSCAD.g4
	go generate ./...

test: internal/parser/*.go
	go test ./...

wasm: cmd/*.go internal/parser/*.go internal/formatter/*.go internal/logutil/*.go
	GOOS=js GOARCH=wasm go build -o scadformat.wasm cmd/wasm/main_wasm_js.go

clean:
	rm -f internal/parser/*
	rm -rf internal/parser/.antlr
	rm -f scadformat.wasm
	rm cmd/version.txt

.PHONY: clean test
