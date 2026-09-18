# 本地 Docker 一键启动

Docker Compose 现在包含 MySQL、Redis、4 个 RPC/API 服务。

启动：

docker compose up -d --build

查看：

docker compose ps
docker compose logs -f drama-api

前台：
http://127.0.0.1:8080

后台 API：
http://127.0.0.1:8081

PowerShell 测试：

.scripts	est-api.ps1

停止：

docker compose down

如果需要连同数据一起重置：

docker compose down -v
