import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import http from "./api";


// use this.$http alias
Vue.prototype.$http = http;

Vue.config.productionTip = false;

new Vue({
    router,
    vuetify,
    store,
    render: h => h(App)
}).$mount('#app');
