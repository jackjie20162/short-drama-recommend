import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/home/Home.vue'
import Discover from '../views/discover/Discover.vue'
import DramaDetail from '../views/drama/DramaDetail.vue'
import Watch from '../views/watch/Watch.vue'
import Checkout from '../views/checkout/Checkout.vue'
import Orders from '../views/orders/Orders.vue'
import Profile from '../views/profile/Profile.vue'
import Login from '../views/auth/Login.vue'
import { token } from '../api/client'

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/login', component: Login },
    { path: '/home', component: Home },
    { path: '/discover', component: Discover },
    { path: '/drama/:id', component: DramaDetail },
    { path: '/watch/:dramaId/:episodeId', component: Watch },
    { path: '/checkout/:dramaId', component: Checkout },
    { path: '/orders', component: Orders },
    { path: '/orders/:id', component: () => import('../views/orders/OrderDetail.vue') },
    { path: '/profile', component: Profile },
  ],
}).beforeEach((to)=>{ if (['/checkout','/orders','/profile'].some(p=>to.path.startsWith(p)) && !token()) return '/login' })