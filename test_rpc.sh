#!/bin/bash

# Скрипт для тестирования RPC API news_portal
# Запуск: ./test_rpc.sh

# Загрузка конфигурации из .env файла
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [ -f "$SCRIPT_DIR/.env" ]; then
    set -a
    source "$SCRIPT_DIR/.env"
    set +a
fi

# Формирование BASE_URL из HTTP_ADDR (по умолчанию :8080)
HTTP_ADDR="${HTTP_ADDR:-:8080}"
BASE_URL="http://localhost${HTTP_ADDR}/rpc/"
FAILED=0
PASSED=0

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для отправки JSON-RPC запроса
rpc_call() {
    local method=$1
    local params=$2
    local id=${3:-1}

    curl -s -w "\n%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"method\":\"$method\",\"params\":$params,\"id\":$id}" \
        "$BASE_URL"
}

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

# Функция для проверки отсутствия ошибки в JSON-RPC ответе
check_no_error() {
    local response=$1
    local test_name=$2

    if echo "$response" | jq -e ".error" > /dev/null 2>&1; then
        echo -e "${RED}✗ FAIL${NC}: $test_name (unexpected error: $(echo "$response" | jq -r '.error.message'))"
        ((FAILED++))
    else
        echo -e "${GREEN}✓ PASS${NC}: $test_name (no error)"
        ((PASSED++))
    fi
}

# Функция для проверки наличия ошибки в JSON-RPC ответе
check_error() {
    local response=$1
    local test_name=$2

    if echo "$response" | jq -e ".error" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name (error present)"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name (expected error, but got success)"
        ((FAILED++))
    fi
}

# Функция для проверки null результата (объект не найден)
check_null_result() {
    local response=$1
    local test_name=$2

    actual=$(echo "$response" | jq -r ".result")
    if [ "$actual" = "null" ]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name (result is null)"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name (expected null, got '$actual')"
        ((FAILED++))
    fi
}

echo -e "${YELLOW}======================================${NC}"
echo -e "${YELLOW}    Testing News Portal RPC API${NC}"
echo -e "${YELLOW}======================================${NC}"
echo ""

# -------------------------------------------
# Test 1: news.categories
# -------------------------------------------
echo -e "${YELLOW}Test: news.categories${NC}"
response=$(rpc_call "news.categories" "{}")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "categories returns 200"
check_no_error "$body" "categories has no error"
check_field "$body" "result" "categories has result"
check_array_length "$body" "result" 5 "categories returns 5 categories"

# -------------------------------------------
# Test 2: news.tags
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.tags${NC}"
response=$(rpc_call "news.tags" "{}")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "tags returns 200"
check_no_error "$body" "tags has no error"
check_field "$body" "result" "tags has result"
check_array_length "$body" "result" 6 "tags returns 6 tags"

# -------------------------------------------
# Test 3: news.count
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count${NC}"
response=$(rpc_call "news.count" "{}")
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "count returns 200"
check_no_error "$body" "count has no error"
check_value "$body" "result" "3" "count returns 3 (published news)"

