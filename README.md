# GoWeb

GoWeb 是个简单的HTTP Handler绑定的代码生成工具。目的是减少人为编写代码，减少反射的使用，把资源让给需要使用反射的地方。

## 安装

使用go工具安装，GoWeb默认安装在$GOPATH

```shell
go get "github.com/YiCodes/goweb/goweb"
go get "github.com/YiCodes/goweb/web"
```

## 使用方法

创建一个go包（例如:handlers），并且在包内创建handler.go 文件

```handler.go
package handlers

import (
    "fmt"
)

type Order struct {
    OrderId string
}

func Home() string {
    return "home."
}

func Hello(name string) string {
    return fmt.Sprintf("hello %v", name)
}

func AddOrder(order *Order) {
}

func GetOrder() *Order {
    return &Order{}
}

// 也可以带http.ResponseWriter, *http.Request类型的参数。
func CustomHandle(response http.ResponseWriter) {
    response.Write([]byte("hello world."))
}

// 也可以带http.ResponseWriter, *http.Request类型的参数。
func CustomHandle2(request *http.Request, response http.ResponseWriter) {
    key := request.FormValue("key")
    response.Write([]byte(key))
}
```

## 生成代码

在命令行输入

```shell
goweb -in="handlers"
```

所有handlers包内的对外公开(exported)方法都是视为Http Handler，会生成相应的绑定代码。

执行成功后，会在handlers内生成一个<包名>.gen.go的文件，里面只有一个ConfigServeMuxHandler的方法，用于设置http.ServeMux。

接着在main.go中使用

```main.go

import (
    "handlers"
)

func main() {
    mux := http.NewServeMux()

    handlerOpts := web.NewHandleSetupOptions()

    // 默认所有方法的url是<host>/<方法名>，也可以设置HandleSetupOptions.Route，自定义url
    handlerOpts.Route["Home"] = "/"

    handlers.ConfigServeMuxHandler(mux, handlerOpts)

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## 例子

例子见[这里](https://github.com/YiCodes/goweb/tree/master/example/webapi)