

参考：
----------
https://colobu.com/2017/06/19/advanced-command-execution-in-Go-with-os-exec/



```golang
import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

// TarDir 打包指定的目录
func TarDir(srcFile, dstFile string) error {
	// 源文件路径
	srcFilePath := filepath.Dir(srcFile)
	// 源文件名
	srcFileName := filepath.Base(srcFile)
	// 打包命令
	command := fmt.Sprintf("cd %s && tar -zcf %s %s", srcFilePath, dstFile, srcFileName)
	cmd := exec.Command("/bin/bash", "-c", command)

	return Command(cmd)
}

// Command 执行shell
func Command(cmd *exec.Cmd) error {
	var stdOut, stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if err := cmd.Run(); err != nil {
		return errors.New(stdErr.String())
	}

	return nil
}

func Command2(cmd *exec.Cmd) error {
	// cmd := "cat /proc/cpuinfo | egrep '^model name' | uniq | awk '{print substr($0, index($0,$4))}'"
	// cmd := exec.Command("/bin/bash", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}

	fmt.Println(string(out))
	return nil
}```