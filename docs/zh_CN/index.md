<!-- 这是 index 页面, 由位于 Introduction/_dummy-index.md 的虚拟侧边栏文件链接 -->
# 什么是 ClashT?

欢迎访问 ClashT 内核项目的官方说明文档.

ClashT 是一个跨平台的基于规则的代理工具, 在网络和应用层运行, 支持各种代理和反审查协议的开箱即用.

在一些互联网受到严格审查或封锁的国家和地区, 它已被互联网用户广泛采用. 无论如何, 任何想要改善其 Internet 体验的人都可以使用 ClashT.

目前, ClashT 包含一个版本:

- ~~[ClashT](https://github.com/Dreamacro/clash): 发布于[github.com/Dreamacro/clash](https://github.com/Dreamacro/clash)的开源版本~~
- ~~[ClashT Premium 版本](https://github.com/Dreamacro/clash/releases/tag/premium): 具有[TUN 和更多支持](/zh_CN/premium/introduction) 的专有内核 (免费)~~
- [ClashT](https://github.com/DryPeng/clashT): 发布于[github.com/DryPeng/clashT](https://github.com/DryPeng/clashT)的开源版本

虽然这个 Wiki 涵盖了上述版本的内容, 然而对于普通用户来说, ClashT 的使用可能仍是一种挑战. 而对于考虑使用 GUI 客户端的用户, 我们确实有一些建议:

- ~~[ClashT for Windows](https://github.com/Fndroid/clash_for_windows_pkg/releases) (Windows 和 macOS)~~
- ~~[ClashT for Android](https://github.com/Kr328/ClashTForAndroid)~~
- ~~[ClashTX](https://github.com/yichengchen/clashX) 或 [ClashTX Pro](https://install.appcenter.ms/users/clashx/apps/clashx-pro/distribution_groups/public) (macOS)~~

## 特点概述

- 入站连接支持: HTTP, HTTPS, SOCKS5 服务端, TUN 设备<sup>*1</sup>
- 出站连接支持: Shadowsocks(R), VMess, Trojan, Snell, SOCKS5, HTTP(S), Wireguard<sup>*1</sup>
- 基于规则的路由: 动态脚本、域名、IP地址、进程名称和更多<sup>*1</sup>
- Fake-IP DNS: 尽量减少 DNS 污染的影响, 提高网络性能
- 透明代理: 使用自动路由表/规则管理 Redirect TCP 和 TProxy TCP/UDP<sup>*1</sup>
- Proxy Groups 策略组: 自动化的可用性测试 (fallback)、负载均衡 (load balance) 或 延迟测试 (url-test)
- 远程 Providers: 动态加载远程代理列表
- RESTful API: 通过一个全面的 API 就地更新配置

<small>*1 只在免费的 Premium 版本中提供。</small>

## License

ClashT 是根据 [GPL-3.0](https://github.com/Dreamacro/clash/blob/master/LICENSE) 开源许可证发布的. 在 [v0.16.0](https://github.com/Dreamacro/clash/releases/tag/v0.16.0) 或 [e5284c](https://github.com/Dreamacro/clash/commit/e5284cf647717a8087a185d88d15a01096274bc2) 提交之前, 其基于 MIT 许可证授权.
