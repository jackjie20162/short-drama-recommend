import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/home/Home.vue'
import Discover from '../views/discover/Discover.vue'
import DramaDetail from '../views/drama/DramaDetail.vue'
import Watch from '../views/watch/Watch.vue'
import Checkout from '../views/checkout/Checkout.vue'
import Orders from '../views/orders/Orders.vue'
import Profile from '../views/profile/Profile.vue'

export default createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/home', component: Home },
    { path: '/discover', component: Discover },
    { path: '/drama/:id', component: DramaDetail },
    { path: '/watch/:dramaId/:episodeId', component: Watch },
    { path: '/checkout/:dramaId', component: Checkout },
    { path: '/orders', component: Orders },
    { path: '/profile', component: Profile },
  ],
})