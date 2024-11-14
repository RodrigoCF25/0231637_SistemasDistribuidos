# Usa la imagen base de Golang
FROM golang:1.22.5

# Configura el directorio de trabajo
WORKDIR /app

# Copia el proyecto completo
COPY . .

# Descarga las dependencias y construye el binario
RUN go mod download
RUN go build -o test-runner

# Ejecuta el binario
CMD ["./test-runner"]
