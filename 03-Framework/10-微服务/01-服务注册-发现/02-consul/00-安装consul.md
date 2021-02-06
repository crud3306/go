

安装Consul
=============


1、下载Consul
-------------
```sh
wget https://releases.hashicorp.com/consul/0.7.5/consul_0.7.5_linux_amd64.zip

#若无wget，请先安装
#yum -y install wget
```

也可以访问 https://releases.hashicorp.com/consul，选择自已需要的consul版本



2、解压consul_0.7.5_linux_amd64.zip
-------------
```sh
unzip consul_0.7.5_linux_amd64.zip

#有可能会出现-bash: unzip: 未找到命令,解决方案
#yum -y install unzip
```


3、执行以下 ./consul 看是否安装成功（是一个启动文件，不是一个目录）
-------------
> .consul



4、启动consul
-------------
我的ip地址是192.168.100.11
```sh
./consul agent -dev -ui -node=consul-dev -client=192.168.100.11
#consul agent -server -bootstrap-expect=1 -data-dir=data -node=consul -bind=x192.168.100.11-ui -client=0.0.0.0

#关闭临时防火墙
systemctl stop firewalld
```

输出
```sh
==> Starting Consul agent...
==> Starting Consul agent RPC...
==> Consul agent running!
           Version: 'v0.7.5'
           Node ID: '564da851-08d0-5ca7-732d-220ea19adae9'
         Node name: 'consul-dev'
        Datacenter: 'dc1'
            Server: true (bootstrap: false)
       Client Addr: 10.11.3.161 (HTTP: 8500, HTTPS: -1, DNS: 8600, RPC: 8400)
      Cluster Addr: 192.168.100.11 (LAN: 8301, WAN: 8302)
    Gossip encrypt: false, RPC-TLS: false, TLS-Incoming: false
             Atlas: <disabled>

==> Log data will now stream in as it occurs:

    2017/07/02 19:14:43 [DEBUG] Using unique ID "564da851-08d0-5ca7-732d-220ea19adae9" from host as node ID
    2017/07/02 19:14:43 [INFO] raft: Initial configuration (index=1): [{Suffrage:Voter ID:192.168.100.11:8300 Address:192.168.100.11:8300}]
    2017/07/02 19:14:43 [INFO] raft: Node at 192.168.100.11:8300 [Follower] entering Follower state (Leader: "")
```


5、访问consul
-------------
http://192.168.100.11:8500


启动时如果指定了-bind参数，则可以通过web查看consul运行状况：http://192.168.100.11:8500/ui/dc1/nodes




