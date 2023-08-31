<template>
    <div class="pt100 pb50">
        <div class="max-w-1280 mx-auto  rounded-23 bg-#fff p30">
            <div class="text-26 font-600 ">
                友链
            </div>
            <el-row :gutter="20">
                <el-col :span="8" v-for=" in 7">
                    <div class="b-2 b-blue b-solid p20 h100 rounded-20 text-16 mt30">
                        <el-space alignment="flex-start" :size="0">
                            <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" :size="50"
                                class="mr15" />
                            <div>
                                <el-space direction="vertical" alignment="flex-start">
                                    <div>工作笔记</div>
                                    <div class="text-#000/50">
                                        集结工作前端和全栈技术、行业视野、技巧总结，个人生活与摄影，养成记录和分享的习惯。
                                    </div>
                                </el-space>

                            </div>
                        </el-space>
                    </div>
                </el-col>
            </el-row>
        </div>
        <div class="max-w-1280 mx-auto  rounded-23 bg-#fff p30 mt30">
            <div class="text-26 font-600 ">
                留言交友 (至少填昵称和网址)
            </div>
            <el-form :model="form" class="mt30">
                <el-row :gutter="50">
                    <el-col :span="12">
                        <el-form-item label="昵称">
                            <el-input v-model="form.nickName" placeholder="请输入昵称" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="邮箱">
                            <el-input v-model="form.email" placeholder="请输入邮箱地址(用于接收通知)" clearable />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-form-item label="头像链接">
                    <el-input v-model="form.avatarUrl" placeholder="请输入头像链接" clearable />
                </el-form-item>
                <el-form-item label="网站链接">
                    <el-input v-model="form.link" placeholder="请输入网站链接" clearable />
                </el-form-item>
                <el-form-item>
                    <el-input show-word-limit :maxlength="30" clearable type="textarea" v-model="form.introduce"
                        placeholder="请输入网站介绍" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="onSubmit">提交申请</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>


<script lang="ts" setup>
definePageMeta({
    layout: "home"
})
const form = ref({
    nickName: '',
    email: '',
    avatarUrl: '',
    link: '',
    introduce: ''
})

const onSubmit = () => {
    console.log('submit!')
}
</script>

<style scoped></style>