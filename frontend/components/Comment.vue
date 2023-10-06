<template>
    <div
        class="bg-#fff rounded-10 mt30 pb30 dark_bg_gray dark_text_white cursor-pointer ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--5">
        <div class="text-16 pt20 px20">
            最新评论
        </div>
        <el-divider class="important:mt24" />
        <div class="text-14 px20 " v-for="item in dataList">
            <el-text truncated class="dark_text_white">
                {{ item.name }}：{{ item.content }}
            </el-text>
            <nuxt-link :to='"/post/" + item.post_id' class="hover:bg-green-200 block mt8">
                <el-space class="c-#000/40 dark_text_white">
                    <div class="i-grommet-icons:article"></div>
                    <div>{{ item.post_title }}</div>
                </el-space>
            </nuxt-link>
            <el-divider class="important:my10" />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { getLatestComments, ILatestComment, } from "../api/comment"
import { IResponse, IListData } from "../api/http";

const dataList = ref([] as ILatestComment[]);
const commentsInfos = async () => {
    try {
        let postRes: any = await getLatestComments()
        let res: IResponse<IListData<ILatestComment>> = postRes.data.value
        dataList.value = res.data?.list || []
    } catch (error) {
        console.log(error);
    }
};
commentsInfos() 
</script>

<style scoped></style>