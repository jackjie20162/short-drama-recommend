<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { logout } from './api/client'
</script>
<template><div class="client-shell"><header class="nav"><RouterLink to="/home" class="logo">SHORT<span>DRAMA</span></RouterLink><nav><RouterLink to="/home">首页</RouterLink><RouterLink to="/discover">发现</RouterLink><RouterLink to="/orders">我的订单</RouterLink><RouterLink to="/profile">我的</RouterLink></nav><div class="locale">{{localStorage.getItem("auth_user") ? JSON.parse(localStorage.getItem("auth_user")||"{}").email : "Guest"}} · <RouterLink v-if="!localStorage.getItem("auth_token")" to="/login">Login</RouterLink><a v-else href="#" @click.prevent="logout();location.reload()">Logout</a></div></header><main class="content"><RouterView/></main></div></template>