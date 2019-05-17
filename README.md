# 加密算法说明

### 登陆验证算法
`sha256`

哈希值计算过程： 

1. 真实密码 = 用户密码 X 4
2. sha256( 真实密码 )

### 数据加密算法
`AES256`

数据加密密码为创建用户时从系统读取的32字节随机字符串。

数据加密密码被加密后存储在数据库中。

用于加密数据加密密码的密码形成算法：

1. 用户真实密码 = 用户密码 X 40
2. 真实密码 = sha256( 用户真实密码 )
