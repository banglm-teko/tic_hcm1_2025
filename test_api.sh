 #!/bin/bash

# Test script for the login API
API_URL="http://localhost:8080/api"

echo "=== Testing Login API ==="
echo

# Test successful login
echo "1. Testing successful login (admin/admin123):"
curl -X POST $API_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' \
  | jq .
echo
echo

# Test successful login with different user
echo "2. Testing successful login (user1/password1):"
curl -X POST $API_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user1", "password": "password1"}' \
  | jq .
echo
echo

# Test failed login
echo "3. Testing failed login (wrong password):"
curl -X POST $API_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "wrongpassword"}' \
  | jq .
echo
echo

# Test failed login with non-existent user
echo "4. Testing failed login (non-existent user):"
curl -X POST $API_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username": "nonexistent", "password": "anypassword"}' \
  | jq .
echo
echo

# Test health endpoint
echo "5. Testing health endpoint:"
curl -X GET $API_URL/health | jq .
echo
echo

echo "=== API Test Complete ==="