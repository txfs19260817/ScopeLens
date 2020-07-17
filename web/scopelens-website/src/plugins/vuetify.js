import Vue from 'vue';
import Vuetify from 'vuetify/lib';

Vue.use(Vuetify);

export default new Vuetify({
    theme:{
        themes: {
            light: {
                primary: '#4768A1',
                warning: '#EEE354'
            },
        }
    }
});
