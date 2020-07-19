<template>
    <v-bottom-sheet v-model="sheet">
        <template v-slot:activator="{ on, attrs }">
            <v-btn color="primary" dark v-bind="attrs" v-on="on"> {{$t('lang.btn')}} </v-btn>
        </template>
        <v-list>
            <v-subheader> {{$t('lang.sheetHeader')}} </v-subheader>
            <v-list-item v-for="(lang, i) in languages" :key="i" @click="switchLang(i)">
                <v-list-item-title>{{ lang.name }}</v-list-item-title>
            </v-list-item>
        </v-list>
    </v-bottom-sheet>
</template>

<script>
    export default {
        name: "LanguageSelector",
        data() {
            return {
                sheet: false,
                languages: [
                    {
                        name:"简体中文",
                        value:"zh"
                    },
                    {
                        name:"English",
                        value:"en"
                    }
                ],
            }
        },
        methods: {
            switchLang(i) {
                this.sheet = false
                this.$i18n.locale = this.languages[i].value;
                // save lang in localStorage
                localStorage.setItem('lang', this.languages[i].value);
            }
        },
        created() {
            this.$i18n.locale = localStorage.getItem('lang') || 'zh-hans'
        }
    }
</script>

<style scoped>

</style>