<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>
<template><div class="client-shell"><header class="nav"><RouterLink to="/home" class="logo">SHORT<span>DRAMA</span></RouterLink><nav><RouterLink to="/home">首页</RouterLink><RouterLink to="/discover">发现</RouterLink><RouterLink to="/orders">我的订单</RouterLink><RouterLink to="/profile">我的</RouterLink></nav><div class="locale">US / English</div></header><main class="content"><RouterView/></main></div></template>