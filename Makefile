# Define la ruta del archivo proto
PROTO_FILES=api/v1/*.proto

# Regla para compilar los archivos .proto
compile:
	protoc $(PROTO_FILES) \
	--go_out=. \
	--go_opt=paths=source_relative \
	--proto_path=.

# Regla para correr las pruebas
test:
	go test -race ./...
