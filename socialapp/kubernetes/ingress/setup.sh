kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.0/deploy/static/provider/cloud/deploy.yaml
# openssl genrsa -out server.key 2048
# for cert at website example.com
openssl req -new -x509 -sha256 -key server.key -out asteriskserver.crt \
-days 3650 -subj "/CN=*.example.com" -addext "subjectAltName=DNS:*.example.com,DNS:www.*.example.com"

openssl req -new -x509 -sha256 -key server.key -out prometheusserver.crt \
-days 3650 -subj "/CN=prometheus.example.com" -addext "subjectAltName=DNS:prometheus.example.com,DNS:www.prometheus.example.com"

openssl req -new -x509 -sha256 -key server.key -out grafana.crt \
-days 3650 -subj "/CN=grafana.example.com" -addext "subjectAltName=DNS:grafana.example.com,DNS:www.grafana.example.com"

# delete secrets first
kubectl delete secret prometheus-cert
kubectl delete secret grafana-cert
kubectl delete secret asterisk-cert

# create secrets
kubectl create secret tls prometheus-cert --key server.key --cert prometheusserver.crt -n default
kubectl create secret tls grafana-cert --key server.key --cert grafana.crt -n default
kubectl create secret tls asterisk-cert --key server.key --cert asteriskserver.crt -n default
