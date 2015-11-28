# phpfpmbeat

The [Beat](https://www.elastic.co/products/beats) for PHP-FPM monitoring.

As Go doesn't currently provide possibility to connect directly to FPM via FastCGI, PHP-FPM must be configured to show its status via some webserver. Sample for Nginx can be found [here](https://easyengine.io/tutorials/php/fpm-status-page/).

### ElasticSearch template

To create template for phpfpmbeat (using cURL):

    curl -XPUT 'http://localhost:9200/_template/phpfpmbeat' -d@etc/phpfpmbeat.template.json

### Build & run

    # Build
    export GO15VENDOREXPERIMENT=1
    GOPATH=<your go path> make

    # Test
    GOPATH=<your go path> make test

    # Run
    ./phpfpmbeat -c etc/phpfpmbeat.yml

### Exported fields

Phpfpmbeat exports all pool information fields provided by PHP-FPM, except FPM start time. To make `slow_requests` counter work, [enable `request_slow_log` feature](https://easyengine.io/tutorials/php/fpm-slow-log/) in you FPM config.

Document sample:

    {
      "@timestamp": "2015-11-28T22:12:04.367Z",
      "beat": {
        "hostname": "Valentins-iMac.local",
        "name": "Valentins-iMac.local"
      },
      "count": 1,
      "phpfpm": {
        "accepted_conn": 218,
        "active_processes": 1,
        "idle_processes": 1,
        "listen_queue": 0,
        "listen_queue_len": 128,
        "max_active_processes": 1,
        "max_children_reached": 0,
        "max_listen_queue": 0,
        "pool": "www",
        "process_manager": "dynamic",
        "slow_requests": 0,
        "start_since": 1176,
        "total_processes": 2
      },
      "type": "phpfpm"
    }

There is no support for per-process info for now.

## Thanks to

- Elastic for their Beat creation tutorial
- [mrkschan](https://github.com/mrkschan) for his [nginxbeat](https://github.com/mrkschan/nginxbeat), which I used as a prototype

## Contribution

Feel free to fork, create merge requests and open issues. I had no experience with Go language previously, so there should be a lot of things to improve.
