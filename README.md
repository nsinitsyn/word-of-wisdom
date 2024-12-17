# “Word of Wisdom” tcp server

## How to run
Build and run server and client by commands:
```
cd deploy/dev
docker compose up -d
```

Check client output:
```
docker logs wow-client
>>2024/12/17 14:58:45 Life isn't about getting and having, it's about giving and being. –Kevin Kruse
```

## Explanation of choice the hashcash POW algorithm
hashcash - most popular and simple POW algorithm using, for example, in Bitcoin