# cURL Examples for API Testing

## Health Check Endpoint

### Basic Health Check
```bash
curl http://localhost:8080/api/v1/health
```

**Expected Response:**
```
OK
```

**Expected Status Code:** `200 OK`

### Health Check with Verbose Output
```bash
curl -v http://localhost:8080/api/v1/health
```

### Health Check with Headers Only
```bash
curl -I http://localhost:8080/api/v1/health
```

### Health Check (Pretty Print)
```bash
curl -s http://localhost:8080/api/v1/health && echo ""
```

### Health Check for Production/Docker
```bash
# If running in Docker
curl http://localhost:8080/api/v1/health

# If deployed to AWS/Production
curl https://your-domain.com/api/v1/health
```

## Other Common Endpoints

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@campusambassador.com",
    "password": "password123"
  }'
```

### Get Current User Profile (Requires Auth)
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Get All Colleges
```bash
curl http://localhost:8080/api/v1/colleges
```

### Get All States
```bash
curl http://localhost:8080/api/v1/states
```

## Quick Test Script

Save this as `test-health.sh`:

```bash
#!/bin/bash

ENDPOINT="${1:-http://localhost:8080/api/v1/health}"

echo "Testing health endpoint: $ENDPOINT"
echo ""

response=$(curl -s -w "\n%{http_code}" "$ENDPOINT")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" -eq 200 ]; then
    echo "✅ Health check passed!"
    echo "Response: $body"
    exit 0
else
    echo "❌ Health check failed!"
    echo "HTTP Code: $http_code"
    echo "Response: $body"
    exit 1
fi
```

Usage:
```bash
chmod +x test-health.sh
./test-health.sh
./test-health.sh https://your-production-domain.com/api/v1/health
```
