import axios from 'axios'
import router from '../router'
import store from '../store'


// log errors
export const logErrors = (error) => {
    // Error
    if (error.response) {
        /*
         * The request was made and the server responded with a
         * status code that falls out of the range of 2xx
         */
        console.log(error.response.data);
        console.log(error.response.status);
        console.log(error.response.headers);
    } else if (error.request) {
        /*
         * The request was made but no response was received, `error.request`
         * is an instance of XMLHttpRequest in the browser and an instance
         * of http.ClientRequest in Node.js
         */
        console.log(error.request);
    } else {
        // Something happened in setting up the request and triggered an Error
        console.log('Error', error.message);
    }
    console.log(error.config);
}

// Response codes
export const ERROR = -1;
export const SUCCESS = 0;

// Axios instance
const http = axios.create({
    baseURL: process.env.VUE_APP_URL,
    timeout: 50000,
});

// Response interceptors
http.interceptors.response.use(
    response => {
        return response
    },
    error => {
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    console.log(error.response)
                    store.commit('LOADING_OFF');
                    store.dispatch("user/logout")
                    store.dispatch('snackbar/openSnackbar', {
                        "msg": "Token is invalid: " + error.response.data.msg + " Please login.",
                        "color": "error"
                    });
                    router.replace({
                        path: '/login',
                        query: {redirect: "/"} // 将跳转的路由path作为参数，登录成功后跳转到该路由
                    }).then(r => null)
            }
        }
        return Promise.reject(error.response)
    }
);

export default http;
