## 木犀通行证 v2

![travis-ci](https://travis-ci.org/Muxi-X/muxi_auth_service_v2.svg?branch=master)

### 简介

木犀通行证旨在构建统一的木犀内外门户，本仓库为原Python版本基础之上修改而来，使用Go语言重构。

主要依赖：gin + gorm + viper + lexkong/log

支持Go语言版本： Golang 1.12 及以上

### 构建和运行 Build and run

```
make && ./main
```

### 测试 Testing

```
make test
```

### 包管理

go module

```shell
go mod tidy
```

### APIs

详见 [文档](./api.yaml)

### OAuth

采用 OAuth2.0 标准，使用**授权码模式**进行认证。

#### 登录

授权码模式，客户端要求是前后端分离的应用。

流程：

```
Frontend              Backend                 Auth server

    +                     +                       +
    |                     |                       |
    |                     |                       |
    |                     |                       |
    |           1) login and auth                 |
    |  +--------------------------------------->  |
    |                     |                       |
    |                     |                       |
    |                     |                       |
    |           2) return auth code               |
    |  <------------------+--------------------+  |
    |                     |                       |
    |                     |                       |
    |     3) login        |                       |
    |  +--------------->  |                       |
    |                     |                       |
    |                     | 4) get access token   |
    |                     | +-------------------> |
    |                     |                       |
    |                     |                       |
    |                     | 5)return access token |
    |                     | <-------------------+ |
    |                     |                       |
    | 6)login successfully|                       |
    | <-----------------+ |                       |
    |                     |                       |
    |                     |                       |
    |                     |                       |
    +                     +                       +
```

1. 客户端前端向 Auth 服务器请求 `auth code`，通过 `.../oauth/auth` [API](#登录--获取授权码)
2. 登录成功后，Auth 服务器返回 auth code；
3. 前端向后端请求登录；
4. 后端向 Auth 服务器请求 access token，通过 `.../oauth/token` [API](#get-access-token)
5. 验证通过，Auth 服务器返回 access token；
6. 后端生成 token（客户端应用所用的），返回给前端；
7. 登录成功。


#### 获取用户信息

使用 `access token`，通过 `.../auth/api/user` API 获取用户信息。

#### 更新 access token

客户端通过 [refresh token API](#fresh-access-token) `refresh token`，进行 `access token` 的更新。

#### 客户端注册

使用 [客户端注册 API](#客户端注册与存储) （`.../oauth/store`） 进行客户端注册，获取 `client_id` 和 `client_secret`。

#### OAuth APIs

##### 登录 & 获取授权码

| Path | Method | Header |
| ---  | ---    | ---    |
| /auth/api/oauth | POST | - |

Query Param:
```
    response_type: code （固定字段）
    client_id:
    token_exp: token过期时间，可选
```

Body Data:
```json
{
    "username": "",
    "password": "" // 密码（base64）
}
```

Response:
```json
{
    "code": "",
    "expired": 0, // 过期时间（s）
}
```

##### Get access token

| Path | Method | Header |
| ---  | ---    | ---    |
| /auth/api/oauth/token | POST | - |

Query Param:
```
    grant_type: authorization_code （固定字段）
    response_type: token （固定字段）
    client_id:
```

Body Data (Forms):
```
    client_secret:
    code: 授权码
```

Response Data:
```json
{
    "access_token": "",
    "access_expired": 0, // 过期时间（s）
    "refresh_token": "",
    "refresh_expired": 0 // 过期时间（s）
}
```

##### Refresh access token

| Path | Method | Header |
| ---  | ---    | ---    |
| /auth/api/oauth/token/refresh | POST | - |

Query Param:
```
    grant_type: refresh_token （固定字段）
    client_id:
```

Body Data (Forms):
```
    client_secret:
    refresh_token:
```

Response Data:
```json
{
    "access_token": "",
    "access_expired": 0, // 过期时间（s）
    "refresh_token": "",
    "refresh_expired": 0 // 过期时间（s）
}
```

##### 客户端注册与存储

| Path | Method | Header |
| ---  | ---    | ---    |
| /auth/api/oauth/store | POST | - |

Body Data:
```json
{
    "domain": "" // 域名
}
```

Response Data:
```json
{
    "client_id": "",
    "client_secret": ""
}
```
