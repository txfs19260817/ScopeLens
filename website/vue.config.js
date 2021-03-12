const BASE_URL = process.env.NODE_ENV === 'production'//If want to put the website into second-level directory, do change it to dir name
  ? '/'
  : '/';
module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  publicPath: BASE_URL,
  pwa: {//Setting PWA
    name: 'ScopeLens',
    themeColor: '#4768a1',
    msTileColor: '#4768a1',
    appleMobileWebAppCapable: 'yes',
    appleMobileWebAppStatusBarStyle: '#4768a1',
    // configure the workbox plugin
    workboxPluginMode: 'GenerateSW',//Auto generate service worker(will be registered in registerServiceWorker.js)
    workboxOptions: {
      runtimeCaching: [
        {
          urlPattern: new RegExp('.*?'),//cache if hited
          handler: 'NetworkFirst',//only fetch cache when network is down
          options: {
            networkTimeoutSeconds: 20,//wait the network for 20s 
            cacheName: 'api-cache',
            cacheableResponse: {
              statuses: [0, 200]
            }
          }
        }
      ]
    }},
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
          vuetify: "Vuetify",
          "vue-i18n": "VueI18n",
          "vue-router": "VueRouter",
          vuex: "Vuex",
          "vue-disqus": "VueDisqus",
          axios: 'axios',
          echarts: 'echarts'
        });
      });
    },
    "publicPath": './',
    pluginOptions: {
      i18n: {
        locale: 'zh',
        fallbackLocale: 'en',
        localeDir: 'locales',
        enableInSFC: false
      }
    }
  }
