<template>
    <div
        class="group relative mb25 rounded-10 overflow-hidden  ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--5">
        <div class=" bg-#fff  p25  dark_bg_gray dark_text_white">

            <el-space :size="30">
                <div class="w300 h150 overflow-hidden rounded-10">
                    <el-image @click="router.push(`post/${postData.sug}`)" class="cursor-pointer w-full h-full"
                        position="top" fit="cover"
                        src='https://fuss10.elemecdn.com/e/5d/4a731a90594a4af544c0c25941171jpeg.jpeg' />
                </div>

                <div class=" text-16 flex justify-between items-end w700">
                    <div class="flex flex-col justify-between h150 ">
                        <div>
                            <div @click="router.push(`post/${postData.sug}`)"
                                class="cursor-pointer text-24 c_title_blue group-hover:text-red ease-in duration-200">
                                {{ postData.title }}
                            </div>
                            <div class="text-16 ">
                                {{ postData.summary }}
                            </div>
                        </div>
                        <div>
                            <el-space :size="20">
                                <div>
                                    <el-space :size="5">
                                        <div class="i-grommet-icons:view" />
                                        <div>13</div>
                                    </el-space>
                                </div>
                                <div>
                                    <el-space :size="5">
                                        <div class="i-grommet-icons:like" />
                                        <div>55</div>
                                    </el-space>
                                </div>
                                <div>
                                    <el-space :size="5">
                                        <div class="i-grommet-icons:chat-option" />
                                        <div>66</div>
                                    </el-space>
                                </div>
                            </el-space>
                        </div>
                    </div>
                    <div class="text-center w-full">
                        {{ dayjs(postData.createTime).format('YYYY-MM-DD HH:mm:ss') }}
                    </div>
                </div>

            </el-space>
        </div>
        <div class="mx-auto h3 bg-blue w0 group-hover:w-full ease-in duration-300" />
        <div
            class="absolute top-15 right--100 p5 rounded-3 c_text_white bg-blue text-13 group-hover:right-10 ease-linear duration-150">
            {{ postData.tags[0] }}
        </div>
    </div>
</template>


<script lang="ts" setup>
import { dayjs } from 'element-plus'
const props = defineProps(['postData'])
const router = useRouter()

</script>

<style scoped></style>