---
sidebarTitle: ClashTT 服务运行
sidebarOrder: 3
---

# ClashTT 服务运行

ClashTT 需要在后台运行, 但是目前 Golang 还没有很好的守护进程实现, 因此我们推荐使用第三方工具来创建 ClashTT 的守护进程.

## systemd

使用以下命令将 ClashTT 二进制文件复制到 `/usr/local/bin`, 配置文件复制到 `/etc/clash`:

```shell
cp clash /usr/local/bin
cp config.yaml /etc/clash/
cp Country.mmdb /etc/clash/
```

创建 systemd 配置文件 `/etc/systemd/system/clash.service`:

```ini
[Unit]
Description=ClashTT 守护进程, Go 语言实现的基于规则的代理.
After=network-online.target

[Service]
Type=simple
Restart=always
ExecStart=/usr/local/bin/clash -d /etc/clash

[Install]
WantedBy=multi-user.target
```

之后, 您应该使用以下命令重新加载 systemd:

```shell
systemctl daemon-reload
```

使用以下命令在系统启动时启动 ClashTT:

```shell
systemctl enable clash
```

使用以下命令立即启动 ClashTT:

```shell
systemctl start clash
```

使用以下命令检查 ClashTT 的运行状况和日志:

```shell
systemctl status clash
journalctl -xe
```

本指南贡献者为 [ktechmidas](https://github.com/ktechmidas). ([#754](https://github.com/Dreamacro/clash/issues/754))

## Docker

本项目提供了预构建的 ClashTT 和 ClashTT Premium Docker 镜像. 因此, 在 Linux 上您可以使用 [Docker Compose](https://docs.docker.com/compose/) 部署 ClashTT. 但是, 您应该知道在容器中运行 **ClashTT Premium** 是[不被推荐的](https://github.com/Dreamacro/clash/issues/2249#issuecomment-1203494599)

::: warning
由于 Mac 版 Docker 中缺少[主机网络和 TUN 支持](https://github.com/Dreamacro/clash/issues/770#issuecomment-650951876), 此设置将无法在 macOS 系统上运行.
:::

::: code-group

```yaml [ClashTT]
services:
  clash:
    image: ghcr.io/dreamacro/clash
    restart: always
    volumes:
      - ./config.yaml:/root/.config/clash/config.yaml:ro
      # - ./ui:/ui:ro # 仪表盘 Volume 映射
    ports:
      - "7890:7890"
      - "7891:7891"
      # - "8080:8080" # 外部控制 (RESTful API)
    network_mode: "bridge"
```

```yaml [ClashTT Premium]
services:
  clash:
    image: ghcr.io/dreamacro/clash-premium
    restart: always
    volumes:
      - ./config.yaml:/root/.config/clash/config.yaml:ro
      # - ./ui:/ui:ro # 仪表盘 Volume 映射
    ports:
      - "7890:7890"
      - "7891:7891"
      # - "8080:8080" # 外部控制 (RESTful API)
    cap_add:
      - NET_ADMIN
    devices:
      - /dev/net/tun
    network_mode: "host"
```

:::

保存为 `docker-compose.yaml`, 并将您的 `config.yaml` 放在同一目录下.

::: tip
在继续操作之前, 请参考您的平台关于时间同步的文件 - 如果时间不同步, 某些协议可能无法正常工作.
:::

准备就绪后, 运行以下命令以启动 ClashTT:

```shell
docker-compose up -d
```

您可以使用以下命令查看日志:

```shell
docker-compose logs
```

Stop ClashTT with:

```shell
docker-compose stop
```
