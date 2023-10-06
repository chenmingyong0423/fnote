<template>
    <div>
        <div class="py15">
            <div class="text-30 mb20 ml10">分类</div>
            <my-tag v-for="tag in homeStore.menuList" :key="tag.name" @click="router.push(tag.route)">
                {{ tag.name }}
            </my-tag>
        </div>
        <div class="py15">
            <div class="text-30 mb20 ml10">标签</div>
            <my-tag v-for="tag in homeStore.menuList" :key="tag.name" @click="router.push(tag.route)">
                {{ tag.name }}
            </my-tag>
        </div>
        <div class="mt40 bg-#e5e5e5/40 p20 rounded-10">
            <div class="text-30 mb20 ml10">我的文章</div>
            <div v-for="item in dataList ">
                <Content @click="router.push(`post/${item.sug}`)" :postData="item"></Content>
                <!-- <el-divider /> -->
            </div>
        </div>
        <Pagination />
    </div>
</template>

<script lang="ts" setup>
import { useHomeStore } from '~/store/home';
import { IPost, getPosts, PageRequest } from '~/api/post';
import { IResponse, IPageData } from "~/api/http";

const homeStore = useHomeStore()
const router = useRouter()
const dataList = ref([] as IPost[])
const rq = ref({
    pageNo: 1,
    pageSize: 5,
    sortField: '',
    sortOrder: '',
    search: '',
    category: '',
    tags: []
} as PageRequest)
const postInfos = async () => {
    try {
        let postRes: any = await getPosts(rq)
        let res: IResponse<IPageData<IPost>> = postRes.data.value
        dataList.value = res.data?.list || []
    } catch (error) {
        console.log(error);
    }
};
postInfos()


</script>

<style scoped>
.el-pagination {
    justify-content: center;
}
</style>