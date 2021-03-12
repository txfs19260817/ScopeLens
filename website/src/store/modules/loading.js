const loading = {
    namespaced: false,
    state: {
        loading: false,
    },
    mutations: {
        LOADING_ON(state) {
            state.loading = true;
        },
        LOADING_OFF(state) {
            state.loading = false;
        },
    },
    actions: {
        loadingOn(context) {
            context.commit('LOADING_ON')
        },
        loadingOff(context) {
            context.commit('LOADING_OFF')
        },
    },
};

export default loading;