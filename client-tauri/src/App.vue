<script setup lang="ts">
import {ref,onMounted} from 'vue'
import {ElMessage} from 'element-plus'
const API='http://127.0.0.1:8080', data=ref<any[]>([]), detail=ref<any>(null), loading=ref(false)
const episodes=ref<any[]>([])
async function load(){loading.value=true;try{const r=await fetch(API+'/api/v1/feed?user_id=1&country=US&language=en&page_size=20');const j=await r.json();data.value=j.items||[]}catch{ElMessage.error('无法连接服务')}finally{loading.value=false}}
async function openDrama(x:any){const r=await fetch(API+'/api/v1/dramas/'+x.drama_id);const j=await r.json();detail.value=j.drama;episodes.value=[]}
function watch(e:any){if(e.is_paid||detail.value.is_paid)ElMessage.info('该剧集需要购买后观看');else if(e.video_url)window.open(e.video_url,'_blank');else ElMessage.warning('该集暂无视频地址')}
onMounted(load)
</script>
<template>
<div class="app">
<header><div class="brand">Short Drama</div><div class="locale">US · English</div></header>
<main><div class="hero"><h1>Short Drama</h1><p>Discover stories made for you.</p></div><div class="toolbar"><b>For You</b><el-button text :loading="loading" @click="load">刷新</el-button></div>
<div class="feed"><el-card v-for="x in data" :key="x.drama_id" class="drama" shadow="hover" @click="openDrama(x)"><div class="cover"><img v-if="x.cover" :src="x.cover"><span v-else>DRAMA</span></div><h3>{{x.title}}</h3><small>{{x.reasons?.join(' · ')}}</small></el-card></div></main>
<el-dialog v-model="detail" width="760px" title="短剧详情"><div v-if="detail"><h2>{{detail.title}}</h2><p>{{detail.description}}</p><el-tag :type="detail.is_paid?'warning':'success'">{{detail.is_paid?'付费短剧':'免费观看'}}</el-tag><el-divider/><h3>剧集</h3><el-empty v-if="!episodes.length" description="剧集接口即将接入"/><el-button v-for="e in episodes" :key="e.id" @click="watch(e)">第{{e.episode_no}}集 {{e.title}}</el-button><div class="watch"><el-button type="primary" @click="ElMessage.info(detail.is_paid?'进入购买流程':'开始播放')">{{detail.is_paid?'购买观看':'开始观看'}}</el-button></div></div></el-dialog>
</div>
</template>