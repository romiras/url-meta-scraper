# url-meta-scraper

url-meta-scraper is a worker which primary task is downloading URLs that consumed from MQ.

On success it will push the content of URL to queue `urls`.
If downloading of URL failed, it will push to separate queue `failed-urls` for further retry by another worker.

## Message Queue adapters

* AMQP

## Status

WIP

- [x] ~~Producer for downloaded URLs~~
- [ ] Producer for failed URLs

## License

See file LICENSE.
