<script setup lang="ts">import {RouterLink,RouterView} from 'vue-router'</script>
<template>
<div class="admin-shell">
 <header class="topbar"><strong>SHORT DRAMA ADMIN</strong><nav><RouterLink to="/dashboard">概览</RouterLink><RouterLink to="/dramas">短剧</RouterLink><RouterLink to="/episodes">剧集</RouterLink><RouterLink to="/orders">订单</RouterLink><RouterLink to="/users">用户</RouterLink><RouterLink to="/settings">系统设置</RouterLink></nav></header>
 <main class="page"><RouterView/></main>
</div>
</template>
<style scoped>
.admin-shell{min-height:100vh;background:#f5f7fa}.topbar{height:64px;background:#111827;color:#fff;display:flex;align-items:center;padding:0 28px;gap:40px}.topbar nav{display:flex;gap:22px}.topbar a{color:#cbd5e1;text-decoration:none}.topbar a.router-link-active{color:#fff;font-weight:700}.page{padding:28px}
</style>
