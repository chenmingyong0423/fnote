<template>
    <!-- bg-#F0F2F5 -->
    <div>
        <!-- 后期加  首页背景图 -->
        <!-- <div class="text-0 relative">
            <div class="text-center w-full font-bold text-70 absolute c-#fff top-30% text-shadow-xl">
                xxx's Blog
                <TypeIt></TypeIt>
            </div>
            <img src="../assets/images/bg.png" class="w-full h-100vh object-cover" />
        </div> -->
        <div class="pt70 pb40 px25">
            <el-row :gutter="20" class="lt-lg:important:display-block">
                <el-col :span="17" class="lt-lg:important:max-w-100%">
                    <Content v-for="item in dataList" :postData="item" :key="item.sug">
                    </Content>
                </el-col>
                <el-col :span="7" class="lt-lg:important:display-none">
                    <div class="w-full">
                        <Profile />
                        <el-affix :offset="0">
                            <Comment />
                        </el-affix>
                    </div>
                </el-col>
            </el-row>
        </div>
        <!-- <SearchDialog /> -->
    </div>
</template>

<script lang="ts" setup>
import { getLatestPosts, IPost } from "~/api/post"
import { IResponse, IListData } from "~/api/http";
import 'element-plus/theme-chalk/dark/css-vars.css'
const dataList = ref<IPost[]>([]);
const postInfos = async () => {
    try {
        let postRes: any = await getLatestPosts()
        let res: IResponse<IListData<IPost>> = postRes.data.value
        dataList.value = res.data?.list || []
    } catch (error) {
        console.log(error);
    }
};
postInfos()

definePageMeta({
    layout: "home"
})
</script>

<style scoped></style>