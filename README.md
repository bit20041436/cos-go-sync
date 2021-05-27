# cos-go-sync

本地文件同步到腾讯云cos

## TODO

钉钉接口判断
企业微信通知
上传文件夹处理

## env变量

```shell
# .env.local
Sync_SecretID=SecretID
Sync_SecretKey=SecretKey
Sync_Url=https://存储桶名称.cos.ap-存储桶地域.myqcloud.com
Sync_WebhookUrl=钉钉机器人Webhook地址
```

## 调试

```shell
go run main.go
```

## 编译

```shell
docker build  --no-cache -t cossync .
```

## 测试

```shell
docker run -it  --rm --env-file .env.local -p 443:443 -v $(pwd)/data/:/data/   cossync
```
