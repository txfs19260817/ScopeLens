<template>
    <v-container></v-container>
</template>

<script>
    export default {
        name: "Logout",
        mounted() {
            if (this.$store.state.user.isLogin) this.$store.dispatch("user/logout");
            this.$router.push("/");
        }
    }
</script>

<style scoped>

</style>