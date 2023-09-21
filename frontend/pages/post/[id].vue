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
                            <div>
                                {{ dayjs(data.data.createTime).format('YYYY-MM-DD HH:mm:ss') }}
                            </div>
                            <div>·</div>
                            <div>阅读 </div>
                            <div>{{ data.data.visit }}</div>
                        </el-space>
                    </div>
                    <v-md-preview :text="data2.data.content" @copy-code-success="handleCopyCodeSuccess"></v-md-preview>
                    <!-- 点赞 -->
                    <div
                        class=" mx-auto cursor-pointer text-22 w45 h45 flex justify-center items-center rounded-50% hover:bg-#e5e5e5 active:bg-#999 active:c-#fff">
                        <div class="i-grommet-icons:like"></div>
                    </div>
                </div>
                <div class="rounded-15 bg-#fff pt20 pb80 px30 mt25">
                    <!-- 发布评论 -->
                    <SubmitComments />
                    <div class="text-16">{{ data2.data.commentInfo.length + ' 评论' }}</div>
                    <PostCommentsFather v-for="item, index in data2.data.commentInfo" :key="index" class=""
                        :commentInfo="item" />
                </div>
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
import { ElMessage } from 'element-plus'
import { dayjs } from 'element-plus'
import type { } from 'element-plus'
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
        "id": "198e8383-8c7a-4a1b-af73-7317c556bd66",
        "author": "陈明勇",
        "title": "Go 函数选项模式（Functional Options Pattern）",
        "summary": "本文对 Go 函数选项模式（Functional Options Pattern）进行了详细介绍，并通过封装一个消息结构体的例子，展示了如何使用函数选项模式进行代码实现。",
        "content": "\u003e 扫码关注公众号，手机阅读更方便\n\u003e \n\u003e ![Go技术干货](https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/img/Go%E6%8A%80%E6%9C%AF%E5%B9%B2%E8%B4%A7.jpg) \n\n# 前言\n\n在日常开发中，有些函数可能需要接收许多参数，其中一些参数是必需的，而其他参数则是可选的。当函数参数过多时，函数会变得臃肿且难以理解。此外，如果在将来需要添加更多参数，就必须修改函数签名，这将影响到已有的调用代码。\n\n而函数选项模式（`functional options`）的出现解决了这个问题，本文将对其进行讲解，准备好了吗？准备一杯你最喜欢的饮料或茶，随着本文一探究竟吧。\n\n# 函数选项模式\n\n## 什么是函数选项模式\n\n在 `Go` 语言中，函数选项模式是一种优雅的设计模式，用于处理函数的可选参数。它提供了一种灵活的方式，允许用户在函数调用时传递一组可选参数，而不是依赖于固定数量和顺序的参数列表。\n\n## 函数选项模式的好处\n\n*   易于使用：调用者可以选择性的设置函数参数，而不需要记住参数的顺序和类型；\n*   可读性强：函数选项模式的代码有着自文档化的特点，调用者能够直观地理解代码的功能；\n*   扩展性好：通过添加新的 `Option` 参数选项，函数可以方便地扩展功能，无需修改函数的签名；\n*   函数选项模式可以提供默认参数值，以减少参数传递的复杂性。\n\n# 函数选项模式的实现\n\n函数选项模式的实现一般包含以下几个部分：\n\n*   选项结构体：用于存储函数的配置参数\n*   选项函数类型：接收选项结构体参数的函数\n*   定义功能函数：接收 0 个或多个固定参数和可变的选项函数参数\n*   设置选项的函数：定义多个设置选项的函数，用于设置选项\n\n代码示例：\n\n```go\ntype Message struct {\n   // 标题、内容、信息类型\n   title, message, messageType string\n\n   // 账号\n   account     string\n   accountList []string\n\n   // token\n   token     string\n   tokenList []string\n}\n\ntype MessageOption func(*Message)\n\nfunc NewMessage(title, message, messageType string, opts ...MessageOption) *Message {\n   msg := \u0026Message{\n      title:       title,\n      message:     message,\n      messageType: messageType,\n   }\n\n   for _, opt := range opts {\n      opt(msg)\n   }\n\n   return msg\n}\n\nfunc WithAccount(account string) MessageOption {\n   return func(message *Message) {\n      message.account = account\n   }\n}\n\nfunc WithAccountList(accountList []string) MessageOption {\n   return func(message *Message) {\n      message.accountList = accountList\n   }\n}\n\nfunc WithToken(token string) MessageOption {\n   return func(message *Message) {\n      message.token = token\n   }\n}\n\nfunc WithTokenList(tokenList []string) MessageOption {\n   return func(message *Message) {\n      message.tokenList = tokenList\n   }\n}\n\nfunc main() {\n   // 单账号推送\n   _ = NewMessage(\n      \"来自陈明勇的信息\",\n      \"你好，我是陈明勇\",\n      \"单账号推送\",\n      WithAccount(\"123456\"),\n   )\n\n   // 多账号推送\n   _ = NewMessage(\n      \"来自陈明勇的信息\",\n      \"你好，我是陈明勇\",\n      \"多账号推送\",\n      WithAccountList([]string{\"123456\", \"654321\"}),\n   )\n}\n```\n\n上述例子中，使用了函数选项模式来创建 `Message` 结构体，并根据消息类型配置不同消息的属性。\n\n首先定义了 `Message` 结构体，其包含 7 个字段；\n\n其次定义 `MessageOptionm`选项函数类型，用于接收 `Message` 参数的函数；\n\n再次定义 `NewMessage` 函数，用于创建一个 `Message` 指针变量，在 `NewMessage` 函数中，固定参数包括 `title`、`message` 和 `messageType`，它们是必需的参数。然后，通过可选参数 `opts ...MessageOption` 来接收一系列的函数选项；\n\n然后定义了四个选项函数：`WithAccount`、`WithAccountList`、`WithToken` 和 `WithTokenList`。这些选项函数分别用于设置被推送消息的账号、账号列表、令牌和令牌列表。\n\n最后，在 `main` 函数中，展示了两种不同的用法。第一个示例是创建单账号推送的消息，通过调用 `NewMessage` 并传递相应的参数和选项函数（`WithAccount`）来配置消息。第二个示例是创建多账号推送的消息，同样通过调用 `NewMessage` 并使用不同的选项函数（`WithAccountList`）来配置消息。\n\n这种使用函数选项模式的方式可以根据需要消息类型去配置消息的属性，使代码更具灵活性和可扩展性。\n\n# 函数选项模式的缺点\n\n前面提到了函数选项模式的优势（好处），但也必须承认它存在一些缺点。\n\n*   复杂性：函数选项模式引入了更多的类型和概念，需要更多的代码和逻辑来处理。这增加了代码的复杂性和理解的难度，尤其是对于初学者来说。\n\n*   可能存在错误的选项组合：由于函数选项模式允许在函数调用中指定多个选项，某些选项之间可能存在冲突或不兼容的情况。这可能导致意外的行为或错误的结果。\n\n*   不适用于所有情况：函数选项模式适用于有大量可选参数或者可配置选项的函数，但对于只有几个简单参数的函数，使用该模式可能过于复杂和冗余。在这种情况下，简单的命名参数可能更直观和易于使用。\n\n# 小结\n\n本文对 `Go` 函数选项模式（`Functional Options Pattern`）进行了详细介绍，并通过封装一个消息结构体的例子，展示了如何使用函数选项模式进行代码实现。\n\n在合适的情况下，我们可以使用函数选项模式来封装一些功能，定制函数的行为，提高代码的可读性和可扩展性。\n\n你是否在实际开发中使用过函数选项模式？欢迎评论区或关注公众号进群留言探讨。",
        "tags": [
            "Go"
        ],
        "publishDate": "2023-06-08 09:59:16",
        "category": "后端",
        "columns": [
            "Go 实用技巧"
        ],
        "cover": "https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/img/%E5%B0%81%E9%9D%A2-Go%20%E5%87%BD%E6%95%B0%E9%80%89%E9%A1%B9%E6%A8%A1%E5%BC%8F%EF%BC%88Functional%20Options%20Pattern%EF%BC%89.jpg",
        "likesCount": 4,
        "viewsCount": 414,
        "commentInfo": [
            {
                "id": "f950477f-b2c9-4969-b2f5-a96a6fb901cc",
                "content": "手机浏览器查看时，首页样式有错乱，后面可以考虑兼容下移动端",
                "username": "张~",
                "commentTime": "2023-06-14 19:55:03",
                "replies": [
                    {
                        "id": "9355309c-8f51-4a21-bfbc-b7cdf8674dc9",
                        "commentId": "f950477f-b2c9-4969-b2f5-a96a6fb901cc",
                        "content": "好的，感谢您的建议！",
                        "username": "陈明勇",
                        "replyToId": "",
                        "replyTo": "",
                        "replyTime": "2023-06-25 00:15:31"
                    }
                ]
            },
            {
                "id": "08eb2603-1c4f-464c-b1cd-7a947ade583c",
                "content": "占个位置，相信会火",
                "username": "~",
                "commentTime": "2023-06-11 14:33:25"
            }
        ],
        "isLike": false
    }
}

const handleCopyCodeSuccess = () => {
    ElMessage({
        message: 'copy success',
        type: 'success',
        customClass: 'text-16'
    });
}
definePageMeta({
    layout: "home"
})
</script>

<style scoped></style>