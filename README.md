# Short URL Service
A short url service

## Prerequisite
Install docker.

Follow the instruction on the official page: https://docs.docker.com/desktop/setup/install/mac-install/

Test
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