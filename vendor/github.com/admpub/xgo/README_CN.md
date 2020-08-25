# Docker 操作备忘录

## 进入容器的方式

```bash
docker run -it --entrypoint /bin/bash  admpub/xgo:latest
```

## 添加国内源

点击路径：`Preferences` -> `Docker Engine`，输入如下配置后点击按钮`Apply & Restart`

```js
{
  "experimental": false,
  "debug": true,
  "registry-mirrors": [
    "https://fkswhxnd.mirror.aliyuncs.com",
    "http://hub-mirror.c.163.com",
    "https://docker.mirrors.ustc.edu.cn",
    "https://docker.mirrors.ustc.edu.cn"
  ]
}
```
