<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { loadStripe } from '@stripe/stripe-js'
import { ElMessage } from 'element-plus'

const API = 'http://127.0.0.1:8080'
const page = ref('home')
const loading = ref(false)
const dramas = ref<any[]>([])
const detail = ref<any>(null)
const episodes = ref<any[]>([])
const payment = ref<any>(null)
const provider = ref<'STRIPE'|'PAYPAL'>('STRIPE')
const userId = 1
const orders = ref<any[]>([])
const orderIdInput = ref('')
const cardNumber=ref(''), cardExpiry=ref(''), cardCvc=ref(''), stripeLoading=ref(false)

async function request(path:string, init?:RequestInit){
  const r = await fetch(API + path, init)
  const j = await r.json().catch(()=>({}))
  if(!r.ok) throw new Error(j.message || j.error || '请求失败')
  return j
}
async function loadHome(){
  loading.value=true
  try { const j=await request('/api/v1/feed?user_id=1&country=US&language=en&page_size=20'); dramas.value=j.items||[] }
  catch(e:any){ ElMessage.error(e.message) } finally { loading.value=false }
}
async function openDrama(x:any){
  try {
    const j=await request('/api/v1/dramas/'+x.drama_id)
    detail.value=j.drama
    const ep=await request('/api/v1/dramas/'+x.drama_id+'/episodes?user_id='+userId); episodes.value=ep.items||[]
    page.value='detail'
  } catch(e:any){ElMessage.error(e.message)}
}
async function buy(){
  if(!detail.value) return
  try {
    const j=await request('/api/v1/payments/orders',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({
      user_id:userId, drama_id:detail.value.id, provider:provider.value, currency:'USD',
      return_url:'http://127.0.0.1:1421/?payment=success', cancel_url:'http://127.0.0.1:1421/?payment=cancel'
    })})
    payment.value=j
    if(provider.value==='PAYPAL' && j.approve_url) window.open(j.approve_url,'_blank')
    ElMessage.success('订单已创建')
    page.value='checkout'
  } catch(e:any){ElMessage.error(e.message)}
}
function watchEpisode(e:any){
  if(!e.unlocked){ ElMessage.info('该集需要购买本剧后解锁'); return }
  if(e.video_url) window.open(e.video_url,'_blank'); else ElMessage.warning('暂无视频地址')
}
async function payStripe(){
 if(!payment.value?.client_secret){ElMessage.error('Stripe PaymentIntent 尚未创建');return}
 const key=(import.meta as any).env?.VITE_STRIPE_PUBLISHABLE_KEY
 if(!key){ElMessage.error('请配置 VITE_STRIPE_PUBLISHABLE_KEY');return}
 stripeLoading.value=true
 try{
  const stripe=await loadStripe(key);if(!stripe)throw Error('Stripe.js 加载失败')
  const result=await stripe.confirmCardPayment(payment.value.client_secret,{payment_method:{card:{number:cardNumber.value,exp_month:Number(cardExpiry.value.split('/')[0]),exp_year:Number(cardExpiry.value.split('/')[1]),cvc:cardCvc.value}}})
  if(result.error)throw Error(result.error.message||'Stripe 支付失败')
  ElMessage.success('Stripe 支付成功，等待订单确认')
  await queryOrderById(payment.value.order_id)
 }catch(e:any){ElMessage.error(e.message)}finally{stripeLoading.value=false}
}
async function queryOrderById(id:number){const j=await request('/api/v1/payments/orders/'+id);orders.value=[j];return j}
function openPaypal(){if(payment.value?.approve_url)window.open(payment.value.approve_url,'_blank')}
async function queryOrder(){
  if(!orderIdInput.value) return
  try { const j=await request('/api/v1/payments/orders/'+orderIdInput.value); orders.value=[j] }
  catch(e:any){ElMessage.error(e.message)}
}
onMounted(async()=>{await loadHome();const q=new URLSearchParams(window.location.search);if(q.get('payment')==='success'&&q.get('order_id')&&q.get('token')){try{const j=await request('/api/v1/payments/capture',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({order_id:Number(q.get('order_id')),provider_order_id:q.get('token')})});ElMessage.success('PayPal 支付状态：'+j.status)}catch(e:any){ElMessage.error(e.message)}}})
</script>

