import {loginRequest, registerRequest} from "../../api/auth";
import {ERROR} from "../../api";

const user = {
    namespaced: true,
    state: {
        username: localStorage.getItem('username'),
        token: localStorage.getItem('token'),
        isLogin: localStorage.getItem('token') !== null,
    },
    mutations: {
        LOGIN(state, options) {
            state.username = options.username;
            state.token = options.token;
            state.isLogin = true;
        },
        LOGOUT(state) {
            state.username = null;
            state.token = null;
            state.isLogin = false;
        },
        UPDATE_LOGIN_STATUS(state, options) {
            state.isLogin = options.isLogin;
        }
    },
    actions: {
        async register({commit, dispatch}, options) {
            let success = false;
            commit('LOADING_ON', null, {root: true})

            const res = await registerRequest(options.data, options.recaptcha);
            if (res.data.code === ERROR) {
                dispatch('snackbar/openSnackbar', {
                    "msg": res.data.msg,
                    "color": "error"
                }, {root: true});
            } else {
                dispatch('snackbar/openSnackbar', {
                    "msg": "Register success! Please login with your new account. ",
                    "color": "success"
                }, {root: true});
                success = true;
            }

            commit('LOADING_OFF', null, {root: true})
            return success
        },
        async login({commit, dispatch}, options) {
            let success = false;
            commit('LOADING_ON', null, {root: true})

            const res = await loginRequest(options.data);
            console.log(res);
            if (res.data.code === ERROR) {
                dispatch('snackbar/openSnackbar', {
                    "msg": "Authentication failed. Please check your username and password. Error: " + res.data.msg,
                    "color": "error"
                }, {root: true});
            } else {
                saveToken(options.data.username, res.data.data, commit);
                success = true;
            }

            commit('LOADING_OFF', null, {root: true})
            return success
        },
        logout(context, options) {
            localStorage.removeItem('username');
            localStorage.removeItem('token');
            context.commit('LOGOUT')
        }
    },
};

function saveToken(username, token, commit) {
    token = "Bearer " + token
    localStorage.setItem('username', username)
    localStorage.setItem('token', token)

    commit('LOGIN', {
        username: username,
        token: token,
    })
}

export default user;