const ParcelProxyServer = require('parcel-proxy-server');

process.env.NODE_ENV = 'development';

// configure the proxy server
const server = new ParcelProxyServer({
  entryPoint: 'index.html',
  parcelOptions: {
    // provide parcel options here
    // these are directly passed into the
    // parcel bundler
    //
    // More info on supported options are documented at
    // https://parceljs.org/api
    https: false
  },
  proxies: {
    // add proxies here
    '/setConfig': {
      target: 'http://nodemcu-controller'
    },
    '/getColorPaletteList': {
      target: 'http://nodemcu-controller'
    },
    '/getConnectedNodeMCUs': {
      target: 'http://nodemcu-controller'
    },
    '/getEffectList': {
      target: 'http://nodemcu-controller'
    },
  }
});

// the underlying parcel bundler is exposed on the server
// and can be used if needed
server.bundler.on('buildEnd', () => {
  console.log('Build completed!');
});

// start up the server
server.listen(1234, () => {
  console.log('Parcel proxy server has started');
});