

python快速入门教程

本章主要介绍在python中grpc的基本用法。

在开始之前先安装一些东西。


前置知识点：
- protobuf语法
- grpc基本概念


1.环境安装
================
grpc支持Python 2.7 或者 Python 3.4 或者更高的版本。

首先确认pip的版本大于9.0.1
> python -m pip install --upgrade pip

如果因为某些原因无法直接升级pip, 可以选择在virtualenv中运行grpc例子。
```sh
$ python -m pip install virtualenv
$ virtualenv venv
$ source venv/bin/activate
$ python -m pip install --upgrade pip
```

1.1. 安装gRPC
-----------------
> python -m pip install grpcio

在mac系统中可能会出现下面的错误：
```sh
$ OSError: [Errno 1] Operation not permitted: '/tmp/pip-qwTLbI-uninstall/System/Library/Frameworks/Python.framework/Versions/2.7/Extras/lib/python/six-1.4.1-py2.7.egg-info'
```

可以使用下面的命令解决问题：
> python -m pip install grpcio --ignore-installed


1.2. 安装gRPC工具
-----------------
gRPC工具包括protocol buffer编译器（protoc）和 python代码生成插件，python代码生成插件通过.proto服务定义文件，生成python的grpc服务端和客户端代码。

> python -m pip install grpcio-tools



2.下载例子代码
================

我们直接从grpc的github地址，下载grpc代码，里面包含了python的例子。
```sh
$ # 将grpc代码clone到本地
$ git clone -b v1.23.0 https://github.com/grpc/grpc

$ # 切换到python的helloworld例子目录。
$ cd grpc/examples/python/helloworld
```


3.运行grpc应用
================
首先运行服务端

> python greeter_server.py


打开另外一个命令窗口，运行客户端

> python greeter_client.py


到目前为止，我们已经运行了服务端和客户端的例子。



4.更新rpc服务
================
现在我们添加一个新的rpc接口，看看python的服务端和客户端代码怎么修改。

grpc的接口是通过protocol buffers定义的，通常都保存在.proto文件中。

如果你打开，前面例子的服务端和客户端代码看，你会发现有一个SayHello的方法，这个方法接受HelloRequest参数，并返回HelloReply参数，下面看看对应的.proto文件，是如何定义rpc接口的。

打开：examples/protos/helloworld.proto 协议文件。
```golang
// 定义Greeter，你可以当成类
service Greeter {
  // 定义一个rpc方法SayHello，接受HelloRequest消息，返回HelloReply消息
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// 定义请求参数，HelloRequest消息
message HelloRequest {
  string name = 1;
}

// 定义响应参数，HelloReply消息
message HelloReply {
  string message = 1;
}
```

现在我们添加一个新的rpc方法，现在.proto定义变成这样。
```golang
// 定义Greeter，你可以当成类
service Greeter {
  // SayHello方法
  rpc SayHello (HelloRequest) returns (HelloReply) {}
  // 定义一个SayHelloAgain方法，接受HelloRequest消息，返回HelloReply消息
  rpc SayHelloAgain (HelloRequest) returns (HelloReply) {}
}

// 定义请求参数，HelloRequest消息
message HelloRequest {
  string name = 1;
}

// 定义响应参数，HelloReply消息
message HelloReply {
  string message = 1;
}
```
修改后记得保存.proto文件。



5.生成grpc代码
================
现在我们需要根据新定义的.proto文件，生成新的python代码（这里其实只是生成一个类库，相当于更新类库，我们的业务代码实现刚才定义的rpc方法即可）

切换到examples/python/helloworld目录，运行命令：

> python -m grpc_tools.protoc -I../../protos --python_out=. --grpc_python_out=. ../../protos/helloworld.proto

命令说明：
```sh
#-I proto协议文件目录  
#--python_out和--grpc_python_out 生成python代码的目录  

#命令最后面的参数是proto协议文件名
```
命令执行后生成helloworld_pb2.py文件和helloworld_pb2_grpc.py文件。

helloworld_pb2.py - 主要包含proto文件定义的消息类。  
helloworld_pb2_grpc.py - 包含服务端和客户端代码  



6.修改并运行应用程序
================
现在我们已经根据proto文件，生成新的python类库，但是我们还没实现新定义的rpc方法，下面介绍服务端和客户端如果升级代码。

6.1. 修改服务端代码
-----------------
在同样目录打开greeter_server.py文件，实现类似如下代码。
```python
class Greeter(helloworld_pb2_grpc.GreeterServicer):
  # 实现SayHello方法
  def SayHello(self, request, context):
    return helloworld_pb2.HelloReply(message='Hello, %s!' % request.name)
  # 实现SayHelloAgain方法
  def SayHelloAgain(self, request, context):
    return helloworld_pb2.HelloReply(message='Hello again, %s!' % request.name)
...
```


6.2. 修改客户端代码
-----------------
在同样的目录打开greeter_client.py文件，实现代码如下：  
```py
def run():
  # 配置grpc服务端地址
  channel = grpc.insecure_channel('localhost:50051')
  stub = helloworld_pb2_grpc.GreeterStub(channel)
  # 请求服务端的SayHello方法
  response = stub.SayHello(helloworld_pb2.HelloRequest(name='you'))
  print("Greeter client received: " + response.message)
  # 请求服务端的SayHelloAgain方法
  response = stub.SayHelloAgain(helloworld_pb2.HelloRequest(name='you'))
  print("Greeter client received: " + response.message)
```


6.3. 运行代码
-----------------
在examples/python/helloworld目录下面，首先运行服务端：  
> python greeter_server.py

运行客户端  
> python greeter_client.py



