#!/bin/bash

# Скрипт для тестирования API news_portal
# Запуск: ./test_api.sh

# Загрузка конфигурации из .env файла
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/.env" ]; then
    set -a
    source "$SCRIPT_DIR/.env"
    set +a
fi

# Формирование BASE_URL из HTTP_ADDR (по умолчанию :8080)
HTTP_ADDR="${HTTP_ADDR:-:8080}"
BASE_URL="http://localhost${HTTP_ADDR}"
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
check_field "$body" "categories" "getCategories has Categories"
check_array_length "$body" "categories" 5 "getCategories returns 5 categories"

# -------------------------------------------
# Test 2: GET /v1/getTags
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getTags${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getTags")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "getTags returns 200"
check_field "$body" "tags" "getTags has Tags"
check_array_length "$body" "tags" 6 "getTags returns 6 tags"

# -------------------------------------------
# Test 3: GET /v1/countNews
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "countNews returns 200"
check_field "$body" "count" "countNews has Count"
check_value "$body" "count" "3" "countNews returns 3 (published news)"

# -------------------------------------------
# Test 4: GET /v1/countNews with filters
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews?category=5 (Наука)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews?category=5")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "countNews with category filter returns 200"
check_value "$body" "count" "1" "countNews with category=5 returns 1"

# -------------------------------------------
# Test 5: GET /v1/listNews
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews returns 200"
check_field "$body" "news" "listNews has News"
check_array_length "$body" "news" 3 "listNews returns 3 news"

# -------------------------------------------
# Test 6: GET /v1/listNews with category filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?category=1 (Технологии)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?category=1")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with category filter returns 200"
check_array_length "$body" "news" 1 "listNews with category=1 returns 1 news"
check_value "$body" "news[0].title" "Прорыв в области искусственного интеллекта" "listNews category filter returns correct news"

# -------------------------------------------
# Test 7: GET /v1/listNews with tag filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?tag=5 (Космос)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?tag=5")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with tag filter returns 200"
check_array_length "$body" "news" 1 "listNews with tag=5 returns 1 news"
check_value "$body" "news[0].title" "Открыта новая экзопланета" "listNews tag filter returns correct news"

# -------------------------------------------
# Test 8: GET /v1/listNews with pagination
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?page=1&limit=2${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?page=1&limit=2")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with pagination returns 200"
check_array_length "$body" "news" 2 "listNews with limit=2 returns 2 news"

# -------------------------------------------
# Test 9: GET /v1/listNews with date range (from)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?from=2026-03-01T00:00:00Z${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?from=2026-03-01T00:00:00Z")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with from filter returns 200"

# -------------------------------------------
# Test 10: GET /v1/listNews with date range (to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?to=2026-02-28T23:59:59Z${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?to=2026-02-28T23:59:59Z")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with to filter returns 200"

# -------------------------------------------
# Test 11: GET /v1/listNews with date range (from + to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?from=2026-02-01T00:00:00Z&to=2026-02-28T23:59:59Z${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?from=2026-02-01T00:00:00Z&to=2026-02-28T23:59:59Z")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with from+to filters returns 200"

# -------------------------------------------
# Test 12: GET /v1/listNews with combined filters (category + tag)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?category=5&tag=5 (Наука + Космос)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?category=5&tag=5")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with category+tag filters returns 200"
check_array_length "$body" "news" 1 "listNews with category=5&tag=5 returns 1 news"

# -------------------------------------------
# Test 13: GET /v1/countNews with tag filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews?tag=5 (Космос)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews?tag=5")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "countNews with tag filter returns 200"
check_value "$body" "count" "1" "countNews with tag=5 returns 1"

# -------------------------------------------
# Test 14: GET /v1/countNews with date range (from)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews?from=2026-03-01T00:00:00Z${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews?from=2026-03-01T00:00:00Z")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "countNews with from filter returns 200"

# -------------------------------------------
# Test 15: GET /v1/countNews with combined filters (category + tag)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews?category=1&tag=1 (Технологии + AI)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews?category=1&tag=1")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "countNews with category+tag filters returns 200"
check_value "$body" "count" "1" "countNews with category=1&tag=1 returns 1"

# -------------------------------------------
# Test 16: GET /v1/listNews with negative category ID
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?category=-1 (invalid category)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?category=-1")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 400 "$status" "listNews with category=-1 returns 400"
check_value "$body" "error" "invalid data" "listNews category error message"

# -------------------------------------------
# Test 17: GET /v1/listNews with invalid date range (from > to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?from=2026-12-31T00:00:00Z&to=2026-01-01T00:00:00Z (invalid range)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?from=2026-12-31T00:00:00Z&to=2026-01-01T00:00:00Z")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 400 "$status" "listNews with from > to returns 400"
check_value "$body" "error" "invalid data" "listNews date range error message"

# -------------------------------------------
# Test 18: GET /v1/countNews with negative category ID
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews?category=-5 (invalid category)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews?category=-5")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 400 "$status" "countNews with category=-5 returns 400"
check_value "$body" "error" "invalid data" "countNews category error message"

# -------------------------------------------
# Test 19: GET /v1/countNews with invalid date range (from > to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/countNews?from=2026-06-01T00:00:00Z&to=2026-01-01T00:00:00Z (invalid range)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews?from=2026-06-01T00:00:00Z&to=2026-01-01T00:00:00Z")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 400 "$status" "countNews with from > to returns 400"
check_value "$body" "error" "invalid data" "countNews date range error message"

# -------------------------------------------
# Test 20: GET /v1/listNews - no results for filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/listNews?category=999 (no results)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews?category=999")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "listNews with non-existent category returns 200"
check_array_length "$body" "news" 0 "listNews with category=999 returns empty array"

# -------------------------------------------
# Test 21: GET /v1/getNews
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getNews?id=1${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getNews?id=1")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "getNews returns 200"
check_field "$body" "news" "getNews has News"
check_value "$body" "news.id" "1" "getNews returns correct ID"
check_value "$body" "news.title" "Прорыв в области искусственного интеллекта" "getNews returns correct title"
check_field "$body" "news.content" "getNews includes Content"
check_field "$body" "news.preamble" "getNews includes Preamble"
check_field "$body" "news.category" "getNews includes Category"
check_field "$body" "news.tags" "getNews includes Tags"

# -------------------------------------------
# Test 22: GET /v1/getNews - not found
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getNews?id=999 (not found)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getNews?id=999")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 404 "$status" "getNews with invalid id returns 404"
check_field "$body" "error" "getNews error response has error field"

# -------------------------------------------
# Test 23: GET /v1/getNews - missing id
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: GET /v1/getNews (missing id)${NC}"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/getNews")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 400 "$status" "getNews without id returns 400"
check_value "$body" "error" "invalid data" "getNews error message is 'invalid data'"

# -------------------------------------------
# Test 24: GET /v1/countNews vs /v1/listNews consistency
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: CountNews equals ListNews length${NC}"
count_response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/countNews")
count_status=$(echo "$count_response" | tail -n1)
count_body=$(echo "$count_response" | sed '$d')
count_value=$(echo "$count_body" | jq -r ".count")

list_response=$(curl -s -w "\n%{http_code}" "$BASE_URL/v1/listNews")
list_status=$(echo "$list_response" | tail -n1)
list_body=$(echo "$list_response" | sed '$d')
list_length=$(echo "$list_body" | jq ".news | length")

check_status 200 "$count_status" "countNews returns 200"
check_status 200 "$list_status" "listNews returns 200"

if [ "$count_value" -eq "$list_length" ]; then
    echo -e "${GREEN}✓ PASS${NC}: CountNews ($count_value) equals ListNews length ($list_length)"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL${NC}: CountNews ($count_value) != ListNews length ($list_length)"
    ((FAILED++))
fi

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
