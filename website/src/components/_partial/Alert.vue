<template>
    <div>
        <v-alert
                v-model="alert"
                border="left"
                close-text="Close Alert"
                color="red lighten-2"
                dark
                dismissible
                icon="mdi-pokemon-go"
                outlined
                transition="fade-transition"
        >
            {{ $t("alert") }}
        </v-alert>
    </div>
</template>

<script>
    export default {
        name: "Alert",
        data () {
            return {
                alert: true,
            }
        },
    }
</script>

<style scoped>

</style>