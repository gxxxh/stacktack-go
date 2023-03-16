# stacktack-go
A notification monitor for openstack nova




## 配置
#### 获取rabbitmq配置
```shell
cat /etc/rabbitmq/rabbitmq.config
cat /etc/rabbitmq/rabbitmq-env.conf
```

#### 获取binding
```shell
rabbitmqctl list_bindings | grep nova
```

#### 获取nova exchange配置
```shell
rabbitmqctl list_exchanges name type durable auto_delete internal arguments | grep nova
```

#### 获取queue的配置
```shell
rabbitmqctl list_queues name durable auto_delete arguments policy  exclusive exclusive_consumer_pid exclusive_consumer_tag| grep notific
```

#### 查看rabbitmq配置
连接需要修改防火墙暴露端口5672
```shell
lsof -i | grep beam
```
