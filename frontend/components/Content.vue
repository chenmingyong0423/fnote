<template>
    <div class="group bg-#fff mb25 p25 rounded-10 dark_bg_gray dark_text_white cursor-pointer ">
        <el-space :size="30">
            <div class="w300 h150 overflow-hidden rounded-10">
                <el-image class=" w-full h-full  group-hover:w105% group-hover:h105% ease-in duration-200" position="top"
                    fit="cover" src='https://fuss10.elemecdn.com/e/5d/4a731a90594a4af544c0c25941171jpeg.jpeg' />
            </div>

            <div class=" text-16 flex justify-between w-650 items-end">
                <div class="flex flex-col justify-between h150 ">
                    <div>
                        <div class="text-24 origin-center group-hover:translate-y--5 ease-in duration-200">
                            {{ props.postData.title }}
                        </div>
                        <div class="text-16">
                            {{ props.postData.summary }}
                        </div>
                    </div>
                    <div>
                        <el-space>
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
                <div>
                    2023-08-13 20:08:09
                </div>
            </div>

        </el-space>
    </div>
</template>


<script lang="ts" setup>
const props = defineProps(['postData'])


</script>

<style scoped></style>