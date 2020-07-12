import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'
import Home from '../views/Home.vue'

Vue.use(VueRouter);

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue')
    },
    {
        path: '/upload',
        name: 'Upload',
        component: () => import('../views/Upload.vue')
    },
    {
        path: '/forum',
        name: 'Forum',
        component: () => import('../views/Forum.vue')
    },
    {
        path: '/about',
        name: 'About',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
    }
]

const router = new VueRouter({
    //mode: 'history',
    routes: routes,
});

router.beforeEach((to, from, next) => {
    // Already logged in, redirect to Homepage
    if (to.name === 'Login' && store.state.user.isLogin) next({name: 'Home'})
    // Not login, redirect to Login page
    else if (to.name === 'Upload' && !store.state.user.isLogin) next({name: 'Login'})
    else next()
})

export default router
