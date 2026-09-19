<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { request,userId } from '../../api/client'
import { ElMessage } from 'element-plus'
const route=useRoute(), router=useRouter(), episode=ref<any>(null), episodes=ref<any[]>([]), video=ref<HTMLVideoElement|null>(null)
let hls:any=null
async function setupVideo(){await nextTick();if(!episode.value?.video_url||!video.value)return;const url=episode.value.video_url;if(video.value.canPlayType('application/vnd.apple.mpegurl')){video.value.src=url;return}try{const mod=await import('hls.js');const Hls=mod.default;if(Hls.isSupported()){hls=new Hls();hls.loadSource(url);hls.attachMedia(video.value);return}}catch{}video.value.src=url}
onMounted(async()=>{try{episodes.value=(await request('/api/v1/dramas/'+route.params.dramaId+'/episodes?user_id='+userId()')).items||[];episode.value=episodes.value.find(x=>String(x.id)===String(route.params.episodeId));if(!episode.value){ElMessage.error('剧集不存在或未解锁');return}if(!episode.value.unlocked&&episode.value.is_paid){ElMessage.info('请先购买本剧');router.replace('/checkout/'+route.params.dramaId);return}await setupVideo()}catch{ElMessage.error('剧集加载失败')}})
onBeforeUnmount(()=>{try{hls?.destroy()}catch{}})
</script>
<template>
<section class="watch-page"><el-button text @click="router.push('/drama/'+route.params.dramaId)">← 返回详情</el-button>
<div class="video-stage"><video ref="video" v-if="episode?.video_url" controls autoplay playsinline :poster="episode.poster_url||undefined"></video><el-empty v-else description="视频地址未配置"/></div>
<div class="watch-info"><h1>{{episode?.title||'正在加载...'}}</h1><p>Episode {{episode?.episode_no}}</p></div>
<el-card><template #header>剧集</template><div class="episode-list"><el-button v-for="e in episodes" :key="e.id" :type="String(e.id)===String(route.params.episodeId)?'primary':'default'" :disabled="e.is_paid&&!e.unlocked" @click="e.is_paid&&!e.unlocked?ElMessage.info('请先购买本剧'):router.push('/watch/'+route.params.dramaId+'/'+e.id)">{{e.episode_no}}. {{e.title}} <span v-if="e.is_paid&&!e.unlocked">🔒</span></el-button></div></el-card></section>
</template>