import { defineConfig } from '#q-app/wrappers';

export default defineConfig((_ctx) => {
  return {
    boot: ['axios'],

    css: ['app.scss'],

    extras: ['roboto-font', 'material-icons'],

    build: {
      target: {
        browser: ['es2022', 'firefox115', 'chrome115', 'safari14'],
        node: 'node20',
      },
      typescript: {
        strict: true,
      },
      vueRouterMode: 'history',
    },

    devServer: {
      port: 9000,
      proxy: {
        '/api': {
          target: 'http://localhost:3000',
          changeOrigin: true,
        },
      },
    },

    framework: {
      config: { dark: true },
      plugins: ['Notify', 'Dialog', 'Loading'],
    },

    animations: [],

    ssr: {
      pwa: false,
      prodPort: 3001,
      middlewares: ['render'],
    },

    pwa: {
      workboxMode: 'GenerateSW',
    },

    cordova: {},

    capacitor: {
      hideSplashscreen: true,
    },

    electron: {
      preloadScripts: ['electron-preload'],
      inspectPort: 5858,
      bundler: 'packager',
      packager: {},
      builder: {
        appId: 'frontend',
      },
    },

    bex: {
      extraScripts: [],
    },
  };
});
