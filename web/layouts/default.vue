<template>
  <div>
    <div ref="myDom">
      <Header/>
    </div>
    <!-- <el-scrollbar max-height="100vh" @scroll="headerScroll" /> -->
    <div class="bg-#F0F2F5 dark_bg_black">
      <slot></slot>
    </div>
  </div>
</template>