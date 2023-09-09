<template>
    <div>
        <span ref="text" class="msg text-30 text-shadow-xl"></span>
    </div>
</template>
    
<script lang="ts" setup>
import TypeIt from 'typeit'
const props = defineProps({ strings: { type: Array } })
const text = ref(null)
onMounted(() => {
    new (TypeIt)(text.value, {
        strings: '在学习的过程中赚钱，在赚钱的过程中学习。',
        cursorChar: "<span class='text-30'>|<span>",//用于光标的字符。HTML也可以
        speed: 150,
        lifeLike: true,// 使打字速度不规则
        cursor: true,//在字符串末尾显示闪烁的光标
        breakLines: false,// 控制是将多个字符串打印在彼此之上，还是删除这些字符串并相互替换
        loop: false,//是否循环
    }).go()
})
</script>
    
<style scoped></style>
