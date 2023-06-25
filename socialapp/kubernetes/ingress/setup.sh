kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.0/deploy/static/provider/cloud/deploy.yaml
kubectl delete secret selfsigned-cert
openssl genrsa -out server.key 2048
# for cert at website example.com
openssl req -new -x509 -sha256 -key server.key -out server.crt \
-days 3650 -subj "/CN=prometheus.example.com" -addext "subjectAltName=DNS:prometheus.example.com,DNS:www.prometheus.example.com"
kubectl create secret tls selfsigned-cert --key server.key --cert server.crt -n default
