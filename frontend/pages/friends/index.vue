<template>
    <div class="pt100 pb50">
        <div class="max-w-1280 mx-auto  rounded-23 bg-#fff p30">
            <div class="text-26 font-600 ">
                友链
            </div>
            <el-row :gutter="20">
                <el-col :span="8" v-for="item in testData.data.list" :key="item.name">
                    <a :href="item.url" target="blank"
                        class="block b-2 b-blue b-solid p20 h100 rounded-20 text-16 mt30 cursor-pointer hover:b-lightblue  hover:bg-#e5e5e5/30 active:bg-#e5e5e5 group">
                        <el-space alignment="flex-start" :size="0">
                            <el-avatar :src="item.logo" :size="50" class="mr15" />
                            <div>
                                <el-space direction="vertical" alignment="flex-start">
                                    <div class="text-#000 group-active:text-#fff">{{ item.name }}</div>
                                    <div class="text-#000/50 group-active:text-#fff">
                                        {{ item.description }}
                                    </div>
                                </el-space>

                            </div>
                        </el-space>
                    </a>
                </el-col>
            </el-row>
        </div>
        <div class="max-w-1280 mx-auto  rounded-23 bg-#fff p30 mt30">
            <div class="text-26 font-600 ">
                留言交友 (至少填昵称和网址)
            </div>
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
                <el-form-item label="头像链接" prop="avatarUrl">
                    <el-input v-model="form.avatarUrl" placeholder="请输入头像链接" clearable />
                </el-form-item>
                <el-form-item label="网站链接" prop="link">
                    <el-input v-model="form.link" placeholder="请输入网站链接" clearable />
                </el-form-item>
                <el-form-item prop="introduce">
                    <el-input show-word-limit :maxlength="30" clearable type="textarea" v-model="form.introduce"
                        placeholder="请输入网站介绍" />
                </el-form-item>
                <el-form-item>
                    <div class="w-full text-center">
                        <el-button type="primary" @click="onSubmit">提交申请</el-button>
                        <el-button @click="resetForm(ruleFormRef)">清空</el-button>
                    </div>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>


<script lang="ts" setup>
import type { FormInstance, FormRules } from 'element-plus'
definePageMeta({
    layout: "home"
})
const ruleFormRef = ref<FormInstance>()
const form = ref({
    nickName: '',
    email: '',
    avatarUrl: '',
    link: '',
    introduce: ''
})
const testData = {
    "code": 200,
    "message": "OK",
    "data": {
        "list": [
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            },
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            },
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            },
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            },
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            },
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            },
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            },
            {
                "name": "陈明勇",
                "url": "https:chenmingyong.cn",
                "logo": "https:chenmingyong.cn/logo.jpg",
                "description": "一个热爱技术，喜欢钻研技术的程序员。",
                "status": 2,
                "priority": 0
            }
        ]
    }
}
const onSubmit = () => {
    console.log('submit!')
}

const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}
</script>

<style scoped></style>