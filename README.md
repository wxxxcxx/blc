# BLC

自动下载 Bilibili 收藏的视频（防止失效）。


# 安装

## 准备

此项目使用 [lux](https://github.com/iawia002/lux) 下载视频，请确保已正确安装 lux 及其依赖。

## 手动安装

1. 安装 [lux](https://github.com/iawia002/lux)，并确保可以正常使用 。
2. 从 [Release](https://github.com/meetcw/blc/releases) 下载编译好的可执行文件 。

## Docker

```
docker pull meetcw/blc
```

### Docker compose
```
version: '3' 
 
services: 
  blc: 
    container_name: blc 
    image: meetcw/blc:latest 
    restart: unless-stopped 
    user: ${UID:?}:${GID:?} 
    volumes: 
      - ${STORAGE:?}:/data
```

Docker 镜像支持的环境变量：
```
ROOT=/data              # 默认下载目录
COOKIE=/data/cookie     # Cookie 文件路径
INTERVAL=3600           # 扫描间隔
```

## 使用

命令参数如下：

```
Usage of blc:
  -cookie string
        Cookie 文件路径 (default "cookie")
  -interval int
        扫描间隔，单位秒 (default 3600)
  -lux string
        Lux 可执行文件路径 (default "lux")
  -root string
        Directory for download
```

Cookie 文件只支持 [Netscape Cookie](https://curl.se/rfc/cookie_spec.html) 格式（ Chrome 或 Firefox 可以使用 [Edit-This-Cookie](https://github.com/ETCExtensions/Edit-This-Cookie) 导出）。