#!/usr/bin/env bash

manifest=$(cat << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 4
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
EOF)

INPUT_ACCOUNTS="kubernetes-eks-prod" \
INPUT_APPLICATION="front50" \
INPUT_MANIFESTPATH="hack/manifests" \
INPUT_CLOUDPROVIDER="kubernetes" \
INPUT_BASEURL="https://spinnaker-api.armory.io" \
INPUT_CONFIGPATH="/Users/ethanfrogers/.spin/config" \
INPUT_WAIT="true" \
go run main.go