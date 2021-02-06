

Centos7下单节点部署etcd服务

一台Centos7的服务器，假设IP：172.16.1.11

登陆到服务器，切换到root用户



下载安装包
=============
下载地址：https://github.com/etcd-io/etcd/releases
```sh
cd /usr/local/src/
wget https://github.com/etcd-io/etcd/releases/download/v3.3.13/etcd-v3.3.13-linux-amd64.tar.gz
```


解压文件
=============
```sh
tar -zxvf etcd-v3.3.13-linux-amd64.tar.gz -C /usr/local

cd /usr/local/etcd-v3.3.13-linux-amd64
```
目录下有两个可执行文件etcd 和 etcdctl，把它们复制到/usr/bin/下


```sh
cp etcd /usr/bin/
cp etcdctl /usr/bin/
```



配置
=============
 
1、配置etcd.service
-------------
在/usr/lib/systemd/system/目录下新建etcd.service文件， 执行命令：

> touch /usr/lib/systemd/system/etcd.service

配置内容如下：
```sh
[Unit]
Description=Etcd Server
After=network.target

[Service]
Type=simple
WorkingDirectory=/var/lib/etcd/
EnvironmentFile=-/etc/etcd/etcd.conf

ExecStart=/usr/bin/etcd

[Install]
WantedBy=multi-user.target
````


2、新建etcd工作目录
-------------
在/var/lib/目录下新建etcd的工作目录etcd，执行命令：

> mkdir /var/lib/etcd


3、配置etcd.conf
-------------
新建/etc/etcd/etcd.conf文件，执行以下命令：

> mkdir /etc/etcd  
> touch /etc/etcd/etcd.conf  

配置内容如下：
```sh
#[member]
ETCD_NAME=default
ETCD_DATA_DIR="/var/lib/etcd/default.etcd"
ETCD_LISTEN_CLIENT_URLS="http://172.16.1.11:2379"

ETCD_ADVERTISE_CLIENT_URLS="http://172.16.1.11:2379"
```


启动并验证
================
配置完成后，执行以下命令，启动etcd服务。
```sh
#重载所有修改过的配置文件；
systemctl daemon-reload 

#将etcd服务加入开机启动列表
systemctl enable etcd.service 

#启动etcd服务
systemctl start etcd.service 
```


启动后执行以下命令验证：
```sh
etcdctl cluster-health

#输出：

member 8e9e05c52164694d is healthy: got healthy result from http://172.16.1.11:2379
cluster is healthy
```






etcd集群
==================
etcd 作为一个高可用键值存储系统，天生就是为集群化而设计的。由于 Raft 算法在做决策时需要多数节点的投票，所以 etcd 一般部署集群推荐奇数个节点，推荐的数量为 3、5 或者 7 个节点构成一个集群。

搭建一个3节点集群示例：
在每个etcd节点指定集群成员，为了区分不同的集群最好同时配置一个独一无二的token。

下面是提前定义好的集群信息，其中n1、n2和n3表示3个不同的etcd节点。

```sh
TOKEN=token-01
CLUSTER_STATE=new
CLUSTER=n1=http://10.240.0.17:2380,n2=http://10.240.0.18:2380,n3=http://10.240.0.19:2380
```

在n1这台机器上执行以下命令来启动etcd：
```sh
etcd --data-dir=data.etcd --name n1 \
    --initial-advertise-peer-urls http://10.240.0.17:2380 --listen-peer-urls http://10.240.0.17:2380 \
    --advertise-client-urls http://10.240.0.17:2379 --listen-client-urls http://10.240.0.17:2379 \
    --initial-cluster ${CLUSTER} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

在n2这台机器上执行以下命令启动etcd：
```sh
etcd --data-dir=data.etcd --name n2 \
    --initial-advertise-peer-urls http://10.240.0.18:2380 --listen-peer-urls http://10.240.0.18:2380 \
    --advertise-client-urls http://10.240.0.18:2379 --listen-client-urls http://10.240.0.18:2379 \
    --initial-cluster ${CLUSTER} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

在n3这台机器上执行以下命令启动etcd：
```sh
etcd --data-dir=data.etcd --name n3 \
    --initial-advertise-peer-urls http://10.240.0.19:2380 --listen-peer-urls http://10.240.0.19:2380 \
    --advertise-client-urls http://10.240.0.19:2379 --listen-client-urls http://10.240.0.19:2379 \
    --initial-cluster ${CLUSTER} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

etcd 官网提供了一个可以公网访问的 etcd 存储地址。你可以通过如下命令得到 etcd 服务的目录，并把它作为-discovery参数使用。

```sh
curl https://discovery.etcd.io/new?size=3
https://discovery.etcd.io/a81b5818e67a6ea83e9d4daea5ecbc92
 
# grab this token
TOKEN=token-01
CLUSTER_STATE=new
DISCOVERY=https://discovery.etcd.io/a81b5818e67a6ea83e9d4daea5ecbc92
 
 
etcd --data-dir=data.etcd --name n1 \
    --initial-advertise-peer-urls http://10.240.0.17:2380 --listen-peer-urls http://10.240.0.17:2380 \
    --advertise-client-urls http://10.240.0.17:2379 --listen-client-urls http://10.240.0.17:2379 \
    --discovery ${DISCOVERY} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
 
 
etcd --data-dir=data.etcd --name n2 \
    --initial-advertise-peer-urls http://10.240.0.18:2380 --listen-peer-urls http://10.240.0.18:2380 \
    --advertise-client-urls http://10.240.0.18:2379 --listen-client-urls http://10.240.0.18:2379 \
    --discovery ${DISCOVERY} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
 
 
etcd --data-dir=data.etcd --name n3 \
    --initial-advertise-peer-urls http://10.240.0.19:2380 --listen-peer-urls http://10.240.0.19:2380 \
    --advertise-client-urls http://10.240.0.19:2379 --listen-client-urls http:/10.240.0.19:2379 \
    --discovery ${DISCOVERY} \
    --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```
到此etcd集群就搭建起来了，可以使用etcdctl来连接etcd。

```sh
export ETCDCTL_API=3
HOST_1=10.240.0.17
HOST_2=10.240.0.18
HOST_3=10.240.0.19
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379
```

> etcdctl --endpoints=$ENDPOINTS member lis

