<!-- This is the index page, linked by the dummy sidebar item at Introduction/_dummy-index.md -->
# What is ClashT?

Welcome to the official knowledge base of the alt-Clash core project ("ClashT").

ClashT is a cross-platform rule-based proxy utility that runs on the network and application layer, supporting various proxy and anti-censorship protocols out-of-the-box.

It has been adopted widely by the Internet users in some countries and regions where the Internet is heavily censored or blocked. Either way, ClashT can be used by anyone who wants to improve their Internet experience.

There are currently one editions of Clash:

- ~~[Clash](https://github.com/Dreamacro/clash): the open-source version released at [github.com/Dreamacro/clash](https://github.com/Dreamacro/clash)~~
- ~~[Clash Premium](https://github.com/Dreamacro/clash/releases/tag/premium): proprietary core with [TUN support and more](/premium/introduction) (free of charge)~~
- [ClashT](https://github.com/DryPeng/clashT): the open-source version released at [github.com/DryPeng/clashT](https://github.com/DryPeng/clashT)

While this wiki covers both, however, the use of Clash could be challenging for the average users. Those might want to consider using a GUI client instead, and we do have some recommendations:

- ~~[Clash for Windows](https://github.com/Fndroid/clash_for_windows_pkg/releases) (Windows and macOS)~~
- ~~[Clash for Android](https://github.com/Kr328/ClashForAndroid)~~
- ~~[ClashX](https://github.com/yichengchen/clashX) or [ClashX Pro](https://install.appcenter.ms/users/clashx/apps/clashx-pro/distribution_groups/public) (macOS)~~

## Feature Overview

- Inbound: HTTP, HTTPS, SOCKS5 server, TUN device<sup>*1</sup>.
- Outbound: Shadowsocks(R), VMess, Trojan, Snell, SOCKS5, HTTP(S), Wireguard<sup>*1</sup>.
- Rule-based Routing: dynamic scripting, domain, IP addresses, process name and more<sup>*1</sup>.
- Fake-IP DNS: minimises impact on DNS pollution and improves network performance.
- Transparent Proxy: Redirect TCP and TProxy TCP/UDP with automatic route table/rule management<sup>*1</sup>.
- Proxy Groups: automatic fallback, load balancing or latency testing.
- Remote Providers: load remote proxy lists dynamically.
- RESTful API: update configuration in-place via a comprehensive API.

<small>*1: Only available in the Premium edition (Awaiting maintenance).</small>


## License

Clash is released under the [GPL-3.0](https://github.com/Dreamacro/clash/blob/master/LICENSE) open-source license. Prior to [v0.16.0](https://github.com/Dreamacro/clash/releases/tag/v0.16.0) or commit [e5284c](https://github.com/Dreamacro/clash/commit/e5284cf647717a8087a185d88d15a01096274bc2), it was licensed under the MIT license.
