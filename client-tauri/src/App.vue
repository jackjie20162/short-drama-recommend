<script setup lang="ts">
import{computed}from'vue';import{RouterLink,RouterView}from'vue-router';import{logout}from'./api/client'
const loggedIn=computed(()=>!!localStorage.getItem('auth_token'));const email=computed(()=>{try{return JSON.parse(localStorage.getItem('auth_user')||'{}').email||'Guest'}catch{return'Guest'}})
function signout(){logout();location.reload()}
</script>
<template><div class="client-shell"><header class="nav"><RouterLink to="/home" class="logo">SHORT<span>DRAMA</span></RouterLink><nav><RouterLink to="/home">首页</RouterLink><RouterLink to="/discover">发现</RouterLink><RouterLink to="/orders">我的订单</RouterLink><RouterLink to="/profile">我的</RouterLink></nav><div class="locale">{{email}} · <RouterLink v-if="!loggedIn" to="/login">Login</RouterLink><a v-else href="#" @click.prevent="signout">Logout</a></div></header><main class="content"><RouterView/></main></div></template>