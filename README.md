# graproxy
simple implementation for grafana auth.proxy with openid


## grafana config example

    [auth.proxy]
    enabled = true
    header_name = X-WEBAUTH-USER
    header_property = email
    auto_sign_up = true


## usage

    Usage of ./graproxy:
      -endpoint string
        openid endpoint
      -grafana string
        grafana host and port (default "localhost")
      -l string
        proxy server listen address (default ":8080")
      -n string
        proxy domain name openid will return to (default "localhost")

## living example

    ./graproxy -endpoint "https://login.provider.com/openid" -grafana "127.0.0.1:3000" -l ":8080" -n monitor.example.com


## build

    cd src && go build -o ~/graproxy
