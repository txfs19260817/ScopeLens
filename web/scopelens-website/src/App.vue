<template>
    <v-app>
        <Navbar/>
        <v-main>
            <Alert v-if="true"></Alert>
            <Snackbar style="height: 0"/>
            <v-slide-x-reverse-transition mode="out-in">
                <keep-alive v-if="isRouterAlive" :include="r">
                    <router-view/>
                </keep-alive>
            </v-slide-x-reverse-transition>
        </v-main>
    </v-app>
</template>

<script>
    import Navbar from "./components/Navbar";
    import Snackbar from "./components/_partial/Snackbar";
    import Alert from "./components/_partial/Alert";

    export default {
        name: 'App',
        provide () {
            return {
                reload: this.reload
            }
        },
        components: {
            Navbar,
            Snackbar,
            Alert
        },
        data: () => ({
            r: ["Home", "Search", "MyTeams"],
            isRouterAlive: true
        }),
        methods: {
            reload () {
                this.isRouterAlive = false
                this.$nextTick(function () {
                    this.isRouterAlive = true
                })
            }
        }
    };
</script>
