const path = require('path');
const fs = require('fs');

module.exports = {
  head: {
    title: 'Phosphorescence Trance Playlist Builder',
    meta: [
      { charset: 'utf-8' },
      { name: 'description', content: 'Intelligent trance music playlists for Spotify.' },
      { name: 'viewport', content: 'width=device-width, height=device-height, initial-scale=1, maximum-scale=1, viewport-fit=cover' },
      { name: 'msapplication-TileColor', content: '#281b3d' },
      { name: 'msapplication-TileImage', content: '/ms-icon-144x144.png' },
      { name: 'theme-color', content: '#281b3d' },
      { name: 'apple-mobile-web-app-capable', content: 'yes' },
      { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' }
    ],
    link: [
      {rel: 'apple-touch-icon', sizes: '57x57', href: '/apple-icon-57x57.png' },
      {rel: 'apple-touch-icon', sizes: '60x60', href: '/apple-icon-60x60.png' },
      {rel: 'apple-touch-icon', sizes: '72x72', href: '/apple-icon-72x72.png' },
      {rel: 'apple-touch-icon', sizes: '76x76', href: '/apple-icon-76x76.png' },
      {rel: 'apple-touch-icon', sizes: '114x114', href: '/apple-icon-114x114.png' },
      {rel: 'apple-touch-icon', sizes: '120x120', href: '/apple-icon-120x120.png' },
      {rel: 'apple-touch-icon', sizes: '144x144', href: '/apple-icon-144x144.png' },
      {rel: 'apple-touch-icon', sizes: '152x152', href: '/apple-icon-152x152.png' },
      {rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-icon-180x180.png' },
      {rel: 'icon', type: 'image/png', sizes: '192x192', href: '/android-icon-192x192.png' },
      {rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
      {rel: 'icon', type: 'image/png', sizes: '96x96', href: '/favicon-96x96.png' },
      {rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' },
      {href: "/iphone5_splash.png", media: "(device-width: 320px) and (device-height: 568px) and (-webkit-device-pixel-ratio: 2)",  rel: "apple-touch-startup-image"},
      {href: "/iphone6_splash.png", media: "(device-width: 375px) and (device-height: 667px) and (-webkit-device-pixel-ratio: 2)",  rel: "apple-touch-startup-image"},
      {href: "/iphoneplus_splash.png", media: "(device-width: 621px) and (device-height: 1104px) and (-webkit-device-pixel-ratio: 3)",  rel: "apple-touch-startup-image"},
      {href: "/iphonex_splash.png", media: "(device-width: 375px) and (device-height: 812px) and (-webkit-device-pixel-ratio: 3)",  rel: "apple-touch-startup-image"},
      {href: "/iphonexr_splash.png", media: "(device-width: 414px) and (device-height: 896px) and (-webkit-device-pixel-ratio: 2)",  rel: "apple-touch-startup-image"},
      {href: "/iphonexsmax_splash.png", media: "(device-width: 414px) and (device-height: 896px) and (-webkit-device-pixel-ratio: 3)",  rel: "apple-touch-startup-image"},
      {href: "/ipad_splash.png", media: "(device-width: 768px) and (device-height: 1024px) and (-webkit-device-pixel-ratio: 2)",  rel: "apple-touch-startup-image"},
      {href: "/ipadpro1_splash.png", media: "(device-width: 834px) and (device-height: 1112px) and (-webkit-device-pixel-ratio: 2)",  rel: "apple-touch-startup-image"},
      {href: "/ipadpro3_splash.png", media: "(device-width: 834px) and (device-height: 1194px) and (-webkit-device-pixel-ratio: 2)",  rel: "apple-touch-startup-image"},
      {href: "/ipadpro2_splash.png", media: "(device-width: 1024px) and (device-height: 1366px) and (-webkit-device-pixel-ratio: 2)",  rel: "apple-touch-startup-image"},
      {rel: 'manifest', href: '/manifest.json'}
    ],
    script: [
      {src: 'https://www.google.com/recaptcha/api.js?render=6LdfBboUAAAAAFv0977A1dWeer-eTy0IBmynzHcS'}
    ]
  },
  plugins: [
    {ssr: false, src: '~plugins/monaco.js'},
    {ssr: false, src: '~plugins/ios-rubberband-bg.js'},
    {ssr: false, src: '~plugins/consola.js'},
    {ssr: false, src: '~plugins/facebook.js'},
    {ssr: false, src: '~plugins/twitter.js'},
    {ssr: false, src: '~/directives/spotify-uri.js'}
  ],
  css: ['~/css/main.css', '~/css/fontawesome.css'],
  mode: 'spa',
  build: {
    extractCSS: true,
    extend(config) {
      config.output.globalObject = 'this';
    }
  },
  render: {
    bundleRenderer: {
      shouldPreload: (file, type) => {
        return ['script', 'style', 'font'].includes(type)
      }
    },
    static: {
      setHeaders: function(res) {
        if (process.env.NODE_ENV !== 'production') {
          res.setHeader('Expires', 'Tue, 21 Oct 2025 07:28:00 GMT');
        }
      }
    }
  },
  router: {
    middleware: ['gdpr', 'facebookurlcleaner', 'eos']
  },
  loading: '~/components/loading-bar.vue',
  loadingIndicator: '~/static/loading.html',
  env: {
    EOS_ORIGIN: process.env.EOS_ORIGIN,
    API_ORIGIN: process.env.API_ORIGIN,
    SCRIPTS_ORIGIN: process.env.SCRIPTS_ORIGIN,
    CONSOLA_LEVEL: process.env.NODE_ENV === 'production' ? process.env.CONSOLA_LEVEL : 0
  },
  modules: [
    ['@nuxtjs/google-analytics', {
      id: 'UA-148749300-1',
      autoTracking: {
        exception: true
      },
      fields: {
        allowLinker: true
      },
      disabled: () => localStorage.getItem('gdpr') !== 'true'
    }]
  ]
};
