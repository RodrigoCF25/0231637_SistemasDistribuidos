docker build -t avance-proyecto-computo-distribuido:test .

docker run --rm -it -d --name Prueba avance-proyecto-computo-distribuido:test sh


kind create cluster --name mycluster

Verifica que el cluster ya se creó:
kind get clusters

kind load docker-image avance-proyecto-computo-distribuido:test --name mycluster


Verifica que el clúster y la imagen se hayan cargado correctamente:
kubectl get nodes


Ir a la carpeta específica del proyecto
helm create mychart

cd mychart

Modificar los archivos:

values.yaml para que use la imagen creada
# mychart/values.yaml
image:
  repository: <imageName>
  tag: <tagDeLaImagen>
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80


Edita el archivo templates/deployment.yaml para configurar un pod que ejecute go test y registre la salida de los tests.

# mychart/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}-test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ .Release.Name }}
  template:
    metadata:
      labels:
        app: {{ .Release.Name }}
    spec:
      containers:
      //ESTO COPIARLO
        - name: server-test
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          command: ["go", "test", "-v", "./..."]
          imagePullPolicy: {{ .Values.image.pullPolicy }}



helm install my-test-deployment ./mychart


Para limpiar el entorno, puedes desinstalar el Chart de Helm y eliminar el clúster de Kind:

helm uninstall my-test-deployment
kind delete cluster --name mycluster
