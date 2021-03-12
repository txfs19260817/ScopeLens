import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import colors from 'vuetify/lib/util/colors'


Vue.use(Vuetify);

let h = new Date().getHours();

export default new Vuetify({
    theme:{
        dark: h >=18 || h <= 6,
        themes: {
            light: {
                primary: '#4768A1',
                secondary: colors.blueGrey,
                warning: '#EEE354'
            },
            dark: {
                primary: colors.amber,
                secondary: colors.grey,
                error: colors.pink,
                success: '#4caf50'
            }
        }
    }
});
