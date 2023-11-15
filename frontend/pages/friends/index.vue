<template>
    <div class="pt100 pb50 mx20">
        <div class="max-w-1280 mx-auto  rounded-23 bg-#fff p30" v-if="dataList.length > 0">
            <div class="text-26 font-600 ">
                友链
            </div>
            <el-row :gutter="20" class="lt-lg:important:display-block">
                <el-col :span="8" v-for="item in dataList" :key="item.name" class="lt-lg:important:max-w-100%">
                    <a :href="item.url" target="blank"
                        class="block  b-2 b-blue b-solid p20 h100 rounded-20 text-16 mt30 cursor-pointer hover:b-lightblue  hover:bg-#e5e5e5/30 active:bg-#e5e5e5 group">
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
                留言交友
            </div>
            <el-form ref="ruleFormRef" :model="req" class="mt30" :rules="rules">
                <el-row :gutter="50">
                    <el-col :span="12">
                        <el-form-item label="昵称" prop="name">
                            <el-input v-model="req.name" placeholder="请输入昵称" clearable />
                        </el-form-item>
                    </el-col>
                    <el-col :span="12">
                        <el-form-item label="邮箱" prop="email">
                            <el-input v-model="req.email" placeholder="请输入邮箱地址(用于接收通知)" clearable />
                        </el-form-item>
                    </el-col>
                </el-row>
                <el-form-item label="头像链接" prop="logo">
                    <el-input v-model="req.logo" placeholder="请输入头像链接" clearable />
                </el-form-item>
                <el-form-item label="网站链接" prop="url">
                    <el-input v-model="req.url" placeholder="请输入网站链接" clearable />
                </el-form-item>
                <el-form-item prop="description">
                    <el-input show-word-limit :maxlength="200" clearable type="textarea" v-model="req.description"
                        placeholder="请输入网站介绍" />
                </el-form-item>
                <el-form-item>
                    <div class="w-full text-center">
                        <el-button type="primary" @click="onSubmit(ruleFormRef)">提交申请</el-button>
                        <el-button @click="resetForm(ruleFormRef)">清空</el-button>
                    </div>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>


<script lang="ts" setup>
import type { FormInstance, FormRules } from 'element-plus'
import { getFriends, IFriend, applyForFriend, FriendReq } from "~/api/friend"
import { IResponse, IListData } from "~/api/http";
import { ElMessage } from 'element-plus'

definePageMeta({
    layout: "home"
})

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
        { max: 200, message: 'Length should be less than or equal to 200', trigger: 'blur' },
    ],
    email: [
        { pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$" }
    ]
})
const applyforFriendFunc = async (fReq: FriendReq, formEl: FormInstance | undefined) => {
    try {
        let cmtsRes: any = await applyForFriend(fReq)
        if (cmtsRes.error.value != null) {
            if (cmtsRes.error.value.statusCode == 403) {
                ElMessage.warning('友链模块未开放！')
            } else if (cmtsRes.error.value.statusCode == 429) {
                ElMessage.warning('请勿重复提交！')
            }
            return
        }
        let res: IResponse<null> = cmtsRes.data.value
        if ((res == undefined) || (res.code != 200)) {
            console.log("apply friend failed.")
            ElMessage.error('提交失败，如有问题，请联系站长。')
            return
        }
        ElMessage.success('提交成功，站长审核之后，您将会收到一封审核结果的邮件。')
        formEl?.resetFields()
    } catch (error) {
        console.log(error);
        ElMessage.error('发生未知错误，请重新提交。')
    }
};

const onSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid, fields) => {
        if (valid) {
            applyforFriendFunc(req, formEl)
        }
    })
}


const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}

const dataList = ref([] as IFriend[]);
const friends = async () => {
    try {
        let postRes: any = await getFriends()
        let res: IResponse<IListData<IFriend>> = postRes.data.value
        dataList.value = res.data?.list || []
    } catch (error) {
        console.log(error);
    }
};
friends()
</script>

<style scoped></style>