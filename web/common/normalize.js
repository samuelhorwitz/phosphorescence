export async function getTrackTag(track) {
    let normalizedName = normalizeString(track.name);
    let normalizedArtist = getPrimaryArtist(track.artists);
    let tag = `${normalizedName}-${normalizedArtist}`;
    let hash = await sha256(tag);
    return toHex(hash);
}

function normalizeString(str) {
    return str.split(/[(\-\[]/, 1)[0].trim().normalize('NFD').replace(/[^a-zA-Z0-9]/g, '').toLowerCase();
}

function getPrimaryArtist(artists) {
    if (isSignum(artists)) {
        return "signum";
    }
    return normalizeString(artists[0].name);
}

async function sha256(message) {
  let encoder = new TextEncoder();
  let data = encoder.encode(message);
  let hash = await crypto.subtle.digest('SHA-256', data);
  return hash;
}

function toHex(buf) {
  let hashArray = Array.from(new Uint8Array(buf));
  let hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
  return hashHex;
}

// See https://github.com/samuelhorwitz/phosphorescence/issues/36
function isSignum(artists) {
    if (artists.length < 2) {
        return false;
    }
    let primaryArtist = normalizeString(artists[0].name);
    let secondaryArtist = normalizeString(artists[1].name);
    if ((primaryArtist == 'ronhagen' && secondaryArtist == 'pascalm') ||
        (primaryArtist == 'pascalm' && secondaryArtist == 'ronhagen') ||
        primaryArtist == 'ronhagenpascalm' || primaryArtist == 'pascalmronhagen') {
        return true;
    }
    return false;
}
