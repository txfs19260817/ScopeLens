const snackbar = {
    namespaced: true,
    state: {
        msg: '', // snackbar message
        visible: false, // snackbar visible
        showClose: false, // show close button
        timeout: 6000, // visible duration
        color: 'error' // color
    },
    mutations: {
        OPEN_SNACKBAR(state, options) {
            state.visible = true;
            state.msg = options.msg;
        },
        CLOSE_SNACKBAR(state) {
            state.visible = false;
        },
        // set snackbar
        SET_SHOW_CLOSE(state, isShow) {
            state.showClose = isShow;
        },
        SET_TIMEOUT(state, timeout) {
            state.timeout = timeout;
        },
        SET_COLOR(state, color) {
            state.color = color;
        },
    },
    actions: {
        openSnackbar(context, options) {
            context.commit('OPEN_SNACKBAR', {
                msg: options.msg,
            });

            context.commit('SET_COLOR', {
                color: options.color,
            })

            setTimeout(() => {
                context.commit('CLOSE_SNACKBAR')
            }, context.state.timeout);
        }
    }
}

export default snackbar;