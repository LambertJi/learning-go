# Learning-Go

> 学习 Go 语言 + AI 编程助手协作开发

## 简介

本项目是 Go 语言学习的练习代码仓库，用于记录学习过程中的各种语法练习、库封装和验证测试。

**项目目标：**
- 学习 Go 语言基础和常用库的使用
- 探索人机协作开发模式，学习使用 AI 编程助手提升开发效率
- 积累可复用的代码片段和最佳实践

## 开发方式

本项目由 **@LambertJi** 与 **Claude Code + GLM4.7** 共同完成。

## 环境要求

- Go 1.25+

## 项目结构

```
learning-go/
├── internals/           # 内部库封装
│   ├── httpx/          # HTTP 客户端封装
│   └── redisx/         # Redis 客户端封装
├── validation/         # 功能验证与测试
├── output/            # 输出目录
├── main.go            # 主程序
├── basic_syntax.go    # 基础语法练习
├── go.mod             # Go 模块定义
└── go.sum             # 依赖版本锁定
```

## 安装依赖

```bash
go mod download
```

## 使用说明

<!-- 在这里填写使用说明 -->

## 依赖库

- [redis/go-redis/v9](https://github.com/redis/go-redis) - Redis 客户端

## 许可证

MIT
