<template>
    <v-container>
        <v-tabs slider-color="primary" background-color="bg_primary">
            <v-tab v-for="(item, i) in tabs" :key="i" @change="tabChange(i)" class="font-weight-bold title">{{ item }}
            </v-tab>
        </v-tabs>
        <ResultsLayout :teams="teams"></ResultsLayout>
        <v-col>
            <v-pagination v-model="curPage" :length="pageLen" total-visible="8" @input="pageChange"></v-pagination>
        </v-col>
    </v-container>
</template>

<script>
    import {getTeams, getTeamsByLikes} from "../api/team";
    import {SUCCESS, logErrors} from "../api";
    import ResultsLayout from "../components/layouts/ResultsLayout";

    export default {
        name: 'Home',
        components: {
            ResultsLayout
        },
        data: () => ({
            teams: [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},],// dummy obj * pageSize
            // page
            total: 1, // total data size
            pageSize: 12, // data size per page
            curPage: 1,
            curTab: 0,
        }),
        methods: {
            async getTeamsByTime(page) {
                // loading
                this.$store.commit('LOADING_ON')
                await getTeams(page).then(res => {
                    if (res.data.code === SUCCESS) {
                        this.total = res.data.data.total
                        this.teams = res.data.data.teams
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": this.$t('api.thenError') + res.data.msg,
                            "color": "error"
                        });
                    }
                }).catch(error => {
                    logErrors(error)
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t('api.catchError'),
                        "color": "error"
                    });
                }).finally(() => {
                    this.$store.commit('LOADING_OFF')
                })
            },
            async getMostLikedTeams(page) {
                // loading
                this.$store.commit('LOADING_ON')
                await getTeamsByLikes(page).then(res => {
                    if (res.data.code === SUCCESS) {
                        this.total = res.data.data.total
                        this.teams = res.data.data.teams
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": this.$t('api.thenError') + res.data.msg,
                            "color": "error"
                        });
                    }
                }).catch(error => {
                    logErrors(error)
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": this.$t('api.catchError'),
                        "color": "error"
                    });
                }).finally(() => {
                    this.$store.commit('LOADING_OFF')
                })
            },
            tabChange(n) {
                this.curPage = 1
                this.curTab = n
                if (n === 0) {
                    this.getTeamsByTime(1);
                } else {
                    this.getMostLikedTeams(1);
                }
            },
            pageChange(p) {
                if (this.curTab === 0) {
                    return this.getTeamsByTime(p);
                } else {
                    return this.getMostLikedTeams(p);
                }
            }
        },
        created() {
            this.getTeamsByTime(1);
        },
        computed: {
            pageLen() {
                return Math.ceil(this.total / this.pageSize)
            },
            tabs() {
                return [
                    this.$i18n.t('home.tabs.latest'),
                    this.$i18n.t('home.tabs.trending'),
                ]
            }
        }
    }
</script>

<style scoped>

</style>