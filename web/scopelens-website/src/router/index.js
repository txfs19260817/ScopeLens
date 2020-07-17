import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'
import Home from '../views/Home.vue'
import {checkToken} from "../api/auth";
import {SUCCESS} from "../api";

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
        path: '/myteams',
        name: 'MyTeams',
        component: () => import('../views/MyTeams.vue'),
        meta: {
            requireAuth: true
        },
    },
    {
        path: '/upload',
        name: 'Upload',
        component: () => import('../views/Upload.vue'),
        meta: {
            requireAuth: true
        },
    },
    {
        path: '/search',
        name: 'Search',
        component: () => import('../views/Search.vue')
    },
    {
        path: '/usage',
        name: 'Usage',
        component: () => import('../views/Usage.vue')
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
    },
    {
        path: '/team/:id',
        name: 'Team',
        component: () => import('../views/Team.vue')
    },
    {
        path: '/logout',
        name: 'Logout',
        component: () => import('../views/Logout.vue')
    },
]

const router = new VueRouter({
    //mode: 'history',
    routes: routes,
});

router.beforeEach((to, from, next) => {
    // Already logged in, redirect to Homepage
    if (to.name === 'Login' && store.state.user.isLogin) next({name: 'Home'})
    // Not login, redirect to Login page
    else if (to.meta.requireAuth) {
        if (!store.state.user.isLogin) next({name: 'Login'})
        else {
            // check if token is still valid
            checkToken(store.state.user.token).then(res => {
                console.log(res.data)
                if (res.data.code === SUCCESS) next()
                // else guarded by http.interceptors.response
            })
        }
    } else next()
})

export default router