# -------------------------------------------
# Test 4: news.count with category filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count with category=5 (Наука)${NC}"
response=$(rpc_call "news.count" '{"filter":{"CategoryID":5}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "count with category filter returns 200"
check_no_error "$body" "count with category has no error"
check_value "$body" "result" "1" "count with category=5 returns 1"

# -------------------------------------------
# Test 5: news.list
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list${NC}"
response=$(rpc_call "news.list" '{}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list returns 200"
check_no_error "$body" "list has no error"
check_field "$body" "result" "list has result"
check_array_length "$body" "result" 3 "list returns 3 news"

# -------------------------------------------
# Test 6: news.list with category filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with category=1 (Технологии)${NC}"
response=$(rpc_call "news.list" '{"filter":{"CategoryID":1}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with category filter returns 200"
check_no_error "$body" "list with category has no error"
check_array_length "$body" "result" 1 "list with category=1 returns 1 news"
check_value "$body" "result[0].title" "Прорыв в области искусственного интеллекта" "list category filter returns correct news"

# -------------------------------------------
# Test 7: news.list with tag filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with tag=5 (Космос)${NC}"
response=$(rpc_call "news.list" '{"filter":{"TagID":5}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with tag filter returns 200"
check_no_error "$body" "list with tag has no error"
check_array_length "$body" "result" 1 "list with tag=5 returns 1 news"
check_value "$body" "result[0].title" "Открыта новая экзопланета" "list tag filter returns correct news"

# -------------------------------------------
# Test 8: news.list with pagination
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with pager page=1 limit=2${NC}"
response=$(rpc_call "news.list" '{"pager":{"Page":1,"Limit":2}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with pagination returns 200"
check_no_error "$body" "list with pagination has no error"
check_array_length "$body" "result" 2 "list with limit=2 returns 2 news"

# -------------------------------------------
# Test 9: news.list with date range (from)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with From=2026-03-01T00:00:00Z${NC}"
response=$(rpc_call "news.list" '{"filter":{"From":"2026-03-01T00:00:00Z"}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with From filter returns 200"
check_no_error "$body" "list with From has no error"

# -------------------------------------------
# Test 10: news.list with date range (to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with To=2026-02-28T23:59:59Z${NC}"
response=$(rpc_call "news.list" '{"filter":{"To":"2026-02-28T23:59:59Z"}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with To filter returns 200"
check_no_error "$body" "list with To has no error"

# -------------------------------------------
# Test 11: news.list with date range (from + to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with From=2026-02-01T00:00:00Z&To=2026-02-28T23:59:59Z${NC}"
response=$(rpc_call "news.list" '{"filter":{"From":"2026-02-01T00:00:00Z","To":"2026-02-28T23:59:59Z"}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with From+To filters returns 200"
check_no_error "$body" "list with From+To has no error"

# -------------------------------------------
# Test 12: news.list with combined filters (category + tag)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with category=5&tag=5 (Наука + Космос)${NC}"
response=$(rpc_call "news.list" '{"filter":{"CategoryID":5,"TagID":5}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with category+tag filters returns 200"
check_no_error "$body" "list with category+tag has no error"
check_array_length "$body" "result" 1 "list with category=5&tag=5 returns 1 news"

# -------------------------------------------
# Test 13: news.count with tag filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count with tag=5 (Космос)${NC}"
response=$(rpc_call "news.count" '{"filter":{"TagID":5}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "count with tag filter returns 200"
check_no_error "$body" "count with tag has no error"
check_value "$body" "result" "1" "count with tag=5 returns 1"

# -------------------------------------------
# Test 14: news.count with date range (from)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count with From=2026-03-01T00:00:00Z${NC}"
response=$(rpc_call "news.count" '{"filter":{"From":"2026-03-01T00:00:00Z"}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "count with From filter returns 200"
check_no_error "$body" "count with From has no error"

# -------------------------------------------
# Test 15: news.count with combined filters (category + tag)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count with category=1&tag=1 (Технологии + AI)${NC}"
response=$(rpc_call "news.count" '{"filter":{"CategoryID":1,"TagID":1}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "count with category+tag filters returns 200"
check_no_error "$body" "count with category+tag has no error"
check_value "$body" "result" "1" "count with category=1&tag=1 returns 1"

# -------------------------------------------
# Test 16: news.list with negative category ID
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with category=-1 (invalid category)${NC}"
response=$(rpc_call "news.list" '{"filter":{"CategoryID":-1}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with category=-1 returns 200"
check_error "$body" "list with category=-1 returns error"

# -------------------------------------------
# Test 17: news.list with invalid date range (from > to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with From=2026-12-31T00:00:00Z To=2026-01-01T00:00:00Z (invalid range)${NC}"
response=$(rpc_call "news.list" '{"filter":{"From":"2026-12-31T00:00:00Z","To":"2026-01-01T00:00:00Z"}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with From > To returns 200"
check_error "$body" "list with invalid date range returns error"

# -------------------------------------------
# Test 18: news.count with negative category ID
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count with category=-5 (invalid category)${NC}"
response=$(rpc_call "news.count" '{"filter":{"CategoryID":-5}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "count with category=-5 returns 200"
check_error "$body" "count with category=-5 returns error"

# -------------------------------------------
# Test 19: news.count with invalid date range (from > to)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count with From=2026-06-01T00:00:00Z To=2026-01-01T00:00:00Z (invalid range)${NC}"
response=$(rpc_call "news.count" '{"filter":{"From":"2026-06-01T00:00:00Z","To":"2026-01-01T00:00:00Z"}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "count with From > To returns 200"
check_error "$body" "count with invalid date range returns error"

# -------------------------------------------
# Test 20: news.list - no results for filter
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.list with category=999 (no results)${NC}"
response=$(rpc_call "news.list" '{"filter":{"CategoryID":999}}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "list with non-existent category returns 200"
check_no_error "$body" "list with category=999 has no error"
check_array_length "$body" "result" 0 "list with category=999 returns empty array"

# -------------------------------------------
# Test 21: news.get
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.get with id=1${NC}"
response=$(rpc_call "news.get" '{"id":1}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "get returns 200"
check_no_error "$body" "get has no error"
check_field "$body" "result" "get has result"
check_value "$body" "result.id" "1" "get returns correct ID"
check_value "$body" "result.title" "Прорыв в области искусственного интеллекта" "get returns correct title"
check_field "$body" "result.content" "get includes content"
check_field "$body" "result.preamble" "get includes preamble"
check_field "$body" "result.category" "get includes category"
check_field "$body" "result.tags" "get includes tags"

# -------------------------------------------
# Test 22: news.get - not found
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.get with id=999 (not found)${NC}"
response=$(rpc_call "news.get" '{"id":999}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "get with invalid id returns 200"
check_null_result "$body" "get with invalid id returns null"

# -------------------------------------------
# Test 23: news.get - missing id
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.get (missing id)${NC}"
response=$(rpc_call "news.get" '{}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "get without id returns 200"
check_null_result "$body" "get without id returns null"

# -------------------------------------------
# Test 24: news.get - draft news (not published)
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.get with id=4 (draft news)${NC}"
response=$(rpc_call "news.get" '{"id":4}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "get with draft news id returns 200"
check_null_result "$body" "get with draft news id returns null"

# -------------------------------------------
# Test 25: news.get - news in unpublished category
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.get with id=5 (news in unpublished category)${NC}"
response=$(rpc_call "news.get" '{"id":5}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "get with news in unpublished category returns 200"
check_null_result "$body" "get with news in unpublished category returns null"

# -------------------------------------------
# Test 26: news.get - deleted news
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.get with id=6 (deleted news)${NC}"
response=$(rpc_call "news.get" '{"id":6}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "get with deleted news id returns 200"
check_null_result "$body" "get with deleted news id returns null"

# -------------------------------------------
# Test 27: news.get - news in deleted category
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.get with id=7 (news in deleted category)${NC}"
response=$(rpc_call "news.get" '{"id":7}')
status=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

check_status 200 "$status" "get with news in deleted category returns 200"
check_null_result "$body" "get with news in deleted category returns null"

# -------------------------------------------
# Test 28: news.count vs news.list consistency
# -------------------------------------------
echo ""
echo -e "${YELLOW}Test: news.count equals news.list length${NC}"
count_response=$(rpc_call "news.count" '{}')
count_status=$(echo "$count_response" | tail -n1)
count_body=$(echo "$count_response" | sed '$d')
count_value=$(echo "$count_body" | jq -r ".result")

list_response=$(rpc_call "news.list" '{}')
list_status=$(echo "$list_response" | tail -n1)
list_body=$(echo "$list_response" | sed '$d')
list_length=$(echo "$list_body" | jq ".result | length")

check_status 200 "$count_status" "count returns 200"
check_status 200 "$list_status" "list returns 200"

if [ "$count_value" -eq "$list_length" ]; then
    echo -e "${GREEN}✓ PASS${NC}: Count ($count_value) equals List length ($list_length)"
    ((PASSED++))
else
    echo -e "${RED}✗ FAIL${NC}: Count ($count_value) != List length ($list_length)"
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