# go-app-intern

go-app-intern is a web app written in Go that is dockerized, automated using jenkins and deployed on minikube.

Contains:
- Dockerfile
- dockercompose
- Jenkinsfile
- kubernetes files

## Prerequistis

- Docker
- docker compose
- Jenkins
- helm
- Minikube

This is the tree of the repository
```bash
.
├── Dockerfile
├── Jenkinsfile
├── README.md
├── app
│   ├── db.go
│   ├── go.mod
│   ├── go.sum
│   ├── main-app
│   └── main.go
├── charts
│   └── myapp
│       ├── Chart.yaml
│       ├── charts
│       │   └── mysql-8.8.16.tgz
│       ├── templates
│       │   ├── app-deployment.yaml
│       │   ├── app-service.yaml
│       │   ├── mysql-deployment.yaml
│       │   ├── mysql-service.yaml
│       │   ├── pv.yaml
│       │   ├── pvc.yaml
│       │   └── storageclass.yaml
│       └── values.yaml
├── docker-compose.yml
├── init.sql
└── main-app
```

## How to build this app

```bash
go build -o main-app .
```
This will produce an artifiact called main-app then to run it
```bash
./main-app
```
## Dockerize application
### Single-Stage Dockerfile

```bash
FROM golang:1.19

WORKDIR /main-app

COPY ./app . 
RUN GOOS=linux go build -o main-app .

EXPOSE 9090
CMD ["./main-app"]
```
To build image
```bash
docker build . -t app
```
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/d9ee8997-8d39-43fc-9e3a-d0b02d0ded6e)
The size of the image is more than 1 GB!!!
### Multi-Stage Dockerfile
As we saw single stage Dockerfile generated an image that has very big size and that isn`t effiecient
```bash
FROM golang:1.19 AS build-stage

WORKDIR /app

COPY ./app .

RUN GOOS=linux go build -o main-app .

# Deploy the application binary into a lean image
FROM debian:11-slim AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app .

#Create a non-root user
RUN adduser --disabled-password --gecos "" appuser
USER appuser

EXPOSE 9090 
CMD ["./main-app"]
```
To build image
```bash
docker build . -t ahmedelmelegy3570/app-multistage
```
What we benifted from that?
The size of the image is greatly reduced
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/1d839a6b-504b-4ac6-92ed-021e78978b1e)
Now it is less than 0.1 GB!

To run container with the built image, we will map port 9090 from host to port 9090 in the container
```bash
docker run -p 9090:9090 ahmedelmelegy3570/app-multistage
```
### To make the image more secure
to create a user and use it instead of root user
```bash
RUN adduser --disabled-password --gecos "" appuser
USER appuser
```
## Docker-compose
in dockercompose.yaml, I combined both app and mysql db so they have the same network and I made the app depend on mysql
```bash
version: '3'
services:
  app:
    image: app-multistage
    ports:
      - 9090:9090
    depends_on:
      - mysql
  mysql:
    image: mysql:5.7
    environment:
      - MYSQL_ROOT_PASSWORD=1234
      - MYSQL_DATABASE=mydb
      - MYSQL_USER=ahmed
      - MYSQL_PASSWORD=1234
    volumes:
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - 3306:3306
    command: ["--log-error-verbosity=3"]
```
To apply this docker compose file
```bash
docker-compose up
```
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/cf3d91b0-adf3-4407-a804-459fe8517aa1)
## Jenkins
To install jenkins as container and it will run on port 8080
```bash
docker run -p 8080:8080 -p 50000:50000 -d -v /var/run/docker.sock:/var/run/docker.sock -v jenkins_home:/var/jenkins_home jenkins/jenkins:lts
```

![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/985109ab-f079-4a7d-98ac-3151cc6e3f3a)
and this is the configuration of it 
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/0269e7ea-edba-4576-aecc-ff52c5c4a3f6)
and I created credentials to be able to securly login to dockerhub to push image
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/8cc4382a-7426-47a8-9572-a8fb2474e969)
This is the pipeline output after build it 
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/331a87a3-9af8-40d8-b92a-22c8e5fd20fd)

### Push image to dockerhub
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/280b1103-13eb-40f5-8163-98cdb87ebfb5)
## kubernetes and helm
app-deployment.yaml:
Added support for multiple replicas using a configurable value.
Added volume mounts and persistent volume claim (PVC) for volume persistence
```bash
spec:
  replicas: {{ .Values.replicaCount }}
    spec:
    {{- if .Values.persistence.enabled }}
        volumeMounts:
        - name: app-persistent-storage
          mountPath: /path/to/persistent/storage
    {{- end }}
  {{- if .Values.persistence.enabled }}
    volumes:
    - name: app-persistent-storage
      persistentVolumeClaim:
        claimName: app-pv-claim
  {{- end }}
```
app-service.yaml:
service type LoadBalancer to expose the service publicly.
```bash
apiVersion: v1
kind: Service
metadata:
  name: app
spec:
  selector:
    app: app
  ports:
  - name: http
    port: 9090
    targetPort: 9090
  type: LoadBalancer
```
Created values.yaml:
Added configuration values for replica count, autoscaling, and persistence.
```bash
replicaCount: 3

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 5
  targetCPUUtilizationPercentage: 50

persistence:
  enabled: true
  size: 10Gi
```
Created autoscaling.yaml:
Added an autoscaling manifest to scale the number of replicas based on CPU utilization.
```bash
apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  name: app-autoscaler
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: app
  minReplicas: {{ .Values.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.autoscaling.maxReplicas }}
  targetCPUUtilizationPercentage: {{ .Values.autoscaling.targetCPUUtilizationPercentage }}
```
Created app-pvc.yaml:
Added a persistent volume claim definition for the app.
```bash
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: app-pv-claim
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Values.persistence.size }}
  storageClassName: local-storage
```
Created storageclass.yaml:
Defined a storage class for local storage.
```bash
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
```
To deploy this app using helm run this command but before that you must need minikube cluster up and running
```bash
helm install myapp .
```
To browse the app
localhost:9090
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/1fd5fd88-35c0-4f4b-92be-49b54ea7e1e1)
localhost:9090/healthcheck
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/62e8842d-2277-47e0-9e2a-fa1f79776241)
I tried to figure out where is the problem in the api but I couldn`t. 
I tried to debug the code so I put log.Println() in each function but I couldn`t find the problem.
The intenrship database and table named stuff is created so there is no problem in the connection between the app and mysql
```bash
mysql -h 127.0.0.1 -P 3306 -u ahmed -p
Enter password: 1234
mysql> use internship
mysql> SHOW TABLES;
```
![image](https://github.com/gAhmed-Elmelegy/go-app-intern/assets/136341359/01c3a5dc-4f8d-484e-a484-6d43cd10154f)
