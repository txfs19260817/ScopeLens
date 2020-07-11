<template>
    <div>
        <v-snackbar v-model="visible" :color="color.color" top>
            {{ msg }}
            <v-btn v-if="showClose" text @click="close">Close</v-btn>
        </v-snackbar>
    </div>
</template>

<script>
    import {mapState} from 'vuex'
    export default {
        name: "Snackbar",
        computed: {
            ...mapState('snackbar', ['msg', 'showClose', 'color', 'timeout']),
            visible: {
                get() {
                    return this.$store.state.snackbar.visible
                },
                set(v) {
                    this.$store.commit('snackbar/CLOSE_SNACKBAR')
                }
            }
        },
        methods: {
            close() {
                this.$store.commit('snackbar/CLOSE_SNACKBAR')
            }
        },
    }
</script>

<style scoped>

</style>