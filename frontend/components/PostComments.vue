<!--
 *                        _oo0oo_
 *                       o8888888o
 *                       88" . "88
 *                       (| -_- |)
 *                       0\  =  /0
 *                     ___/`---'\___
 *                   .' \\|     |// '.
 *                  / \\|||  :  |||// \
 *                 / _||||| -:- |||||- \
 *                |   | \\\  - /// |   |
 *                | \_|  ''\---/''  |_/ |
 *                \  .-\__  '-'  ___/-. /
 *              ___'. .'  /--.--\  `. .'___
 *           ."" '<  `.___\_<|>_/___.' >' "".
 *          | | :  `- \`.;`\ _ /`;.`/ - ` : | |
 *          \  \ `_.   \_ __\ /__ _/   .-` /  /
 *      =====`-.____`.___ \_____/___.-`___.-'=====
 *                        `=---='
 * 
 * 
 *      ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 * 
 *            佛祖保佑     永不宕机     永无BUG
 * 
 *        佛曰:  
 *                写字楼里写字间，写字间里程序员；  
 *                程序人员写程序，又拿程序换酒钱。  
 *                酒醒只在网上坐，酒醉还来网下眠；  
 *                酒醉酒醒日复日，网上网下年复年。  
 *                但愿老死电脑间，不愿鞠躬老板前；  
 *                奔驰宝马贵者趣，公交自行程序员。  
 *                别人笑我忒疯癫，我笑自己命太贱；  
 *                不见满街漂亮妹，哪个归得程序员？
 -->
<template>
    <div class="rounded-15 bg-#fff pt20 pb80 px30">
        <!-- 发布评论 -->
        <div>
            <el-form ref="ruleFormRef" :model="form" class="mt30">
                <el-row :gutter="50">
                    <el-col :span="12">
                        <el-form-item label="昵称" prop="nickName">
                            <el-input v-model="form.nickName" placeholder="请输入昵称" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="邮箱" prop="email">
                            <el-input v-model="form.email" placeholder="请输入邮箱地址(用于接收通知)" clearable />
                        </el-form-item>
                    </el-col>
                </el-row>

                <el-form-item prop="introduce">
                    <el-input show-word-limit :maxlength="30" clearable type="textarea" v-model="form.content"
                        placeholder="请输入评论内容，支持markdown格式" />
                </el-form-item>
                <el-form-item>
                    <div class="w-full text-center">
                        <el-button type="primary" @click="onSubmit">提交评论</el-button>
                        <el-button @click="resetForm(ruleFormRef)">清空</el-button>
                    </div>
                </el-form-item>
            </el-form>
        </div>
        <!-- 展示评论 -->
        <div>

            <div class="text-16">{{ commentInfo.length + ' 评论' }}</div>
            <div v-for="item, index in commentInfo" :key="index">
                <!-- 一级评论 -->
                <el-divider class="important:my24"></el-divider>
                <div>
                    <el-space :size="8" alignment="flex-start">
                        <el-avatar :size="36" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                        <div>
                            <el-space :size="0">
                                <div class="ml8">
                                    {{ item.username }}
                                </div>
                                <div class="ml8">
                                    发表于
                                </div>
                                <div class="ml4">
                                    {{ item.commentTime }}
                                </div>
                            </el-space>
                        </div>
                    </el-space>
                    <v-md-preview :text="item.content" class="px18"></v-md-preview>
                    <el-button class="ml50">回复</el-button>

                    <!-- 二级评论 -->
                    <div v-if="item.replies" class="mt25 px50">
                        <el-space direction="vertical" :fill="true" class="w-full" :spacer="spacer">

                            <div v-for="replies in item.replies" :key="replies.id"
                                class="b-1 b-#DDDBDB b-solid p20 rounded-10">
                                <el-space :size="8" alignment="flex-start">
                                    <el-avatar :size="36"
                                        src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                                    <div>
                                        <el-space :size="0">
                                            <div class="ml8 flex">
                                                <div>{{ replies.username }}</div>
                                                <div v-if="replies.replyTo" class="ml8">
                                                    <el-space>
                                                        <div>回复</div>
                                                        <div>{{ replies.replyTo }}</div>
                                                    </el-space>
                                                </div>
                                            </div>
                                            <div class="ml8">
                                                发表于
                                            </div>
                                            <div class="ml4">
                                                {{ replies.replyTime }}
                                            </div>
                                        </el-space>
                                    </div>
                                </el-space>
                                <v-md-preview :text="replies.content" class="px18"></v-md-preview>
                                <el-button class="ml50">回复</el-button>
                            </div>
                        </el-space>
                    </div>

                </div>
            </div>
        </div>
    </div>
</template>


<script lang="ts" setup>
import type { FormInstance, FormRules } from 'element-plus'
import { ElDivider } from 'element-plus'
const spacer = h(ElDivider)
const ruleFormRef = ref<FormInstance>()
const form = ref({
    nickName: '',
    email: '',
    content: ''
})
const onSubmit = () => {
    console.log('submit!')
}
const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}
const commentInfo = [
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
]
</script>

<style scoped>
:deep(.el-divider--horizontal) {
    margin: 0
}
</style>