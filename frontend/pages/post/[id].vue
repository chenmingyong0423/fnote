<template>
    <div class="bg-#F0F2F5 dark_bg_black pt100 pb30 px25">
        <el-row :gutter="40">
            <el-col :span="18">
                <div class=" rounded-15 bg-#fff p30">
                    <!-- 文章标题 -->
                    <div class="text-center text-42 font-500">{{ data2.data.title }}</div>
                    <!-- 作者信息 -->
                    <div class="text-center text-16 c-#999/80 my20 ">
                        <el-space :size="4">
                            <div>{{ data.data.author }}</div>
                            <div>{{ data.data.createTime }}</div>
                            <div>·</div>
                            <div>阅读 </div>
                            <div>{{ data.data.visit }}</div>
                        </el-space>
                    </div>
                    <v-md-preview :text="data2.data.content"></v-md-preview>
                    <!-- 点赞 -->
                    <div
                        class=" mx-auto cursor-pointer text-22 w45 h45 flex justify-center items-center rounded-50% hover:bg-#e5e5e5 active:bg-#999 active:c-#fff">
                        <div class="i-grommet-icons:like"></div>
                    </div>
                </div>
                <PostComments class="mt25" />
            </el-col>
            <el-col :span="6">
                <el-affix :offset="90">
                    <div class="w-full">
                        <Anchor :text="data2.data.content"></Anchor>
                    </div>
                </el-affix>
            </el-col>
        </el-row>
    </div>
</template>

<script lang="ts" setup>



// const router = useRouter()
// const route = useRoute()

const data = {
    "code": 200,
    "message": "OK",
    "data": {
        "sug": "post1",
        "author": "陈明勇",
        "title": "Go语言啊",
        "summary": "Summary 1",
        "cover_img": "/images/cover1.jpg",
        "category": "A",
        "tags": [
            "a",
            "b"
        ],
        "likeCount": 2,
        "comments": 3,
        "visit": 0,
        "priority": 1,
        "createTime": 1692806408145,
        "content": "# 切片\r\n哈哈",
        "meta_description": "Description 1",
        "meta_keywords": "Keyword 1",
        "word_count": 0,
        "allow_comment": true,
        "update_time": 1692806408143,
        "is_liked": false
    }
}

