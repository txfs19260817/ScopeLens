<template>
    <v-app>
        <Navbar/>
        <v-main>
            <Snackbar style="height: 0"/>
            <v-slide-x-reverse-transition mode="out-in">
                <keep-alive :include="r">
                    <router-view/>
                </keep-alive>
            </v-slide-x-reverse-transition>
        </v-main>
    </v-app>
</template>

<script>
    import Navbar from "./components/Navbar";
    import Snackbar from "./components/_partial/Snackbar";

    export default {
        name: 'App',
        components: {
            Navbar,
            Snackbar,
        },
        data: () => ({
            r: ["Home", "Search", "MyTeams"]
        }),
    };
</script>
