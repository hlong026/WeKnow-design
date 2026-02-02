# WeKnora Mac 本地部署指南

本文档详细说明如何将修改后的 WeKnora 项目代码部署到 Mac 电脑上。

## 前置要求

### 1. 安装 Docker Desktop for Mac

1. 访问 [Docker Desktop 官网](https://www.docker.com/products/docker-desktop/)
2. 下载 Mac 版本（注意选择对应芯片版本）：
   - Intel 芯片：选择 "Mac with Intel chip"
   - Apple Silicon (M1/M2/M3/M4)：选择 "Mac with Apple chip"
3. 双击下载的 `.dmg` 文件，将 Docker 拖入 Applications 文件夹
4. 启动 Docker Desktop，等待 Docker 引擎启动完成（菜单栏图标变为稳定状态）

### 2. 配置 Docker Desktop 资源

打开 Docker Desktop → Settings → Resources，建议配置：
- CPUs: 4 核以上
- Memory: 8 GB 以上（推荐 12GB）
- Disk image size: 60 GB 以上

点击 "Apply & Restart" 保存设置。

### 3. 验证 Docker 安装

打开终端（Terminal），运行以下命令验证安装：

```bash
docker --version
docker-compose --version
```

如果显示版本号，说明安装成功。

---

## 部署步骤

### 步骤 1：复制项目文件夹

将整个项目文件夹复制到 Mac 电脑上，可以通过以下方式：
- U盘拷贝
- 网络共享
- AirDrop
- 云盘同步

建议放置路径：`~/Projects/WeKnora` 或 `/Users/用户名/WeKnora`

### 步骤 2：打开终端并进入项目目录

```bash
cd ~/Projects/WeKnora
# 或者你实际放置项目的路径
```

### 步骤 3：配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 使用文本编辑器打开 .env 文件
open -e .env
# 或使用 vim/nano
# vim .env
```

### 步骤 4：编辑 .env 文件

以下是必须配置的关键参数：

```bash
# ========== 必填配置 ==========

# 数据库配置
DB_USER=postgres
DB_PASSWORD=postgres123!@#
DB_NAME=WeKnora

# Redis 配置
REDIS_PASSWORD=redis123!@#

# ========== 可选配置 ==========

# 端口配置（如果默认端口被占用可修改）
APP_PORT=8080
FRONTEND_PORT=80

# 禁用用户认证（内网使用可设为 true）
DISABLE_AUTH=false

# 如果 Mac 在国外，注释掉或删除这行（国内镜像源）
# APK_MIRROR_ARG=mirrors.tencent.com
```

保存并关闭文件。

### 步骤 5：构建并启动服务

```bash
# 构建镜像并启动所有服务（首次运行需要较长时间）
docker-compose up -d --build
```

**说明：**
- `--build` 参数会使用本地代码构建镜像
- `-d` 参数让服务在后台运行
- 首次构建可能需要 10-30 分钟，取决于网络速度

### 步骤 6：查看构建和启动状态

```bash
# 查看所有容器状态
docker-compose ps

# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f docreader
```

等待所有服务状态变为 `healthy` 或 `running`。

### 步骤 7：访问服务

- Web 界面：http://localhost
- 后端 API：http://localhost:8080
- API 健康检查：http://localhost:8080/health

---

## 启动可选组件

根据需要启动额外组件：

```bash
# 启动 MinIO 文件存储
docker-compose --profile minio up -d --build

# 启动 Jaeger 链路追踪
docker-compose --profile jaeger up -d --build

# 启动 Neo4j 知识图谱
docker-compose --profile neo4j up -d --build

# 启动 Qdrant 向量数据库
docker-compose --profile qdrant up -d --build

# 启动全部组件
docker-compose --profile full up -d --build
```

---

## 常用运维命令

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷（清空所有数据）
docker-compose down -v

# 重启所有服务
docker-compose restart

# 重启单个服务
docker-compose restart app

# 重新构建并启动（代码更新后）
docker-compose up -d --build

# 只重新构建特定服务
docker-compose build app
docker-compose up -d app
```

---

## 常见问题排查

### 问题 1：端口被占用

错误信息：`Bind for 0.0.0.0:80 failed: port is already allocated`

解决方案：修改 `.env` 文件中的端口配置

```bash
FRONTEND_PORT=3000
APP_PORT=8081
```

### 问题 2：构建失败 - 网络超时

错误信息：`dial tcp: lookup xxx: no such host` 或下载超时

解决方案：

1. 如果在国外，编辑 `docker/Dockerfile.docreader`，注释掉国内镜像源：
```dockerfile
# RUN sed -i 's@http://deb.debian.org@https://mirrors.tuna.tsinghua.edu.cn@g' /etc/apt/sources.list.d/debian.sources
```

2. 配置 Docker 代理（如有需要）

### 问题 3：Apple Silicon 架构问题

错误信息：`exec format error` 或镜像拉取失败

解决方案：

```bash
# 强制使用 ARM64 架构构建
docker-compose build --build-arg TARGETARCH=arm64
docker-compose up -d
```

### 问题 4：内存不足

错误信息：容器频繁重启或 OOM Killed

解决方案：
1. 增加 Docker Desktop 内存分配（建议 12GB 以上）
2. 减少同时运行的服务

### 问题 5：数据库连接失败

错误信息：`connection refused` 或 `password authentication failed`

解决方案：
1. 确认 `.env` 中的数据库配置正确
2. 等待 PostgreSQL 完全启动：
```bash
docker-compose logs postgres
```

### 问题 6：docreader 服务启动慢

这是正常现象，docreader 需要：
- 下载 OCR 模型
- 安装 Playwright 浏览器

可通过日志查看进度：
```bash
docker-compose logs -f docreader
```

---

## 数据持久化

以下数据会持久化保存，即使容器重启也不会丢失：

| 数据类型 | 存储位置 |
|---------|---------|
| PostgreSQL 数据 | Docker Volume: `postgres-data` |
| 上传的文件 | Docker Volume: `data-files` |
| MinIO 数据 | Docker Volume: `minio_data` |
| Neo4j 数据 | Docker Volume: `neo4j-data` |
| Qdrant 数据 | Docker Volume: `qdrant_data` |

查看数据卷：
```bash
docker volume ls | grep weknora
```

---

## 完全卸载

如需完全卸载并清除所有数据：

```bash
# 停止并删除容器、网络、数据卷
docker-compose down -v

# 删除构建的镜像
docker-compose down --rmi local

# 删除所有相关镜像（可选）
docker images | grep weknora | awk '{print $3}' | xargs docker rmi
```

---

## 快速检查清单

部署完成后，确认以下检查项：

- [ ] Docker Desktop 已启动且资源配置充足
- [ ] `.env` 文件已正确配置
- [ ] `docker-compose ps` 显示所有服务为 `running` 或 `healthy`
- [ ] http://localhost 可以正常访问
- [ ] http://localhost:8080/health 返回正常

如有问题，请查看日志：`docker-compose logs -f`
