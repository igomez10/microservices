
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
