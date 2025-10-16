
curl -X POST "http://elasticsearch:9200/_security/role/logstash_writer" \
  -H "Content-Type: application/json" \
  -u "elastic:${ELASTIC_PASSWORD}" \
  -d '{
    "cluster": ["monitor","manage_ilm","manage_index_templates"],
    "indices": [
      {
        "names": ["*"],
        "privileges": ["auto_configure","create_index","write","create","view_index_metadata"]
      }
    ]
  }'

# create role property_writer 
curl -X POST "http://elasticsearch:9200/_security/role/property_writer" \
  -H "Content-Type: application/json" \
  -u "elastic:${ELASTIC_PASSWORD}" \
  -d '{
    "cluster": ["all"],
    "indices": [
      {
        "names": ["properties*"],
        "privileges": ["all"]
      }
    ]
  }'

# create user logstash_service_account
curl -X POST "http://elasticsearch:9200/_security/user/logstash_service_account" \
  -H "Content-Type: application/json" \
  -u "elastic:${ELASTIC_PASSWORD}" \
  -d '{
    "password": "'"${LOGSTASH_HOST_SERVICE_ACCOUNT_PASSWORD}"'",
    "full_name": "Logstash Service Account",
    "email": "logstash_service_account@example.com",
    "metadata": {
      "intake": "logstash"
    },
    "roles": ["logstash_writer", "property_writer"]
  }'

psql -U postgres -h database -d postgres -c "CREATE DATABASE socialapp;"
psql -U postgres -h database -d socialapp -f /schema.sql
