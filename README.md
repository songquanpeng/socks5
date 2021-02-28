# Socks5
> Yep, another Socks5 proxy server implementation of Go.

## Usage
```shell script
# Clone & Build
git clone https://github.com/songquanpeng/socks5
cd socks5
go build

# Usage examples:
./socks5.exe
curl -x socks5://localhost:1080 https://baidu.com

./socks5.exe -username stan -password pAssW0rd
curl -x socks5://stan:pAssW0rd@localhost:1080 https://baidu.com
```


## Reference
1. [RFC: SOCKS Protocol Version 5](https://tools.ietf.org/html/rfc1928)
2. [RFC: Username/Password Authentication for SOCKS V5](https://tools.ietf.org/html/rfc1929)
3. [A small Socks5 Proxy Server in Python3](https://github.com/MisterDaneel/pysoxy)