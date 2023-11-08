<h1 align="center">
  <img src="https://github.com/DryPeng/clashT/raw/master/docs/logo.png" alt="Clash" width="200">
  <br>ClashT<br>
</h1>

<h4 align="center">A rule-based tunnel in Go.</h4>

<p align="center">
  <a href="https://github.com/DryPeng/clashT/actions">
    <img src="https://img.shields.io/github/actions/workflow/status/DryPeng/clashT/release.yml?branch=master&style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/DryPeng/clashT">
    <img src="https://goreportcard.com/badge/github.com/DryPeng/clashT?style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/DryPeng/clashT?style=flat-square">
  <a href="https://github.com/DryPeng/clashT/releases">
    <img src="https://img.shields.io/github/release/DryPeng/clashT/all.svg?style=flat-square">
  </a>
  <!-- <a href="https://github.com/DryPeng/clashT/releases/tag/premium">
    <img src="https://img.shields.io/badge/release-Premium-00b4f0?style=flat-square">
  </a> -->
</p>

> [!WARNING]  
> This is not the [original *Clash*](https://github.com/Dreamacro/clash). This repositories is currently trying maintained. For convenience, *ClashT* is a new repositories.

> [!NOTE]  
> Do you want to contribute? Just [open a PR](https://github.com/DryPeng/clashT/pulls).

> [!IMPORTANT]  
> We are collecting [original *Clash*](https://github.com/Dreamacro/clash) issues and suggestions, [click here to submit](https://github.com/DryPeng/clashT/issues).


## Features

This is a general overview of the features that comes with *Clash*.  

- Inbound: HTTP, HTTPS, SOCKS5 server, TUN device.
- Outbound: Shadowsocks(R), VMess, Trojan, Snell, SOCKS5, HTTP(S), Wireguard.
- Rule-based Routing: dynamic scripting, domain, IP addresses, process name and more.
- Fake-IP DNS: minimises impact on DNS pollution and improves network performance.
- Transparent Proxy: Redirect TCP and TProxy TCP/UDP with automatic route table/rule management.
- Proxy Groups: automatic fallback, load balancing or latency testing.
- Remote Providers: load remote proxy lists dynamically.
- RESTful API: update configuration in-place via a comprehensive API.

*Some of the features may only be available in the [Premium core](https://github.com/DryPeng/clashT/blob/master/docs/premium/introduction.md)* (Awaiting maintenance).


## Documentation
> [!Note]
> We have temporarily ensured the readability of the document. Because the original repo ([Clash](https://github.com/Dreamacro/clash/)) cannot be viewed (404), some "issue" cannot be viewed.

> [!note]
> If you find any issues with the documentation, please feel free to [let us know](https://github.com/DryPeng/clashT/issues/).

You can find the latest documentation at https://clashT.drypeng.io/.


## Credits

- [Dreamacro/clash](https://github.com/Dreamacro/clash) (original *Clash*)
- [riobard/go-shadowsocks2](https://github.com/riobard/go-shadowsocks2)
- [v2ray/v2ray-core](https://github.com/v2ray/v2ray-core)
- [WireGuard/wireguard-go](https://github.com/WireGuard/wireguard-go)


## 	:heart_on_fire: Special thanks 	:heart_on_fire:

#### Outside Developer

[@Baidu0Big](https://github.com/Baidu0Big): Chinese (zh-CN and zh-TW) translation and code review


## License

This software is released under the GPL-3.0 license.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FDryPeng%2FclashT.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2FDryPeng%2FclashT?ref=badge_large&issueType=license)
