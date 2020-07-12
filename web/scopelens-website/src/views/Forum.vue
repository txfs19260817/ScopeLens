<template>
    <div class='forum'>
        <h1>Scope Lens Forum</h1>
        <v-subheader>A real-time discuss fourm powered by Disquz</v-subheader>
        <Disqus shortname='forum' />
    </div>
</template>

<script>
    import {Disqus} from 'vue-disqus'

    export default {
        name: "Forum",
        components: {
            Disqus
        }
    }
</script>

<style scoped>
    .forum{
        padding: 20px;
    }
</style>