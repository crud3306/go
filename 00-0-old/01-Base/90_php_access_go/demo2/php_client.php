<?php


# 注意，执行此php前，需先启动 go_server.go

$msg = "msg";

// 创建 连接 发送消息 接收响应 关闭连接
$socket = socket_create(AF_UNIX, SOCK_STREAM, 0);

socket_connect($socket, '/tmp/keyword_match.sock');

socket_send($socket, $msg, strlen($msg), 0);

$response = socket_read($socket, 1024);

socket_close($socket);

// 有值则为匹配成功
if (strlen($response) > 3) {
    var_dump($response);
}