<template>
    <div class="bg-#F0F2F5 dark_bg_black pt100 pb30 px25">
        <div class=" rounded-15 bg-#fff p30">
            <!-- 文章标题 -->
            <div class="text-center text-42 font-500">{{ data?.title }}</div>
            <!-- 作者信息 -->
            <div class="text-center text-16 c-#999/80 my20 lt-lg:text-14">
                <el-space :size="4">
                    <div>{{ data?.author }}</div>
                    <div>
                        {{ dayjs(data?.create_time).format('YYYY-MM-DD HH:mm:ss') }}
                    </div>
                    <div>·</div>
                    <div>阅读 </div>
                    <div>{{ data?.visit_count }}</div>
                    <div>·</div>
                    <div class="lt-lg-display-none">全文字数 </div>
                    <div class="lt-lg-display-none">{{ data?.word_count }}</div>
                </el-space>
            </div>
            <v-md-preview :text="data?.content" @copy-code-success="handleCopyCodeSuccess"></v-md-preview>
        </div>
        <div class="rounded-15 bg-#fff pt20 pb80 px30 mt25">
            <!-- 发布评论 -->
            <SubmitComments />
            <div class="text-16">{{ 0 + ' 评论' }}</div>
            <PostCommentsFather v-for="item, index in commentList" :key="index" class="" :commentInfo="item" />
        </div>
    </div>
</template>

s
<script lang="ts" setup>
import { ElMessage } from 'element-plus'
import { dayjs } from 'element-plus'
import type { } from 'element-plus'
import { getPostsById, IPostDetail } from "~/api/post"
import { IComment, getComments } from "~/api/comment"
import { IResponse, IListData } from "~/api/http";

// const router = useRouter()
const route = useRoute()

const data = ref<IPostDetail>();
const detail = async () => {
    try {
        let postRes: any = await getPostsById("about-me")
        let res: IResponse<IPostDetail> = postRes.data.value
        data.value = res.data
    } catch (error) {
        console.log(error);
    }
};
detail()

const commentList = ref<IComment[]>([])

const comments = async () => {
    try {
        let cmtsRes: any = await getComments(sug)
        let res: IResponse<IListData<IComment>> = cmtsRes.data.value
        commentList.value = res.data?.list || []
    } catch (error) {
        console.log(error);
    }
};

comments()

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