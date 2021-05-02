import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/components/Login'
import Layout from '@/components/Layout'

Vue.use(Router)

const Home = { template: "<div>Home</div>" }


export default new Router({
  routes: [
    {
      path: '/',
      name: 'Layout',
      component: Layout
    },
    {
      path: '/app',
      name: 'App',
      component: Home
    },
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/layout',
      name: 'Layout',
      component: Layout
    }
  ]
})
