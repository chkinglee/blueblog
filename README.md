# blueblog

Blue系列微服务-博客系统

**文件目录**

- api 接口定义
- assets 静态文件资源
- cmd 命令行main入口
- configs 配置文件
- deploy 容器化部署配置文件
- docs 文档
- internal 代码目录
    - app 应用代码
        - controller 控制层
        - service 服务层
        - repository 持久层
    - pkg 应用公共组件
- output 工程编译产出
- pkg 工程公共组件
- vendor 依赖包

**应用**

| app | description |
| --- | --- |
| blueblog-interface | 请求接入，预处理，用户鉴权等 |
| blueblog-service | 业务处理 |
| blueblog-job | 流式任务处理 |
| blueblog-task | 定时任务处理 |
| blueblog-admin | 运营平台 |
