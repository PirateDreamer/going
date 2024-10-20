基于k8s的grpc的封装

proto文件统一管理，所有proto文件放到同一个git仓库，其他git项目引入这个仓库作为子仓库，仓库名称建议是proto

使用buf管理proto生成代码,插件和proto文件依赖

服务支持http和grpc，其中http使用grpc-gateway实现

错误处理，结果封装、链路追踪、日志 (拦截器)


使用k8s负载均衡、熔断、限流能力

分布式锁、分布式ID、分布式事务
