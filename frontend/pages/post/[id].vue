<template>
    <div class="bg-#F0F2F5 dark_bg_black pt100 pb30 px25">
        <el-row :gutter="40">

            <el-col :span="2">
                <el-affix :offset="350">
                    <div class=" text-center">
                        <el-space direction="vertical" :size="25">
                            <!-- 点赞 -->
                            <div
                                class="bg-#fff drop-shadow-lg mx-auto cursor-pointer text-22 w45 h45 flex justify-center items-center rounded-50% hover:bg-#e5e5e5 active:bg-#999 active:c-#fff">
                                <div class="i-grommet-icons:like"></div>
                            </div>
                            <div
                                class="bg-#fff drop-shadow-lg mx-auto cursor-pointer text-22 w45 h45 flex justify-center items-center rounded-50% hover:bg-#e5e5e5 active:bg-#999 active:c-#fff">
                                <div class="i-grommet-icons:share-rounded"></div>
                            </div>
                            <div
                                class="bg-#fff drop-shadow-lg mx-auto cursor-pointer text-22 w45 h45 flex justify-center items-center rounded-50% hover:bg-#e5e5e5 active:bg-#999 active:c-#fff">
                                <div class="i-grommet-icons:tooltip"></div>
                            </div>
                        </el-space>
                    </div>
                </el-affix>
            </el-col>


            <el-col :span="16">
                <div class=" rounded-15 bg-#fff p30">
                    <!-- 文章标题 -->
                    <div class="text-center text-42 font-500">{{ data?.title }}</div>
                    <!-- 作者信息 -->
                    <div class="text-center text-16 c-#999/80 my20 ">
                        <el-space :size="4">
                            <div>{{ data?.author }}</div>
                            <div>
                                {{ dayjs(data?.create_time).format('YYYY-MM-DD HH:mm:ss') }}
                            </div>
                            <div>·</div>
                            <div>阅读 </div>
                            <div>{{ data?.visit_count }}</div>
                            <div>·</div>
                            <div>全文字数 </div>
                            <div>{{ data?.word_count }}</div>
                        </el-space>
                    </div>
                    <v-md-preview :text="data?.content" @copy-code-success="handleCopyCodeSuccess"></v-md-preview>
                    <!-- 点赞 -->
                    <!-- <div
                        class=" mx-auto cursor-pointer text-22 w45 h45 flex justify-center items-center rounded-50% hover:bg-#e5e5e5 active:bg-#999 active:c-#fff">
                        <div class="i-grommet-icons:like"></div>
                    </div> -->
                </div>
                <div class="rounded-15 bg-#fff pt20 pb80 px30 mt25">
                    <!-- 发布评论 -->
                    <SubmitComments />
                    <div class="text-16">{{ 0 + ' 评论' }}</div>
                    <PostCommentsFather v-for="item, index in commentList" :key="index" class="" :commentInfo="item" />
                </div>
            </el-col>

            <el-col :span="6">
                <el-affix :offset="90">
                    <div class="w-full">
                        <Anchor :text="data?.content"></Anchor>
                    </div>
                </el-affix>
            </el-col>
        </el-row>
    </div>
</template>

<script lang="ts" setup>
import { ElMessage, dayjs } from 'element-plus'
import { getPostsById, IPostDetail } from "~/api/post"
import { IComment, getComments } from "~/api/comment"
import { IResponse, IListData } from "~/api/http";

// const router = useRouter()
const route = useRoute()
const sug: string = route.params.id as string

const data = ref<IPostDetail>();
const detail = async () => {
    try {
        let postRes: any = await getPostsById(sug)
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