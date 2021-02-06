

封装
-------------
```golang
package service

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
)

var ftpClient *QFtp
var ftpClientLock sync.Mutex

// NewFtpClient 初始ftp客户端
func NewFtpClient() (*QFtp, error) {
	ftpClientLock.Lock()
	defer ftpClientLock.Unlock()

	// LoadFtpConfig 自行实现
	config := LoadFtpConfig()

	c, err := ftp.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(config.User, config.Password)
	if err != nil {
		return nil, err
	}

	ftpClient := &QFtp{
		ftpConn: c,
	}

	return ftpClient, nil
}

// QFtp 文件上传
type QFtp struct {
	ftpConn *ftp.ServerConn
}

// Upload 上传
func (f *QFtp) Upload(path string, r io.Reader) error {
	if err := f.ftpConn.Stor(path, r); err != nil {
		return err
	}

	return nil
}

// Close 关闭连接
func (f *QFtp) Close() error {
	if err := f.ftpConn.Quit(); err != nil {
		return err
	}

	return nil
}
```

use
```golang
import (
	"errors"
	"fmt"
	"xxx/service"
	"path/filepath"
)

// UploadToFtp 上传文件至ftp
func UploadToFtp(uploadFile string) error {
	fileName := filepath.Base(uploadFile)

	var file *os.File
	file, err := os.Open(uploadFile)
	if err != nil {
		return err
	}
	defer file.Close()

	ftpClient, err := service.NewFtpClient()
	if err != nil {
		return err
	}
	defer ftpClient.Close()

	if err := ftpClient.Upload(fileName, file); err != nil {
		return err
	}

	return nil
}


uploadFile := "/xx/xx/xx1" 
if err := UploadToFtp(uploadFile); err != nil {
	fmt.Println("err", err.Error())
}

fmt.Println("successs")
```