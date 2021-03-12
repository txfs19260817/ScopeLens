const prodPlugins = [];
if (process.env.NODE_ENV === 'production') {
  prodPlugins.push('transform-remove-console')
}

module.exports = {
  presets: [
    '@vue/cli-plugin-babel/preset'
  ],
  "plugins": [
    ...prodPlugins, // 发布产品时候的插件数组
  ]
}
