<template>
    <div class='forum'>
        <h1>{{ $t('forum.title') }}</h1>
        <v-subheader>{{ $t('forum.subtitle') }}</v-subheader>
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