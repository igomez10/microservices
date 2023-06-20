# SOCIALAPP
curl -i -X POST -H "Accept:application/json" \
    -H  "Content-Type:application/json" \
    http://localhost:8083/connectors/ \
    -d @./connectors/socialapp-register-postgres.json

# PUTTYKNIFE
curl -i -X POST -H "Accept:application/json" \
-H  "Content-Type:application/json" \
http://localhost:8083/connectors/ \
-d @./debezium/config/puttyknife-register-postgres.json
