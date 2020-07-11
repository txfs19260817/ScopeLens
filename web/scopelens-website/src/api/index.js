import axios from 'axios'
import router from '../router'


// Response codes
export const ERROR = -1;
export const SUCCESS = 0;

// Axios instance
const http = axios.create({
    baseURL: process.env.VUE_APP_URL,
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
                    router.replace({
                        path: 'login',
                        query: {redirect: "/"} // 将跳转的路由path作为参数，登录成功后跳转到该路由
                    }).then(r => null)
            }
        }
        return Promise.reject(error.response)
    }
);

export default http;
