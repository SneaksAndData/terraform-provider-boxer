# Running the integration tests manually

1. Apply the terraform configuration
2. Run the following curl command to check to token
```bash
curl -X 'GET' \
  'http://localhost:8081/token/review' \
  -H "Authorization: Bearer $TOKEN"
```

