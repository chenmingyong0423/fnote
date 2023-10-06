<template>
    <div>
        <div class="py15">
            <div class="text-30 mb20 ml10">标签:{{ menu?.name || '未知' }}</div>
        </div>
        <div class="bg-#e5e5e5/40 p20 rounded-10">
            <div v-for="item, index in dataList" :key="index">
                <Content @click="router.push(`/post/${item.sug}`)" :postData="item"></Content>
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
import { IMenu } from "~/api/category"

const route = useRoute()
const router = useRouter()
const homeStore = useHomeStore()

const menu: IMenu | undefined = homeStore.menuList.find((item: IMenu) => {
    if (item.route === route.path)
        return item
})

const dataList = ref<IPost[]>([])
const rq = ref<PageRequest>({
    pageNo: 1,
    pageSize: 5,
    sortField: '',
    sortOrder: '',
    search: '',
    category: '',
    tags: [menu!.name]
})

const postInfos = async () => {
    try {
        let postRes: any = await getPosts(rq.value)
        let res: IResponse<IPageData<IPost>> = postRes.data.value
        dataList.value = res.data?.list || []
    } catch (error) {
        console.log(error);
    }
};

postInfos()
</script>

<style scoped></style>
