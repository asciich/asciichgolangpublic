# signal-cli-rest-api utils

Utils to interact with the [signal-cli-rest-api](https://github.com/bbernhard/signal-cli-rest-api).

## RunReceiveCacheServer

The `RunReceiveCacheServer` function starts a web server that exposes `/messages` endpoint to get cached messages in JSON format. It automatically receives new messages at a given interval and drops the latest message when the cache is full.
