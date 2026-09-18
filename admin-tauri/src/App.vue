<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
const API='http://127.0.0.1:8081'
const loggedIn=ref(false), loading=ref(false), items=ref<any[]>([])
const loginForm=ref({username:'admin',password:''})
const form=ref({title:'',description:'',cover:'',country:'US',language:'en',total_episodes:1,is_paid:false})
async function load(){loading.value=true;try{const r=await fetch(API+'/api/v1/admin/dramas?page=1&page_size=100');if(!r.ok)throw Error('加载失败');const j=await r.json();items.value=j.items||[]}catch(e:any){ElMessage.error(e.message)}finally{loading.value=false}}
async function login(){if(!loginForm.value.username||!loginForm.value.password){ElMessage.warning('请输入账号和密码');return}loggedIn.value=true;await load();ElMessage.success('进入后台')}
async function create(){if(!form.value.title.trim()){ElMessage.warning('请输入短剧标题');return}const r=await fetch(API+'/api/v1/admin/dramas',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(form.value)});if(!r.ok){ElMessage.error('创建失败');return}ElMessage.success('创建成功');form.value={...form.value,title:'',description:'',cover:''};await load()}
async function toggle(x:any){const status=x.status===1?2:1;const r=await fetch(API+'/api/v1/admin/dramas/'+x.id+'/status',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({status})});if(r.ok){ElMessage.success(status===1?'已发布':'已下架');await load()}else ElMessage.error('状态更新失败')}
async function edit(x:any){const {value}=await ElMessageBox.prompt('请输入新的标题','编辑短剧',{inputValue:x.title,confirmButtonText:'保存'});if(value){await fetch(API+'/api/v1/admin/dramas/'+x.id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({...x,title:value})});await load()}}
onMounted(()=>{if(loggedIn.value)load()})
</script>
<template>
<el-container class="shell" v-if="loggedIn">
 <el-header class="top"><div><b>Short Drama Admin</b><span>海外短剧运营后台</span></div><el-button text @click="loggedIn=false">退出</el-button></el-header>
 <el-container>
  <el-aside width="220px"><el-menu default-active="1"><el-menu-item index="1">短剧管理</el-menu-item><el-menu-item index="2">短剧发布</el-menu-item><el-menu-item index="3">订单 / 付费</el-menu-item></el-menu></el-aside>
  <el-main>
   <el-card class="publish"><template #header><b>发布短剧</b></template>
    <el-form label-width="90px"><el-form-item label="标题"><el-input v-model="form.title"/></el-form-item><el-form-item label="简介"><el-input v-model="form.description" type="textarea"/></el-form-item><el-form-item label="封面"><el-input v-model="form.cover" placeholder="CDN 图片 URL"/></el-form-item>
    <el-form-item label="国家"><el-input v-model="form.country"/></el-form-item><el-form-item label="语言"><el-input v-model="form.language"/></el-form-item><el-form-item label="总集数"><el-input-number v-model="form.total_episodes" :min="1"/></el-form-item><el-form-item label="付费"><el-switch v-model="form.is_paid"/></el-form-item><el-button type="primary" @click="create">创建草稿</el-button></el-form>
   </el-card>
   <el-card><template #header><div class="table-head"><b>短剧管理</b><el-button @click="load" :loading="loading">刷新</el-button></div></template>
    <el-table :data="items" stripe><el-table-column prop="id" label="ID" width="80"/><el-table-column prop="title" label="标题" min-width="220"/><el-table-column prop="country" label="国家" width="100"/><el-table-column prop="language" label="语言" width="100"/><el-table-column label="类型" width="100"><template #default="{row}"><el-tag :type="row.is_paid?'warning':'success'">{{row.is_paid?'付费':'免费'}}</el-tag></template></el-table-column><el-table-column label="状态" width="100"><template #default="{row}"><el-tag>{{row.status===1?'已发布':row.status===2?'已下架':'草稿'}}</el-tag></template></el-table-column><el-table-column label="操作" width="190"><template #default="{row}"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="primary" @click="toggle(row)">{{row.status===1?'下架':'发布'}}</el-button></template></el-table-column></el-table>
   </el-card>
  </el-main>
 </el-container>
</el-container>
<div v-else class="login-page"><el-card class="login-card"><h1>Short Drama</h1><p>运营管理后台</p><el-form @submit.prevent="login"><el-form-item><el-input v-model="loginForm.username" placeholder="管理员账号"/></el-form-item><el-form-item><el-input v-model="loginForm.password" type="password" show-password placeholder="管理员密码"/></el-form-item><el-button type="primary" class="full" @click="login">登录</el-button></el-form></el-card></div>
</template>