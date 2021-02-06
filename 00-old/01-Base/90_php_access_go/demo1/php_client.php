<?php

$descriptorspec = [
  0 => ["pipe", "r"],
  1 => ["pipe", "w"]
];

// var_dump($descriptorspec);

// 注意此处第一个参数是 go_server.go 通过 (go build) 生成的二进制可执行文件
$handle = proc_open(
  '/Users/qianyouming/go/src/php_access_go/demo1/go_server',
  $descriptorspec,
  $pipes
);

// var_dump($handle, $pipes);
// $fp = fopen("title.txt", "rb");

// while (!feof($fp)) {
//   fwrite($pipes['0'], trim(fgets($fp))."\n");
//   echo fgets($pipes[1]);
// }

 $flag = fwrite($pipes['0'], "这是一个测试文本\n");
 // var_dump($flag);

 echo fgets($pipes[1]);

// fclose($pipes['0']);
// fclose($pipes['1']);
// proc_close($handle);

