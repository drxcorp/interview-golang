#!/bin/bash

# E-Commerce API Test Script
# Base URL
BASE_URL="http://localhost:8080"

echo "=========================================="
echo "E-Commerce API Test Commands"
echo "=========================================="

echo ""
echo "### Health Check ###"
echo ""

# Health Check
echo "1. Health Check"
curl -X GET "$BASE_URL/health"
echo -e "\n"

echo ""
echo "### Users ###"
echo ""

# Get all users
echo "2. Get all users"
curl -X GET "$BASE_URL/users"
echo -e "\n"

# Search users
echo "3. Search users by name"
curl -X GET "$BASE_URL/users?search=John"
echo -e "\n"

# Create user
echo "4. Create a new user"
curl -X POST "$BASE_URL/users/create" \
  -d "name=Jane Doe" \
  -d "email=jane@example.com" \
  -d "password=password123" \
  -d "role=user"
echo -e "\n"

# Delete user (use with caution)
echo "5. Delete user (ID=2)"
curl -X DELETE "$BASE_URL/users/delete?id=2"
echo -e "\n"

echo ""
echo "### Products ###"
echo ""

# Get all products
echo "6. Get all products"
curl -X GET "$BASE_URL/products"
echo -e "\n"

echo ""
echo "### Orders ###"
echo ""

# Get all orders
echo "7. Get all orders"
curl -X GET "$BASE_URL/orders"
echo -e "\n"

# Get orders by user
echo "8. Get orders by user ID"
curl -X GET "$BASE_URL/orders?user_id=1"
echo -e "\n"

# Create order
echo "9. Create a new order"
curl -X POST "$BASE_URL/orders/create" \
  -d "user_id=1" \
  -d "product_id=1" \
  -d "quantity=2"
echo -e "\n"

echo ""
echo "### Wallet ###"
echo ""

# Get wallet by user
echo "10. Get wallet by user ID"
curl -X GET "$BASE_URL/wallet?user_id=1"
echo -e "\n"

# Get wallet balance
echo "11. Get wallet balance"
curl -X GET "$BASE_URL/wallet/balance?user_id=1"
echo -e "\n"

# Create wallet for user
echo "12. Create wallet for user"
curl -X POST "$BASE_URL/wallet/create" \
  -d "user_id=2" \
  -d "initial_balance=500" \
  -d "currency=USD"
echo -e "\n"

# Top up wallet
echo "13. Top up wallet"
curl -X POST "$BASE_URL/wallet/topup" \
  -d "user_id=1" \
  -d "amount=100"
echo -e "\n"

# Transfer money between users
echo "14. Transfer money (user 1 to user 2)"
curl -X POST "$BASE_URL/wallet/transfer" \
  -H "Content-Type: application/json" \
  -d '{"from_user_id": 1, "to_user_id": 2, "amount": 50, "description": "Test transfer"}'
echo -e "\n"

# Pay for order
echo "15. Pay for order"
curl -X POST "$BASE_URL/wallet/pay" \
  -d "user_id=1" \
  -d "order_id=1"
echo -e "\n"

# Get transaction history
echo "16. Get transaction history"
curl -X GET "$BASE_URL/wallet/transactions?user_id=1"
echo -e "\n"

echo "=========================================="
echo "Test Complete"
echo "=========================================="
