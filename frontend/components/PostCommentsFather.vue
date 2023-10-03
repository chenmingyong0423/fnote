<!--
 *                        _oo0oo_
 *                       o8888888o
 *                       88" . "88
 *                       (| -_- |)
 *                       0\  =  /0
 *                     ___/`---'\___
 *                   .' \\|     |// '.
 *                  / \\|||  :  |||// \
 *                 / _||||| -:- |||||- \
 *                |   | \\\  - /// |   |
 *                | \_|  ''\---/''  |_/ |
 *                \  .-\__  '-'  ___/-. /
 *              ___'. .'  /--.--\  `. .'___
 *           ."" '<  `.___\_<|>_/___.' >' "".
 *          | | :  `- \`.;`\ _ /`;.`/ - ` : | |
 *          \  \ `_.   \_ __\ /__ _/   .-` /  /
 *      =====`-.____`.___ \_____/___.-`___.-'=====
 *                        `=---='
 * 
 * 
 *      ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 * 
 *            佛祖保佑     永不宕机     永无BUG
 * 
 *        佛曰:  
 *                写字楼里写字间，写字间里程序员；  
 *                程序人员写程序，又拿程序换酒钱。  
 *                酒醒只在网上坐，酒醉还来网下眠；  
 *                酒醉酒醒日复日，网上网下年复年。  
 *                但愿老死电脑间，不愿鞠躬老板前；  
 *                奔驰宝马贵者趣，公交自行程序员。  
 *                别人笑我忒疯癫，我笑自己命太贱；  
 *                不见满街漂亮妹，哪个归得程序员？
 -->
<template>
    <div>

        <!-- 展示评论 -->
        <div>

            <div>
                <!-- 一级评论 -->
                <el-divider class="important:my24"></el-divider>
                <div>
                    <el-space :size="8" alignment="flex-start">
                        <el-avatar :size="36" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                        <div>
                            <el-space :size="0">
                                <div class="ml8">
                                    {{ commentInfo.username }}
                                </div>
                                <div class="ml8">
                                    发表于
                                </div>
                                <div class="ml4">
                                    {{ formatTimestamp(commentInfo.comment_time) }}
                                </div>
                            </el-space>
                        </div>
                    </el-space>
                    <v-md-preview :text="commentInfo.content" class="px18"></v-md-preview>
                    <el-button class="ml50" @click="showReply = true">回复</el-button>
                    <el-button class="ml50" v-if="showReply" @click="showReply = false">取消</el-button>
                    <SubmitComments v-if="showReply" />
                    <!-- 二级评论 -->
                    <PostCommentsSon :replies="commentInfo.replies" />
                </div>
            </div>
        </div>
    </div>
</template>


<script lang="ts" setup>
import { IComment } from "~/api/comment"
import { dayjs } from 'element-plus'

const props = defineProps(['commentInfo'])
const commentInfo: IComment = props.commentInfo
let showReply = ref(false)
const formatTimestamp = (timestamp: number): string => {
    return dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss');
};
</script>

<style scoped>
:deep(.el-divider--horizontal) {
    margin: 0
}
</style>