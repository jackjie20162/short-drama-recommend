import {createRouter,createWebHashHistory} from 'vue-router'
import Dashboard from '../views/dashboard/Dashboard.vue'
import DramaList from '../views/drama/DramaList.vue'
import DramaCreate from '../views/drama/DramaCreate.vue'
import EpisodeList from '../views/episode/EpisodeList.vue'
import OrderList from '../views/order/OrderList.vue'
import UserList from '../views/user/UserList.vue'
export default createRouter({history:createWebHashHistory(),routes:[
{path:'/',redirect:'/dashboard'},{path:'/dashboard',component:Dashboard},{path:'/dramas',component:DramaList},
{path:'/dramas/create',component:DramaCreate},{path:'/episodes',component:EpisodeList},{path:'/orders',component:OrderList},{path:'/users',component:UserList}]})