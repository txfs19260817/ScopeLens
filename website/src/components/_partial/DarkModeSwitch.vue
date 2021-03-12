<template>
    <v-btn :class="{'grey--text text--darken-4': $vuetify.theme.dark}" color="primary" fab dark small @click="$vuetify.theme.dark = !$vuetify.theme.dark">
        <span v-if="!$vuetify.theme.dark">
            <v-icon dark>mdi-lightbulb-off</v-icon>
        </span>
        <span v-else>
            <v-icon dark>mdi-lightbulb-on</v-icon>
        </span>
    </v-btn>
</template>

<script>
    export default {
        name: "DarkModeSwitch",
    }
</script>

<style scoped>
    .v-btn{
        margin: 0 5px;
    }
</style>