const data2 = {
    "code": 200,
    "message": "OK",
    "data": {
        "id": "4d70f912-f9a9-49db-ae60-947eb08d356e",
        "author": "陈明勇",
        "title": "一文掌握 Go 并发模式 Context 上下文",
        "summary": "本文详细介绍了 Go 语言中的 Context 上下文，包括核心方法、创建方式以及应用场景等方面的内容。",
        "content": "\n\u003e 扫码关注公众号，手机阅读更方便\n\u003e \n\u003e ![Go技术干货](https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/img/Go%E6%8A%80%E6%9C%AF%E5%B9%B2%E8%B4%A7.jpg) \n\u003e Go version → 1.20.4\n# 前言\n\u003e Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.\u003csup\u003e[1]\u003c/sup\u003e\n\n`Go` 在 **1.7** 引入了 `context` 包，目的是为了在不同的 `goroutine` 之间或跨 `API` 边界传递超时、取消信号和其他请求范围内的值（与该请求相关的值。这些值可能包括用户身份信息、请求处理日志、跟踪信息等等）。\n\n在 `Go` 的日常开发中，`Context` 上下文对象无处不在，无论是处理网络请求、数据库操作还是调用 `RPC` 等场景下，都会使用到 `Context`。那么，你真的了解它吗？熟悉它的正确用法吗？了解它的使用注意事项吗？喝一杯你最喜欢的饮料，随着本文一探究竟吧。\n\n# Context 接口\n`context` 包在提供了一个用于跨 `API` 边界传递超时、取消信号和其他请求范围值的通用数据结构。它定义了一个名为 `Context` 的接口，该接口包含一些方法，用于在多个 `Goroutine` 和函数之间传递请求范围内的信息。\n\n以下是 `Context` 接口的定义：\n```go\ntype Context interface {\n    Deadline() (deadline time.Time, ok bool)\n\n    Done() \u003c-chan struct{}\n\n    Err() error\n\n    Value(key any) any\n}\n```\n# Context 的核心方法\n\n![Context 的核心方法.jpg](https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/img/Context%20%E7%9A%84%E6%A0%B8%E5%BF%83%E6%96%B9%E6%B3%95.jpg)\n\n`Context` 接口中有四个核心方法：`Deadline()`、`Done()`、`Err()`、`Value()`。\n\n## Deadline()\n`Deadline() (deadline time.Time, ok bool)` 方法返回 `Context` 的截止时间，表示在这个时间点之后，`Context` 会被自动取消。如果 `Context` 没有设置截止时间，该方法返回一个零值 `time.Time` 和一个布尔值 `false`。\n```go\ndeadline, ok := ctx.Deadline()\nif ok {\n    // Context 有截止时间\n} else {\n    // Context 没有截止时间\n}\n```\n## Done()\n`Done()` 方法返回一个只读通道，当 `Context` 被取消时，该通道会被关闭。你可以通过监听这个通道来检测 `Context` 是否被取消。如果 `Context` 永不取消，则返回 `nil`。\n\n```go\nselect {\ncase \u003c-ctx.Done():\n    // Context 已取消\ndefault:\n    // Context 尚未取消\n}\n```\n## Err()\n`Err()` 方法返回一个 `error` 值，表示 `Context` 被取消时产生的错误。如果 `Context` 尚未取消，该方法返回 `nil`。\n```go\nif err := ctx.Err(); err != nil {\n    // Context 已取消，处理错误\n}\n```\n## Value()\n`Value(key any) any` 方法返回与 `Context` 关联的键值对，一般用于在 `Goroutine` 之间传递请求范围内的信息。如果没有关联的值，则返回 `nil`。\n```go\nvalue := ctx.Value(key)\nif value != nil {\n    // 存在关联的值\n}\n```\n# Context 的创建方式\n\n![Context 的创建方式.jpg](https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/img/Context%20%E7%9A%84%E5%88%9B%E5%BB%BA%E6%96%B9%E5%BC%8F.jpg)\n## context.Background()\n`context.Background()` 函数返回一个非 `nil` 的空 `Context`，它没有携带任何的值，也没有取消和超时信号。通常作为根 `Context` 使用。\n```go\nctx := context.Background()\n```\n## context.TODO()\n\n`context.TODO()` 函数返回一个非 `nil` 的空 `Context`，它没有携带任何的值，也没有取消和超时信号。虽然它的返回结果和 `context.Background()` 函数一样，但是它们的使用场景是不一样的，如果不确定使用哪个上下文时，可以使用 `context.TODO()`。\n```go\nctx := context.TODO()\n```\n## context.WithValue()\n`context.WithValue(parent Context, key, val any)` 函数接收一个父 `Context` 和一个键值对 `key`、`val`，返回一个新的子 `Context`，并在其中添加一个 `key-value` 数据对。\n```go\nctx := context.WithValue(parentCtx, \"username\", \"陈明勇\")\n```\n## context.WithCancel()\n`context.WithCancel(parent Context) (ctx Context, cancel CancelFunc)` 函数接收一个父 `Context`，返回一个新的子 `Context` 和一个取消函数，当取消函数被调用时，子 `Context` 会被取消，同时会向子 `Context` 关联的 `Done()` 通道发送取消信号，届时其衍生的子孙 `Context` 都会被取消。这个函数适用于手动取消操作的场景。\n```go\nctx, cancelFunc := context.WithCancel(parentCtx)  \ndefer cancelFunc()\n```\n## context.WithCancelCause() 与 context.Cause()\n`context.WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc)` 函数是 `Go 1.20` 版本才新增的，其功能类似于 `context.WithCancel()`，但是它可以设置额外的取消原因，也就是 `error` 信息，返回的 `cancel` 函数被调用时，需传入一个 `error` 参数。\n```go\nctx, cancelFunc := context.WithCancelCause(parentCtx)\ndefer cancelFunc(errors.New(\"原因\"))\n```\n`context.Cause(c Context) error` 函数用于返回取消 `Context` 的原因，即错误值 `error`。如果是通过 `context.WithCancelCause()` 函数返回的取消函数 `cancelFunc(myErr)` 进行的取消操作，我们可以获取到 `myErr` 的值。否则，我们将得到与 `c.Err()` 相同的返回值。如果 `Context` 尚未被取消，将返回 `nil`。\n```go\nerr := context.Cause(ctx)\n```\n## context.WithDeadline()\n`context.WithDeadline(parent Context, d time.Time) (Context, CancelFunc)` 函数接收一个父 `Context` 和一个截止时间作为参数，返回一个新的子 `Context`。当截止时间到达时，子 `Context` 其衍生的子孙 `Context` 会被自动取消。这个函数适用于需要在特定时间点取消操作的场景。\n```go\ndeadline := time.Now().Add(time.Second * 2)\nctx, cancelFunc := context.WithTimeout(parentCtx, deadline)\ndefer cancelFunc()\n```\n## context.WithTimeout()\n`context.WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)` 函数和 `context.WithDeadline()` 函数的功能是一样的，其底层会调用 `WithDeadline()` 函数，只不过其第二个参数接收的是一个超时时间，而不是截止时间。这个函数适用于需要在一段时间后取消操作的场景。\n```go\nctx, cancelFunc := context.WithTimeout(parentCtx, time.Second * 2)\ndefer cancelFunc()\n```\n\n# Context 的使用场景\n## 传递共享数据\n编写中间件函数，用于向 `HTTP` 处理链中添加处理请求 `ID` 的功能。\n```go\ntype key int\n\nconst (\n   requestIDKey key = iota\n)\n\nfunc WithRequestId(next http.Handler) http.Handler {\n   return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {\n      // 从请求中提取请求ID和用户信息\n      requestID := req.Header.Get(\"X-Request-ID\")\n\n      // 创建子 context，并添加一个请求 Id 的信息\n      ctx := context.WithValue(req.Context(), requestIDKey, requestID)\n\n      // 创建一个新的请求，设置新 ctx\n      req = req.WithContext(ctx)\n\n      // 将带有请求 ID 的上下文传递给下一个处理器\n      next.ServeHTTP(rw, req)\n   })\n}\n```\n首先，我们从请求的头部中提取请求 `ID`。然后使用 `context.WithValue` 创建一个子上下文，并将请求 `ID` 作为键值对存储在子上下文中。接着，我们创建一个新的请求对象，并将子上下文设置为新请求的上下文。最后，我们将带有请求 `ID` 的上下文传递给下一个处理器。\n这样，通过使用 `WithRequestId` 中间件函数，我们可以在处理请求的过程中方便地获取和使用请求 `ID`，例如在 **日志记录、跟踪和调试等方面**。\n## 传递取消信号，结束任务\n启动一个工作协程，接收到取消信号就停止工作。\n```go\npackage main\n\nimport (\n   \"context\"\n   \"fmt\"\n   \"time\"\n)\n\nfunc main() {\n   ctx, cancelFunc := context.WithCancel(context.Background())\n   go Working(ctx)\n\n   time.Sleep(3 * time.Second)\n   cancelFunc()\n\n   // 等待一段时间，以确保工作协程接收到取消信号并退出\n   time.Sleep(1 * time.Second)\n}\n\nfunc Working(ctx context.Context) {\n   for {\n      select {\n      case \u003c-ctx.Done():\n         fmt.Println(\"下班啦...\")\n         return\n      default:\n         fmt.Println(\"陈明勇正在工作中...\")\n      }\n   }\n}\n```\n执行结果\n```\n······\n······\n陈明勇正在工作中...\n陈明勇正在工作中...\n陈明勇正在工作中...\n陈明勇正在工作中...\n陈明勇正在工作中...\n下班啦...\n```\n在上面的示例中，我们创建了一个 `Working` 函数，它会不断执行工作任务。我们使用 `context.WithCancel` 创建了一个上下文 `ctx` 和一个取消函数 `cancelFunc`。然后，启动了一个工作协程，并将上下文传递给它。\n\n在主函数中，需要等待一段时间（**3** 秒）模拟业务逻辑的执行。然后，调用取消函数 `cancelFunc`，通知工作协程停止工作。工作协程在每次循环中都会检查上下文的状态，一旦接收到取消信号，就会退出循环。\n\n最后，等待一段时间（**1** 秒），以确保工作协程接收到取消信号并退出。\n## 超时控制\n模拟耗时操作，超时控制。\n```go\npackage main\n\nimport (\n   \"context\"\n   \"fmt\"\n   \"time\"\n)\n\nfunc main() {\n   // 使用 WithTimeout 创建一个带有超时的上下文对象\n   ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)\n   defer cancel()\n\n   // 在另一个 goroutine 中执行耗时操作\n   go func() {\n      // 模拟一个耗时的操作，例如数据库查询\n      time.Sleep(5 * time.Second)\n      cancel()\n   }()\n\n   select {\n   case \u003c-ctx.Done():\n      fmt.Println(\"操作已超时\")\n   case \u003c-time.After(10 * time.Second):\n      fmt.Println(\"操作完成\")\n   }\n}\n```\n执行结果\n```\n操作已超时\n```\n在上面的例子中，首先使用 `context.WithTimeout()` 创建了一个带有 **3** 秒超时的上下文对象 `ctx, cancel := context.WithTimeout(ctx, 3*time.Second)`。\n\n接下来，在一个新的 `goroutine` 中执行一个模拟的耗时操作，例如等待 **5** 秒钟。当耗时操作完成后，调用 `cancel()` 方法来取消超时上下文。\n\n\n最后，在主 `goroutine` 中使用 `select` 语句等待超时上下文的完成信号。如果在 **3** 秒内耗时操作完成，那么会输出 \"操作完成\"。如果超过了 **3** 秒仍未完成，超时上下文的 `Done()` 通道会被关闭，输出 \"操作已超时\"。\n\n# 使用 Context 的一些规则\n使用 `Context` 上下文，应该遵循以下规则，以保持包之间的接口一致，并使静态分析工具能够检查上下文传播:\n- 不要在结构类型中加入 `Context` 参数，而是将它显式地传递给需要它的每个函数，并且它应该是第一个参数，通常命名为 `ctx`:\n    ```go\n    func DoSomething(ctx context.Context, arg Arg) error {\n            // ... use ctx ...\n    }\n    ```\n- 即使函数允许，也不要传递 `nil Context`。如果不确定要使用哪个 `Context`，建议使用 `context.TODO()`。\n\n- 仅将 `Context` 的值用于传输进程和 `api` 的请求作用域数据，不能用于向函数传递可选参数。\u003csup\u003e[1]\u003c/sup\u003e\n# 小结\n本文详细介绍了 `Go` 语言中的 `Context` 上下文，通过阅读本文，相信你们对 `Context` 的功能和使用场景有所了解。同时，你们也应该能够根据实际需求选择最合适的 `Context` 创建方式，并且根据规则，正确、高效地使用它。\n\n# 参考资料\n[1] https://pkg.go.dev/context@go1.20.4\n",
        "tags": [
            "Go"
        ],
        "publishDate": "2023-05-17 02:28:58",
        "category": "后端",
        "columns": [
            "Go 并发"
        ],
        "cover": "https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/img/%E5%B0%81%E9%9D%A2-%E4%B8%80%E6%96%87%E6%8E%8C%E6%8F%A1%20Go%20%E5%B9%B6%E5%8F%91%E6%A8%A1%E5%BC%8F%20Context%20%E4%B8%8A%E4%B8%8B%E6%96%87.jpg",
        "likesCount": 2,
        "viewsCount": 201,
        "commentInfo": [
            {
                "id": "8fbe8dfd-32bd-4901-9eae-4d7f0287b799",
                "content": "超时控制代码的 select 没看明白，case \u003c-time.After()是什么意思？",
                "username": "威",
                "commentTime": "2023-07-27 09:16:55",
                "replies": [
                    {
                        "id": "3050a13b-7639-4f5e-aade-d48430f3a0ef",
                        "commentId": "8fbe8dfd-32bd-4901-9eae-4d7f0287b799",
                        "content": "`case \u003c-time.After(10 * time.Second)`\n\n表示，只要不超时（不触发第一个 `case`），每 `10 s` 都会触发这个 `case`，然后执行 `case` 里的逻辑 `fmt.Println(\"操作完成\")`\n\n不过不建议这样使用哈，本文这个例子举的不好，容易造成内存溢出，推荐使用以下用法：\n\n```go\n    ticker := time.NewTicker(time.Second * 10)\n    for {\n       select {\n           case \u003c-ctx.Done():\n              fmt.Println(\"操作已超时\")\n              return\n          case \u003c-ticker.C:\n             fmt.Println(\"操作完成\")\n        }\n    }\n```",
                        "username": "陈明勇",
                        "replyToId": "",
                        "replyTo": "",
                        "replyTime": "2023-07-27 22:04:35"
                    },
                    {
                        "id": "ef448c18-fed6-41e3-8470-c6f217c30c96",
                        "commentId": "8fbe8dfd-32bd-4901-9eae-4d7f0287b799",
                        "content": "懂了。就是如果没有 cancel 信号，程序会在 select 处阻塞 10 秒，然后打印“操作完成”",
                        "username": "威",
                        "replyToId": "3050a13b-7639-4f5e-aade-d48430f3a0ef",
                        "replyTo": "陈明勇",
                        "replyTime": "2023-07-27 22:30:12"
                    }
                ]
            },
            {
                "id": "a87b6410-ca0a-4040-9d76-d556897303f9",
                "content": "不错不错",
                "username": "zq",
                "commentTime": "2023-05-17 10:20:24"
            }
        ],
        "isLike": false
    }
}
definePageMeta({
    layout: "home"
})
</script>

<style scoped></style>