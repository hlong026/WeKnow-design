# Docker 镜像管理脚本使用指南

本目录包含用于管理 WeKnora Docker 镜像的便捷脚本。

## 脚本列表

### 1. rebuild_and_restart.bat (Windows)
完整的镜像重构和重启脚本，包括停止容器、删除旧镜像、重新构建和启动。

### 2. restart.bat (Windows)
快速重启脚本，不重新构建镜像，适用于配置更改后的快速重启。

### 3. rebuild_and_restart.sh (Linux/Mac)
Linux/Mac 版本的完整重构和重启脚本。

## 使用方法

### Windows 系统

#### 重构并重启所有服务（生产环境）
```cmd
scripts\rebuild_and_restart.bat
```

#### 重构并重启所有服务（开发环境）
```cmd
scripts\rebuild_and_restart.bat dev
```

#### 只重构特定服务
```cmd
scripts\rebuild_and_restart.bat prod app
scripts\rebuild_and_restart.bat prod frontend
scripts\rebuild_and_restart.bat prod docreader
```

#### 不使用缓存重构
```cmd
scripts\rebuild_and_restart.bat prod all no-cache
```

#### 使用特定 profile
```cmd
scripts\rebuild_and_restart.bat prod all minio
scripts\rebuild_and_restart.bat prod all full
```

#### 快速重启（不重构）
```cmd
scripts\restart.bat
scripts\restart.bat dev
scripts\restart.bat prod app
```

### Linux/Mac 系统

首先给脚本添加执行权限：
```bash
chmod +x scripts/rebuild_and_restart.sh
```

#### 重构并重启所有服务（生产环境）
```bash
./scripts/rebuild_and_restart.sh
```

#### 重构并重启所有服务（开发环境）
```bash
./scripts/rebuild_and_restart.sh --dev
```

#### 只重构特定服务
```bash
./scripts/rebuild_and_restart.sh --service app
./scripts/rebuild_and_restart.sh --service frontend
./scripts/rebuild_and_restart.sh --service docreader
```

#### 不使用缓存重构
```bash
./scripts/rebuild_and_restart.sh --no-cache
```

#### 使用特定 profile
```bash
./scripts/rebuild_and_restart.sh --profile minio
./scripts/rebuild_and_restart.sh --profile full
```

## 常用场景

### 场景 1: 代码更新后重新部署
```cmd
REM Windows
scripts\rebuild_and_restart.bat prod all no-cache

# Linux/Mac
./scripts/rebuild_and_restart.sh --no-cache
```

### 场景 2: 只更新前端
```cmd
REM Windows
scripts\rebuild_and_restart.bat prod frontend

# Linux/Mac
./scripts/rebuild_and_restart.sh --service frontend
```

### 场景 3: 开发环境快速重启
```cmd
REM Windows
scripts\restart.bat dev

# Linux/Mac
docker compose -f docker-compose.dev.yml restart
```

### 场景 4: 配置文件更改后重启
```cmd
REM Windows
scripts\restart.bat prod app

# Linux/Mac
docker compose -f docker-compose.yml restart app
```

## 手动 Docker Compose 命令

如果需要更精细的控制，可以直接使用 docker compose 命令：

### 查看服务状态
```bash
docker compose -f docker-compose.yml ps
```

### 查看日志
```bash
docker compose -f docker-compose.yml logs -f
docker compose -f docker-compose.yml logs -f app
```

### 停止所有服务
```bash
docker compose -f docker-compose.yml down
```

### 停止并删除数据卷（危险操作！）
```bash
docker compose -f docker-compose.yml down -v
```

### 重新构建特定服务
```bash
docker compose -f docker-compose.yml build --no-cache app
docker compose -f docker-compose.yml up -d app
```

### 使用 profile 启动
```bash
docker compose -f docker-compose.yml --profile full up -d
docker compose -f docker-compose.yml --profile minio up -d
```

## 注意事项

1. **数据安全**: 重构镜像不会删除数据卷，你的数据是安全的
2. **缓存使用**: 默认使用构建缓存加速，如需完全重新构建请使用 `no-cache` 选项
3. **环境变量**: 确保 `.env` 文件配置正确
4. **端口冲突**: 确保所需端口未被其他程序占用
5. **磁盘空间**: 定期清理未使用的镜像和容器以释放空间

## 清理 Docker 资源

### 清理未使用的镜像
```bash
docker image prune -a
```

### 清理未使用的容器
```bash
docker container prune
```

### 清理未使用的数据卷
```bash
docker volume prune
```

### 清理所有未使用的资源
```bash
docker system prune -a --volumes
```

## 故障排查

### 问题 1: 端口被占用
```bash
# 查看端口占用
netstat -ano | findstr :8080
# 或使用 docker ps 查看是否有容器在运行
docker ps
```

### 问题 2: 构建失败
```bash
# 查看详细构建日志
docker compose -f docker-compose.yml build --no-cache --progress=plain
```

### 问题 3: 服务无法启动
```bash
# 查看服务日志
docker compose -f docker-compose.yml logs app
# 查看容器状态
docker compose -f docker-compose.yml ps
```

### 问题 4: 数据库连接失败
```bash
# 检查数据库是否健康
docker compose -f docker-compose.yml ps postgres
# 查看数据库日志
docker compose -f docker-compose.yml logs postgres
```
