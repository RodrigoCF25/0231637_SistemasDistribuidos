# Etapa de construcción
FROM golang:1.22.5 AS builder

# Configura el directorio de trabajo
WORKDIR /app

# Copia los archivos del proyecto
COPY . .

# Descarga las dependencias y construye el proyecto
RUN go mod download

# Ejecuta `go test` en el contenedor cuando se inicie
CMD ["go", "test", "/app/server", "-v"]

