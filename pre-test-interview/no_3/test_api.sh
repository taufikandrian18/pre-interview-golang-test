#!/bin/bash

# Cache API Testing Script
# This script demonstrates how to test the cache API endpoints

echo "🚀 Testing Cache API Server"
echo "Make sure the server is running with: go run . server"
echo ""

BASE_URL="http://localhost:8080"

echo "📝 Testing SET operations:"
echo "Setting user name..."
curl -s -X POST "$BASE_URL/api/cache/set?key=user" \
     -H "Content-Type: application/json" \
     -d '"Alice Johnson"' | jq '.'

echo ""
echo "Setting user age..."
curl -s -X POST "$BASE_URL/api/cache/set?key=age" \
     -d "value=28" | jq '.'

echo ""
echo "Setting user preferences (JSON)..."
curl -s -X POST "$BASE_URL/api/cache/set?key=preferences" \
     -H "Content-Type: application/json" \
     -d '{"theme":"dark","notifications":true}' | jq '.'

echo ""
echo "🔍 Testing GET operations:"
echo "Getting user name..."
curl -s "$BASE_URL/api/cache/get?key=user" | jq '.'

echo ""
echo "Getting user age..."
curl -s "$BASE_URL/api/cache/get?key=age" | jq '.'

echo ""
echo "Getting preferences..."
curl -s "$BASE_URL/api/cache/get?key=preferences" | jq '.'

echo ""
echo "Getting non-existent key..."
curl -s "$BASE_URL/api/cache/get?key=nonexistent" | jq '.'

echo ""
echo "📊 Testing STATS:"
curl -s "$BASE_URL/api/cache/stats" | jq '.'

echo ""
echo "🗑️  Testing DELETE operations:"
echo "Deleting age..."
curl -s -X DELETE "$BASE_URL/api/cache/delete?key=age" | jq '.'

echo ""
echo "Trying to get deleted age..."
curl -s "$BASE_URL/api/cache/get?key=age" | jq '.'

echo ""
echo "❤️  Testing HEALTH check:"
curl -s "$BASE_URL/health" | jq '.'

echo ""
echo "🎯 Testing TTL functionality (if using TTL cache):"
echo "Setting temporary data..."
curl -s -X POST "$BASE_URL/api/cache/set?key=temp" \
     -H "Content-Type: application/json" \
     -d '"This will expire in 30 seconds"' | jq '.'

echo ""
echo "Getting temporary data immediately..."
curl -s "$BASE_URL/api/cache/get?key=temp" | jq '.'

echo ""
echo "💡 Try these additional commands manually:"
echo "curl -X POST '$BASE_URL/api/cache/set?key=test' -d 'value=hello'"
echo "curl '$BASE_URL/api/cache/get?key=test'"
echo "curl -X DELETE '$BASE_URL/api/cache/delete?key=test'"
echo ""
echo "📖 Open http://localhost:8080 in your browser for interactive documentation!"
