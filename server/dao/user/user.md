# server/dao/user/user.go 技术文档

## 概述

`server/dao/user/user.go` 文件是 AgentGo 项目的用户数据访问层（DAO - Data Access Object），提供了与用户相关的数据库操作功能。该文件实现了用户存在性检查和用户注册功能，通过 MySQL 数据库进行数据持久化。

## 包信息

- **包名**: `user`
- **功能**: 提供用户相关的数据访问功能
- **导入路径**: `server/dao/user`

## 常量定义

### CodeMsg
- **值**: `"AgentGo验证码如下（验证码仅限于2分钟内有效）: "`
- **用途**: 用于生成验证码消息的前缀文本

### UserNameMsg
- **值**: `"AgentGo的账号如下：请保留好，后续可以用账号进行登录 "`
- **用途**: 用于生成账号信息消息的前缀文本

## 全局变量

### ctx
- **类型**: `context.Context`
- **值**: `context.Background()`
- **用途**: 为数据库操作提供上下文支持

## 核心功能

### IsExistUser 函数
- **功能**: 判断用户是否存在（系统仅支持通过账号进行登录）
- **参数**:
    - `username` (string): 用户名
- **返回值**:
    - `bool`: 用户存在返回 true，否则返回 false
    - `*models.User`: 用户存在时返回用户信息，否则返回 nil
- **实现细节**:
    - 调用 `mysql.GetUserByUsername(username)` 从数据库查询用户
    - 检查错误是否为 `gorm.ErrRecordNotFound` 或用户对象是否为 nil
    - 根据查询结果返回相应的布尔值和用户对象

### Register 函数
- **功能**: 用户注册
- **参数**:
    - `username` (string): 用户名
    - `email` (string): 邮箱地址
    - `password` (string): 密码
- **返回值**:
    - `*models.User`: 注册成功时返回用户信息，失败时返回 nil
    - `bool`: 注册成功返回 true，失败返回 false
- **实现细节**:
    - 创建 [User](file:///F:/AgentGo/server/models/user.go#L9-L18) 对象，包含邮箱、用户名、姓名和 MD5 加密后的密码
    - 调用 `mysql.InsertUser()` 将用户信息插入数据库
    - 处理插入错误，成功则返回用户信息和 true，失败则返回 nil 和 false

## 数据模型依赖

该文件依赖于 `server/models/user.go` 中定义的 [User](file:///F:/AgentGo/server/models/user.go#L9-L18) 模型，该模型包含以下字段：
- [ID](file://F:\AgentGo\server\models\user.go#L9-L9): 主键，自增
- [Name](file://F:\AgentGo\server\models\user.go#L10-L10): 用户昵称
- [Email](file://F:\AgentGo\server\models\user.go#L11-L11): 邮箱地址，支持索引
- [Username](file://F:\AgentGo\server\models\user.go#L12-L12): 唯一用户名，用于登录
- [Password](file://F:\AgentGo\server\models\user.go#L13-L13): MD5 加密后的密码
- [CreatedAt](file://F:\AgentGo\server\models\user.go#L14-L14)/`UpdatedAt`: 自动时间戳
- [DeletedAt](file://F:\AgentGo\server\models\user.go#L16-L16): 支持软删除

## 外部依赖

- **GORM**: 用于 ORM 操作
- **MySQL**: 数据库存储
- **Utils 包**: 提供 `MD5` 加密功能

## 安全考虑

- 用户密码在存储前使用 MD5 加密
- 通过唯一索引确保用户名不重复
- 使用 GORM 的软删除功能保护数据

## 架构角色

作为数据访问层，该文件在系统架构中扮演以下角色：
- **业务逻辑与数据存储的桥梁**: 封装数据库操作细节，为上层业务逻辑提供简单的接口
- **数据验证**: 在插入前对用户输入进行基本验证
- **错误处理**: 处理数据库操作中的各种错误情况

