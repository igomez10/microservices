for file in $(find prometheus -type f)
do
    filename=$(basename "$file")
    configmap_name="prometheus-config-${filename}"
    kubectl create configmap "$configmap_name" --from-file="$file"
    kubectl label configmap $configmap_name app=prometheus
done
