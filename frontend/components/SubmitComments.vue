<template>
    <div>
        <el-form ref="ruleFormRef" :model="req" class="mt30" :rules="rules">
            <el-row :gutter="50">
                <el-col :span="12">
                    <el-form-item label="昵称" prop="nickName">
                        <el-input v-model="req.name" placeholder="请输入昵称" clearable />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="邮箱" prop="email">
                        <el-input v-model="req.email" placeholder="请输入邮箱地址(用于接收通知)" clearable />
                    </el-form-item>
                </el-col>
            </el-row>

            <el-form-item prop="content">
                <el-input show-word-limit :maxlength="30" clearable type="textarea" v-model="req.content"
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
</template>


<script lang="ts" setup>
import type { FormInstance, FormRules } from 'element-plus'
import { applyForFriend, FriendReq } from '~/api/friend'
import { IResponse } from "~/api/http";

const ruleFormRef = ref<FormInstance>()
const req = ref<FriendReq>({})

const rules = reactive<FormRules<FriendReq>>({
    name: [
        { required: true, message: 'Please input name', trigger: 'blur' },
    ],
    url: [
        { required: true, message: 'Please input url', trigger: 'blur' },
    ],
    logo: [
        { required: true, message: 'Please input logo', trigger: 'blur' },
    ],
    description: [
        { required: true, message: 'Please input description', trigger: 'blur' },
        { max: 20, message: 'Length should be less than or equal to 20', trigger: 'blur' },
    ],
    email: [
        { pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$" }
    ]
})
const applyforFriendFunc = async () => {
    try {
        let cmtsRes: any = await applyForFriend(req)
        let res: IResponse<null> = cmtsRes.data.value
        if ((res == undefined) || (res.code != 0)) {
            console.log("apply friend failed.")
        }
    } catch (error) {
        console.log(error);
    }
};

const onSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid, fields) => {
        if (valid) {
            console.log('submit!')
        } else {
            console.log('error submit!', fields)
        }
    })
}
const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}
</script>

<style scoped></style>