<template>
<div class="client-shell">
  <header class="nav">
    <div class="logo" @click="page='home';loadHome()">SHORT<span>DRAMA</span></div>
    <nav>
      <el-button text @click="page='home';loadHome()">首页</el-button>
      <el-button text @click="page='discover'">发现</el-button>
      <el-button text @click="page='orders'">我的订单</el-button>
      <el-button text @click="page='profile'">我的</el-button>
    </nav>
    <div class="locale">US / English</div>
  </header>

  <main class="content">
    <section v-if="page==='home'">
      <div class="hero"><div><p class="eyebrow">GLOBAL SHORT DRAMA</p><h1>Stories you can't stop watching.</h1><p>Personalized short dramas for viewers around the world.</p></div><el-button type="primary" size="large" @click="page='discover'">Explore now</el-button></div>
      <div class="section-head"><h2>For You</h2><el-button text :loading="loading" @click="loadHome">刷新</el-button></div>
      <div class="grid">
        <el-card v-for="x in dramas" :key="x.drama_id" class="drama-card" shadow="hover" @click="openDrama(x)">
          <div class="cover"><img v-if="x.cover" :src="x.cover"><span v-else>SHORT DRAMA</span></div>
          <h3>{{x.title}}</h3><p>{{x.reasons?.join(' · ')}}</p>
        </el-card>
      </div>
    </section>

    <section v-else-if="page==='discover'">
      <div class="page-title"><p class="eyebrow">DISCOVER</p><h1>Find your next obsession</h1><p>Browse recommendations by country and language.</p></div>
      <div class="grid"><el-card v-for="x in dramas" :key="x.drama_id" class="drama-card" @click="openDrama(x)"><div class="cover"><span>{{x.title}}</span></div><h3>{{x.title}}</h3><p>{{x.reasons?.join(' · ')}}</p></el-card></div>
    </section>

    <section v-else-if="page==='detail' && detail">
      <el-button text @click="page='home'">← 返回</el-button>
      <div class="detail">
        <div class="detail-cover"><img v-if="detail.cover" :src="detail.cover"><span v-else>SHORT DRAMA</span></div>
        <div class="detail-main"><p class="eyebrow">{{detail.country}} · {{detail.language}}</p><h1>{{detail.title}}</h1><p>{{detail.description}}</p>
          <div class="price" v-if="detail.is_paid">USD {{Number(detail.price_cents||499)/100}}</div><el-tag v-if="detail.is_paid" type="warning">付费剧</el-tag><el-tag v-else type="success">免费</el-tag>
          <div class="actions"><el-button v-if="detail.is_paid" type="primary" size="large" @click="page='checkout'">立即购买</el-button><el-button v-else type="primary" size="large">开始观看</el-button></div>
        </div>
      </div>
      <el-card class="episodes"><template #header><b>剧集列表</b></template><div class="episode-list"><el-button v-for="e in episodes" :key="e.id" @click="watchEpisode(e)">第 {{e.episode_no}} 集 {{e.title}}</el-button><el-empty v-if="!episodes.length" description="剧集数据待接口接入"/></div></el-card>
    </section>

    <section v-else-if="page==='checkout'">
      <div class="page-title"><p class="eyebrow">CHECKOUT</p><h1>Complete your purchase</h1><p>Secure payment · Instant unlock after confirmation</p></div>
      <el-card class="checkout"><h2>{{detail?.title||'Short Drama'}}</h2><div class="checkout-price">USD {{Number(detail?.price_cents||499)/100}}</div>
        <el-radio-group v-model="provider" class="providers"><el-radio-button label="STRIPE">Stripe</el-radio-button><el-radio-button label="PAYPAL">PayPal</el-radio-button></el-radio-group>
        <el-alert v-if="provider==='STRIPE'" title="Stripe 支付" description="订单创建后将返回 PaymentIntent client_secret；下一步接入 Stripe.js 完成卡支付。" type="info" show-icon/>
        <el-alert v-else title="PayPal Sandbox" description="创建订单后打开 PayPal 授权页面，授权完成后回到本客户端并执行 Capture。" type="info" show-icon/>
        <el-button type="primary" size="large" class="pay-btn" @click="buy">创建 {{provider}} 支付订单</el-button>
        <div v-if="payment?.client_secret" class="card-box"><el-input v-model="cardNumber" placeholder="Card number (Sandbox: 4242 4242 4242 4242)"/><div class="card-row"><el-input v-model="cardExpiry" placeholder="MM/YY"/><el-input v-model="cardCvc" placeholder="CVC"/></div><el-button type="success" :loading="stripeLoading" @click="payStripe">确认 Stripe 支付</el-button></div>
        <div v-if="payment" class="payment-result"><p>订单号：{{payment.order_no}}</p><p>状态：{{payment.status}}</p><p v-if="payment.client_secret">PaymentIntent：{{payment.provider_order_id}}</p><el-button v-if="payment.approve_url" type="success" @click="openPaypal">打开 PayPal</el-button></div>
      </el-card>
    </section>

    <section v-else-if="page==='orders'">
      <div class="page-title"><p class="eyebrow">ORDERS</p><h1>My Orders</h1><p>查询支付订单状态与解锁状态。</p></div>
      <el-card><div class="order-search"><el-input v-model="orderIdInput" placeholder="输入订单 ID"/><el-button type="primary" @click="queryOrder">查询</el-button></div><el-table :data="orders"><el-table-column prop="order_id" label="订单ID"/><el-table-column prop="order_no" label="订单号"/><el-table-column prop="provider" label="支付渠道"/><el-table-column prop="amount" label="金额"/><el-table-column prop="currency" label="币种"/><el-table-column prop="status" label="状态"/></el-table></el-card>
    </section>

    <section v-else>
      <div class="page-title"><p class="eyebrow">PROFILE</p><h1>My Profile</h1><p>Account, language and viewing preferences.</p></div>
      <el-card><el-descriptions :column="1" border><el-descriptions-item label="User ID">1</el-descriptions-item><el-descriptions-item label="Country">US</el-descriptions-item><el-descriptions-item label="Language">English</el-descriptions-item><el-descriptions-item label="Timezone">UTC</el-descriptions-item></el-descriptions></el-card>
    </section>
  </main>
</div>
</template>
