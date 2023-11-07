---
sidebarTitle: ClashT as a Service
sidebarOrder: 3
---

# ClashT as a Service

While ClashT is meant to be run in the background, there's currently no elegant way to implement daemons with Golang, hence we recommend you to daemonize ClashT with third-party tools.

## systemd

Copy ClashT binary to `/usr/local/bin` and configuration files to `/etc/clash`:

```shell
cp clash /usr/local/bin
cp config.yaml /etc/clash/
cp Country.mmdb /etc/clash/
```

Create the systemd configuration file at `/etc/systemd/system/clash.service`:

```ini
[Unit]
Description=ClashT daemon, A rule-based proxy in Go.
After=network-online.target

[Service]
Type=simple
Restart=always
ExecStart=/usr/local/bin/clash -d /etc/clash

[Install]
WantedBy=multi-user.target
```

After that you're supposed to reload systemd:

```shell
systemctl daemon-reload
```

Launch clashd on system startup with:

```shell
systemctl enable clash
```

Launch clashd immediately with:

```shell
systemctl start clash
```

Check the health and logs of ClashT with:

```shell
systemctl status clash
journalctl -xe
```

Credits to [ktechmidas](https://github.com/ktechmidas) for this guide. ([#754](https://github.com/Dreamacro/clash/issues/754))

## Docker

We provide pre-built images of ClashT and ClashT Premium. Therefore you can deploy ClashT with [Docker Compose](https://docs.docker.com/compose/) if you're on Linux. However, you should be advised that it's [not recommended](https://github.com/Dreamacro/clash/issues/2249#issuecomment-1203494599) to run **ClashT Premium** in a container.

::: warning
This setup will not work on macOS systems due to the lack of [host networking and TUN support](https://github.com/Dreamacro/clash/issues/770#issuecomment-650951876) in Docker for Mac.
:::


::: code-group

```yaml [ClashT]
services:
  clash:
    image: ghcr.io/dreamacro/clash
    restart: always
    volumes:
      - ./config.yaml:/root/.config/clash/config.yaml:ro
      # - ./ui:/ui:ro # dashboard volume
    ports:
      - "7890:7890"
      - "7891:7891"
      # - "8080:8080" # The External Controller (RESTful API)
    network_mode: "bridge"
```

```yaml [ClashT Premium]
services:
  clash:
    image: ghcr.io/dreamacro/clash-premium
    restart: always
    volumes:
      - ./config.yaml:/root/.config/clash/config.yaml:ro
      # - ./ui:/ui:ro # dashboard volume
    ports:
      - "7890:7890"
      - "7891:7891"
      # - "8080:8080" # The External Controller (RESTful API)
    cap_add:
      - NET_ADMIN
    devices:
      - /dev/net/tun
    network_mode: "host"
```

:::

Save as `docker-compose.yaml` and place your `config.yaml` in the same directory.

::: tip
Before proceeding, refer to your platform documentations about time synchronisation - things will break if time is not in sync.
:::

When you're ready, run the following commands to bring up ClashT:

```shell
docker-compose up -d
```

You can view the logs with:

```shell
docker-compose logs
```

Stop ClashT with:

```shell
docker-compose stop
```

## FreeBSD rc

install clash with `ports(7)` or `pkg(8)`

copy the required files to `/usr/local/etc/clash`

```shell
cp config.yaml /usr/local/etc/clash/
cp Country.mmdb /usr/local/etc/clash/
```

Create the rc configuration file at `/usr/local/etc/rc.d/clash`:

```shell
#!/bin/sh

# PROVIDE: clash
# REQUIRE: NETWORKING DAEMON
# BEFORE: LOGIN
# KEYWORD: shutdown

. /etc/rc.subr

name=clash
rcvar=clash_enable

: ${clash_enable="NO"}
: ${clash_config_dir="/usr/local/etc/clash"}

required_dirs="${clash_config_dir}"
required_files="${clash_config_dir}/config.yaml ${clash_config_dir}/Country.mmdb"

command="/usr/sbin/daemon"
procname="/usr/local/bin/${name}"
pidfile="/var/run/${name}.pid"
start_precmd="${name}_prestart"

clash_prestart()
{
	rc_flags="-T ${name} -p ${pidfile} ${procname} -d ${clash_config_dir} ${rc_flags}"
}

load_rc_config $name
run_rc_command "$1"
```

make the script executable:

```shell
chmod +x /usr/local/etc/rc.d/clash
```

Launch clashd on system startup with:

```shell
service clash enable
```

Launch clashd immediately with:

```shell
service clash onestart
```

Check the status of ClashT with:

```shell
service clash status
```

You can check log in file `/var/log/daemon.log`

::: tip
If you want to change the default config directory add the following lines to /etc/rc.conf :
```shell
clash_enable (bool):        Set it to YES to run clash on startup.
                            Default: NO
clash_config_dir (string):   clash config directory.
                            Default: /usr/loca/etc/clash
```
:::
