<template>
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

            <el-form-item prop="content">
                <el-input :autosize="{ minRows: 2, maxRows: 4 }" type="textarea" show-word-limit :maxlength="30" clearable
                    v-model="form.content" placeholder="请输入评论内容，支持markdown格式" v-if="!review" />
                <!-- 预览窗口 -->
                <div class="bg-#e5e5e5 w-full rounded-8" v-if="review">
                    <v-md-preview :text="form.content"></v-md-preview>
                </div>
            </el-form-item>
            <el-form-item>
                <div class="w-full text-center">
                    <el-button type="primary" @click="onSubmit">提交评论</el-button>
                    <el-button @click="resetForm(ruleFormRef)">清空</el-button>
                    <el-button v-if="!review" @click="review = true">预览</el-button>
                    <el-button v-if="review" @click="review = false">编辑</el-button>
                </div>
            </el-form-item>
        </el-form>
    </div>
</template>


<script lang="ts" setup>
import type { FormInstance } from 'element-plus'
const ruleFormRef = ref<FormInstance>()
const review = ref<boolean>(false)
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
</script>

<style scoped></style>