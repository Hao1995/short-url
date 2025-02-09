# Short URL Service
A short url service

## Prerequisite
### Docker
Follow the instruction on the official page: https://docs.docker.com/desktop/setup/install/mac-install/

Check successful installation
```
docker version
docker-compose version
```

## Quick Start
Start the service
```
make up
```

Shutdown the service
```
make down
```

Test
```
# Upload URL
curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls -d '{
"url": "https://www.dcard.tw/",
"expireAt": "2021-02-08T09:20:41Z"
}'
# Response
{
"id": "<url_id>",
"shortUrl": "http://localhost/<url_id>"
}

# Redirect URL API
# Use the `url_id` from the previous response 
curl -L -X GET http://localhost/<url_id>

```

# Thought Process
> 寫上你的思路與決定使用 db、3rd party lib 的原因

額外問題釐清
1. 想請問系統有預計可以容納多少 short url 嗎？ >> 至少需要能夠存 Millions 以上
2. 有無限制可用的字元類型（例如僅限 english letter、不允許數字）>> 這個請列入您設計的考量
3. 有無限制 url id 的長度？(不包含 domain 的情況下) >> 這個請列入您設計的考量
4. 有預期需要承受多少 RPS 嗎？有的話兩支 API 分別需求是多少？ >> 題目重點是會受到 malicious attack
5. 第 3 點有提到可用的 storage 選項，是僅限於這些選項，還是可以用其他的？ >> 可以用其它的

## DB Choose
由於資料量不多，且目的為預防 `malicious attack`，所以簡易選了 MySQL 作為 DB 的情況下，額外加上 cache 來解決 non-existent shorten URL 的問題。
因為 cache 只能處理 read 需求，所以如果 write 需求很高的話，再來考慮採用 Cassandra 這種分散式資料庫來提升性能即可。

## Address Non-existent Shorten URL
由於不希望惡意用戶一直嘗試不存在的短網址時，因為 key 在 cache 找不到所以一直往 DB request 造成資料庫性能問題，所以我同時把 `record not found` 的結果也存在 cache，所以就算用戶一直嘗試，也不會對服務造成負擔。(當然進一步還有 firewall、ip detect 等預防惡意攻擊的方式可以做。)

## Hash function 的採用
CRC32 為 32 bits，最大可容納資料為 4,294,967,296 (4,294M)，可符合 millions 的需求。
另外，由於一般短網址服務，不會限定同一個 url 不能再次請行短網址產生，所以我而外使用 random 字串來避免 hash collision。

## Cache Library
cache lib 採用的是 [viney-shih/go-cache](https://github.com/viney-shih/go-cache)，GET 請求發出的時候，會先到 local cache 尋找是否有資料，沒有的話再到 shared cache 尋找，而且背後使用 singleflight，同一時間若有多個 requests，只會有一個 request 真的去後面拿資料，相同的請求會等該目前請求結束後一起分享資料，避免 cache miss 的時候，大量 requests 同時往 DB 請求，造成性能瓶頸。

## DB Related Libraries
- Gorm
    - Golang 的大宗 orm 套件，避免 SQL injection 問題。
- Goose
    - DB migration 套件，幫助開發時可以更方便修改 DB schema。

## HTTP Web Framework
- Gin
    - 主流 HTTP web framework，用來處理 request 驗證、router 設定等工作。
- fvbock/endless
    - Gin 官方推薦，用來做 graceful shutdown 的 lib。

## Test
- testify
    - 知名測試套件，方便撰寫測試、驗證 actual、expected 內容。
- dockertest
    - 藉由在測試的時候啟動 mysql、redis 等 container 來達到 integration test，確保服務運作符合預期。

## Configuration
- caarlos0/env
    - 使 env 參數可以直接映射到 object，方便使用各種參數。