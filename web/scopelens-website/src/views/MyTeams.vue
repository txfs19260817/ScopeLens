<template>
    <v-container>
        <v-tabs color="black" slider-color="primary" background-color="bg_primary">
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
    import {getUploadedTeamsByUsername, getLikedTeamsByUsername} from "../api/team";
    import {SUCCESS, logErrors} from "../api";
    import ResultsLayout from "../components/layouts/ResultsLayout";

    export default {
        name: 'MyTeams',
        components: {
            ResultsLayout
        },
        data: () => ({
            teams: [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},],// dummy obj * pageSize
            // page
            total: 1, // total data size
            pageSize: 12, // data size per page
            curPage: 1,
            tabs: [
                'Liked',
                'Uploaded'
            ],
            curTab: 0,
        }),
        methods: {
            async getLikedTeams(page) {
                // loading
                this.$store.commit('LOADING_ON')
                await getLikedTeamsByUsername(page, this.$store.state.user.username).then(res => {
                    if (res.data.code === SUCCESS) {
                        this.total = res.data.data.total
                        this.teams = res.data.data.teams
                    } else {
                        // In this case, if user never liked any team, no error msg will be shown.
                        if (res.data.msg.toString().includes("$in")) {
                            this.$store.dispatch('snackbar/openSnackbar', {
                                "msg": "Failed to retrieve teams from server! " + res.data.msg,
                                "color": "error"
                            });
                        }
                    }
                }).catch(error => {
                    logErrors(error)
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": "Failed to connect to server! ",
                        "color": "error"
                    });
                }).finally(() => {
                    this.$store.commit('LOADING_OFF')
                })
            },
            async getUploadedTeams(page) {
                // loading
                this.$store.commit('LOADING_ON')
                await getUploadedTeamsByUsername(page, this.$store.state.user.username).then(res => {
                    if (res.data.code === SUCCESS) {
                        this.total = res.data.data.total
                        this.teams = res.data.data.teams
                    } else {
                        this.$store.dispatch('snackbar/openSnackbar', {
                            "msg": "Failed to retrieve teams from server! " + res.data.msg,
                            "color": "error"
                        });
                    }
                }).catch(error => {
                    logErrors(error)
                    this.$store.dispatch('snackbar/openSnackbar', {
                        "msg": "Failed to connect to server! ",
                        "color": "error"
                    });
                }).finally(() => {
                    this.$store.commit('LOADING_OFF')
                })
            },
            tabChange(n) {
                this.curPage = 1
                this.curTab = n
                this.teams = [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},]
                if (n === 0) {
                    this.getLikedTeams(1);
                } else {
                    this.getUploadedTeams(1);
                }
            },
            pageChange(p) {
                if (this.curTab === 0) {
                    return this.getLikedTeams(p);
                } else {
                    return this.getUploadedTeams(p);
                }
            }
        },
        created() {
            this.getLikedTeams(1);
        },
        computed: {
            pageLen() {
                return Math.ceil(this.total / this.pageSize)
            }
        }
    }
</script>

<style scoped>
    .v-card {
        transition: opacity .3s ease-in-out;
    }

    .v-card:not(.on-hover) {
        opacity: 0.88;
    }
</style>