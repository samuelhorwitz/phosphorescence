# Phosphorescence Training Module
This stuff is kind of a mess and I had to clean it up for release so a lot of files containing proprietary Spotify data are missing. You will have to get this yourself with your own API token, I don't have the rights to distribute it in this repo. Basically, the idea is that you are going to download a lot of Spotify tracks with audio features and use the training UI to specify your evocativeness opinion. This UI is fairly buggy, the OAuth doesn't properly refresh so after about an hour you'll no longer be able to use the player/playlist buttons. **Your results are stored in localstorage** so please, frequently dump them to disk, using the snippet below under "Get the data!". Once dumped, it's a good idea to clear localstorage to start anew, but please make sure your data is safetly saved first or else your classifying for that period of time will be lost. Running that from the developer console should put the JSON data in your clipboard for you to paste into a file. If your developer tools/OS do not support `copy` (which puts the content in the system clipboard) you'll have manually copy and paste.

Once you have all this, you'll run the Python/Keras code and it should perform regression classification. Each axis is a separate regression. There is a regression performed on the aetherealness spectrum and on the primordialness spectrum.

## Get the data!

When you are tired of training, export your training data out of localstorage. This puts it in your clipboard to paste into a JSON file.

```
var tracks = {};
for(var key in localStorage) {
    if (key === '_spharmony_device_id') continue;
    tracks[key]=JSON.parse(localStorage.getItem(key));
}
copy(tracks);
```

## CouchDB

You will need to have CouchDB set up to import the JSON and export it again into cleaned up CSV files.

## x/y

- x is chthonic (0) to aethereal (1)
- y is transcendental (0) to primordial (1)


