<template>
    <v-app>
        <Navbar />
        <v-main>
            <Snackbar style="height: 0" />
            <v-slide-x-reverse-transition mode="out-in">
                <router-view />
            </v-slide-x-reverse-transition>
        </v-main>
    </v-app>
</template>

<script>
    import Navbar from "./components/Navbar";
    import Snackbar from "./components/_partial/Snackbar";

    export default {
        name: 'App',
        components:{
            Navbar,
            Snackbar,
        },
        data: () => ({
            //
        }),
    };
</script>
