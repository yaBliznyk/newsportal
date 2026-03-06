#!/bin/bash

# Скрипт для тестирования API news_portal
# Запуск: ./test_api.sh

BASE_URL="http://localhost:8080"
FAILED=0
PASSED=0

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для проверки статус кода
check_status() {
    local expected=$1
    local actual=$2
    local test_name=$3

    if [ "$actual" -eq "$expected" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name (status: $actual)"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name (expected: $expected, got: $actual)"
        ((FAILED++))
    fi
}

# Функция для проверки наличия поля в JSON
check_field() {
    local response=$1
    local field=$2
    local test_name=$3

    if echo "$response" | jq -e ".$field" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name (field '$field' exists)"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name (field '$field' missing)"
        ((FAILED++))
    fi
}

# Функция для проверки значения поля
check_value() {
    local response=$1
    local field=$2
    local expected=$3
    local test_name=$4

    actual=$(echo "$response" | jq -r ".$field")
    if [ "$actual" = "$expected" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name ($field = $actual)"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name (expected '$expected', got '$actual')"
        ((FAILED++))
    fi
}

# Функция для проверки длины массива
check_array_length() {
    local response=$1
    local field=$2
    local expected=$3
    local test_name=$4

    actual=$(echo "$response" | jq ".$field | length")
    if [ "$actual" -eq "$expected" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name ($field length = $actual)"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name (expected $expected, got $actual)"
        ((FAILED++))
    fi
}

echo -e "${YELLOW}======================================${NC}"
echo -e "${YELLOW}    Testing News Portal API${NC}"
echo -e "${YELLOW}======================================${NC}"
echo ""

# -------------------------------------------
# Test 1: GET /v1/getCategories
# -------------------------------------------
echo -e "${YELLOW}Test: GET /v1/getCategories${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getCategories")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "getCategories returns 200"
check_field "$body" "Categories" "getCategories has Categories"
check_array_length "$body" "Categories" 5 "getCategories returns 5 categories"

# -------------------------------------------
# Test 2: GET /v1/getTags
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getTags${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getTags")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "getTags returns 200"
check_field "$body" "Tags" "getTags has Tags"
check_array_length "$body" "Tags" 6 "getTags returns 6 tags"

# -------------------------------------------
# Test 3: GET /v1/countNews
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "countNews returns 200"
check_field "$body" "Count" "countNews has Count"
check_value "$body" "Count" "3" "countNews returns 3 (published news)"

# -------------------------------------------
# Test 4: GET /v1/countNews with filters
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews?category=5 (Наука)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews?category=5")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "countNews with category filter returns 200"
check_value "$body" "Count" "1" "countNews with category=5 returns 1"

# -------------------------------------------
# Test 5: GET /v1/listNews
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews returns 200"
check_field "$body" "News" "listNews has News"
check_array_length "$body" "News" 3 "listNews returns 3 news"

# -------------------------------------------
# Test 6: GET /v1/listNews with category filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?category=1 (Технологии)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?category=1")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with category filter returns 200"
check_array_length "$body" "News" 1 "listNews with category=1 returns 1 news"
check_value "$body" "News[0].Title" "Прорыв в области искусственного интеллекта" "listNews category filter returns correct news"

# -------------------------------------------
# Test 7: GET /v1/listNews with tag filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?tag=5 (Космос)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?tag=5")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with tag filter returns 200"
check_array_length "$body" "News" 1 "listNews with tag=5 returns 1 news"
check_value "$body" "News[0].Title" "Открыта новая экзопланета" "listNews tag filter returns correct news"

# -------------------------------------------
# Test 8: GET /v1/listNews with pagination
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?page=1&limit=2${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?page=1&limit=2")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with pagination returns 200"
check_array_length "$body" "News" 2 "listNews with limit=2 returns 2 news"

# -------------------------------------------
# Test 9: GET /v1/getNews
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getNews?id=1${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getNews?id=1")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "getNews returns 200"
check_field "$body" "News" "getNews has News"
check_value "$body" "News.ID" "1" "getNews returns correct ID"
check_value "$body" "News.Title" "Прорыв в области искусственного интеллекта" "getNews returns correct title"
check_field "$body" "News.Content" "getNews includes Content"
check_field "$body" "News.Preamble" "getNews includes Preamble"
check_field "$body" "News.Category" "getNews includes Category"
check_field "$body" "News.Tags" "getNews includes Tags"

# -------------------------------------------
# Test 10: GET /v1/getNews - not found
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getNews?id=999 (not found)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getNews?id=999")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 404 "$status" "getNews with invalid id returns 404"
check_field "$body" "error" "getNews error response has error field"

# -------------------------------------------
# Test 11: GET /v1/getNews - missing id
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getNews (missing id)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getNews")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 400 "$status" "getNews without id returns 400"
check_value "$body" "error" "invalid data" "getNews error message is 'invalid data'"

# -------------------------------------------
# Summary
# -------------------------------------------
echo ""
echo -e "${YELLOW}======================================${NC}"
echo -e "${YELLOW}              Summary${NC}"
echo -e "${YELLOW}======================================${NC}"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi
