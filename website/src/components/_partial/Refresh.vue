<template>
    <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark}" color="primary" fab dark small @click="refreshPage">
        <v-icon dark>mdi-refresh</v-icon>
    </v-btn>
</template>

<script>
    export default {
        name: "Refresh",
        inject: ["reload"],
        methods: {
            refreshPage() {
                this.reload();
            }
        },
    }
</script>

<style scoped>
    .v-btn{
        margin: 0 5px;
    }
</style>