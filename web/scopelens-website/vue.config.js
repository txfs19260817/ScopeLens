module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  "chainWebpack": config => {
    config.when(process.env.NODE_ENV === 'production', config => {
      config
          .entry('app')
          .clear()
          .add('./src/main.js');
      config.devServer.disableHostCheck(true);

      config.plugin('html').tap(args => {
        args[0].isProd = true;
        return args
      });

      config.set('externals', {
        vue: 'Vue',
        vuetify:"Vuetify",
        //"vue-router": "VueRouter",
        vuex: "Vuex",
        "vue-disqus":"VueDisqus",
        axios: 'axios',
        echarts: 'echarts'
      });
    });
  },
  "publicPath": './',
}