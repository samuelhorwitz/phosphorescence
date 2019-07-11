const path = require('path');
const fs = require('fs');

module.exports = {
  head: {
    title: 'Phosphorescence - Trance Playlist Builder',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { name: 'msapplication-TileColor', content: '#281b3d' },
      { name: 'msapplication-TileImage', content: '/ms-icon-144x144.png' },
      { name: 'theme-color', content: '#281b3d' }
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
      {rel: 'manifest', href: '/manifest.json'}
    ]
  },
  plugins: [
    {ssr: false, src: '~plugins/eos.js'},
    {ssr: false, src: '~plugins/monaco.js'}
  ],
  css: ['~/css/main.css'],
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
    }
  },
  env: {
    EOS_ORIGIN: process.env.EOS_ORIGIN,
    API_ORIGIN: process.env.API_ORIGIN,
    SCRIPTS_ORIGIN: process.env.SCRIPTS_ORIGIN
  }
};
