
go 程序 启动docker容器
---------------
```golang
package main
 
import (
    "io"
    "log"
    "os"
    "time"
 
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/mount"
    "github.com/docker/docker/client"
    "github.com/docker/go-connections/nat"
    "golang.org/x/net/context"
)
 
const (
    imageName     string   = "my-gin:latest"                      //镜像名称
    containerName string   = "mygin-latest"                       //容器名称
    indexName     string   = "/" + containerName                  //容器索引名称，用于检查该容器是否存在是使用
    cmd           string   = "./ginDocker2"                       //运行的cmd命令，用于启动container中的程序
    workDir       string   = "/go/src/ginDocker2"                 //container工作目录
    openPort      nat.Port = "7070"                               //container开放端口
    hostPort      string   = "7070"                               //container映射到宿主机的端口
    containerDir  string   = "/go/src/ginDocker2"                 //容器挂在目录
    hostDir       string   = "/home/youngblood/Go/src/ginDocker2" //容器挂在到宿主机的目录
    n             int      = 5                                    //每5s检查一个容器是否在运行
 
)
 
func main() {
    ctx := context.Background()
    cli, err := client.NewEnvClient()
    defer cli.Close()
    if err != nil {
        panic(err)
    }
    checkAndStartContainer(ctx, cli)
}
 
//创建容器
func createContainer(ctx context.Context, cli *client.Client) {
    //创建容器
    cont, err := cli.ContainerCreate(ctx, &container.Config{
        Image:      imageName,     //镜像名称
        Tty:        true,          //docker run命令中的-t选项
        OpenStdin:  true,          //docker run命令中的-i选项
        Cmd:        []string{cmd}, //docker 容器中执行的命令
        WorkingDir: workDir,       //docker容器中的工作目录
        ExposedPorts: nat.PortSet{
            openPort: struct{}{}, //docker容器对外开放的端口
        },
    }, &container.HostConfig{
        PortBindings: nat.PortMap{
            openPort: []nat.PortBinding{nat.PortBinding{
                HostIP:   "0.0.0.0", //docker容器映射的宿主机的ip
                HostPort: hostPort,  //docker 容器映射到宿主机的端口
            }},
        },
        Mounts: []mount.Mount{ //docker 容器目录挂在到宿主机目录
            mount.Mount{
                Type:   mount.TypeBind,
                Source: hostDir,
                Target: containerDir,
            },
        },
    }, nil, containerName)
    if err == nil {
        log.Printf("success create container:%s\n", cont.ID)
    } else {
        log.Println("failed to create container!!!!!!!!!!!!!")
    }
}
 
//启动容器
func startContainer(ctx context.Context, containerID string, cli *client.Client) error {
    err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
    if err == nil {
        log.Printf("success start container:%s\n", containerID)
    } else {
        log.Printf("failed to start container:%s!!!!!!!!!!!!!\n", containerID)
    }
    return err
}
 
//将容器的标准输出输出到控制台中
func printConsole(ctx context.Context, cli *client.Client, id string) {
    //将容器的标准输出显示出来
    out, err := cli.ContainerLogs(ctx, id, types.ContainerLogsOptions{ShowStdout: true})
    if err != nil {
        panic(err)
    }
    io.Copy(os.Stdout, out)
 
    //容器内部的运行状态
    status, err := cli.ContainerStats(ctx, id, true)
    if err != nil {
        panic(err)
    }
    io.Copy(os.Stdout, status.Body)
}
 
//检查容器是否存在并启动容器
func checkAndStartContainer(ctx context.Context, cli *client.Client) {
    for {
        select {
        case <-isRuning(ctx, cli):
            //该container没有在运行
            //获取所有的container查看该container是否存在
            contTemp := getContainer(ctx, cli, true)
            if contTemp.ID == "" {
                //该容器不存在，创建该容器
                log.Printf("the container name[%s] is not exists!!!!!!!!!!!!!\n", containerName)
                createContainer(ctx, cli)
            } else {
                //该容器存在，启动该容器
                log.Printf("the container name[%s] is exists\n", containerName)
                startContainer(ctx, contTemp.ID, cli)
            }
 
        }
    }
}
 
//获取container
func getContainer(ctx context.Context, cli *client.Client, all bool) types.Container {
    containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{All: all})
    if err != nil {
        panic(err)
    }
    var contTemp types.Container
    //找出名为“mygin-latest”的container并将其存入contTemp中
    for _, v1 := range containerList {
        for _, v2 := range v1.Names {
            if v2 == indexName {
                contTemp = v1
                break
            }
        }
    }
    return contTemp
}
 
//容器是否正在运行
func isRuning(ctx context.Context, cli *client.Client) <-chan bool {
    isRun := make(chan bool)
    var timer *time.Ticker
    go func(ctx context.Context, cli *client.Client) {
        for {
            //每n s检查一次容器是否运行
 
            timer = time.NewTicker(time.Duration(n) * time.Second)
            select {
            case <-timer.C:
                //获取正在运行的container list
                log.Printf("%s is checking the container[%s]is Runing??", os.Args[0], containerName)
                contTemp := getContainer(ctx, cli, false)
                if contTemp.ID == "" {
                    log.Print(":NO")
                    //说明container没有运行
                    isRun <- true
                } else {
                    log.Print(":YES")
                    //说明该container正在运行
                    go printConsole(ctx, cli, contTemp.ID)
                }
            }
 
        }
    }(ctx, cli)
    return isRun
}
```