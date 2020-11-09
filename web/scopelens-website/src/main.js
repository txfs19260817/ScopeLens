import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import vuetify from './plugins/vuetify';
import http from "./api";
import i18n from './plugins/i18n';
import VueClipboard from 'vue-clipboard2';
import './registerServiceWorker'

// use this.$http alias
Vue.prototype.$http = http;

Vue.config.productionTip = false;
Vue.use(VueClipboard);
new Vue({
    router,
    vuetify,
    store,
    i18n,
    render: h => h(App)
}).$mount('#app');
