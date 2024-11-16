#!/bin/bash

# Construir la imagen Docker
#echo "Construyendo la imagen Docker..."
#docker build -t avance-proyecto-computo-distribuido:test .

# Ejecutar el contenedor en segundo plano
#echo "Ejecutando el contenedor..."
#docker run --rm -it -d --name Prueba avance-proyecto-computo-distribuido:test sh

# Crear un clúster de Kind
echo "Creando el clúster de Kind..."
kind create cluster --name mycluster

# Verificar que el clúster se haya creado
echo "Verificando clúster..."
kind get clusters

# Cargar la imagen Docker en el clúster Kind
echo "Cargando la imagen Docker en el clúster de Kind..."
kind load docker-image avance-proyecto-computo-distribuido:test --name mycluster

# Verificar que el clúster y la imagen se hayan cargado correctamente
echo "Verificando nodos del clúster..."
kubectl get nodes

# Navegar a la carpeta del proyecto
#cd ruta/a/la/carpeta/del/proyecto || exit

# Verificar si el chart 'mychart' ya existe
if [ -d "mychart" ]; then
  echo "El chart 'mychart' ya existe. Continuando..."
else
  # Crear un nuevo chart de Helm si no existe
  echo "Creando un nuevo Chart de Helm..."
  helm create mychart
fi

# Entrar en la carpeta mychart
cd mychart || exit

# Modificar los archivos YAML si no se ha hecho
#echo "Recuerda modificar los archivos 'values.yaml' y 'templates/deployment.yaml' según las instrucciones."

# Instalar el chart de Helm
echo "Instalando el chart de Helm..."
helm install my-test-deployment .

# Proporcionar instrucciones para limpiar el entorno
echo "Para limpiar el entorno, ejecuta los siguientes comandos:"
echo "helm uninstall my-test-deployment"
echo "kind delete cluster --name mycluster"
