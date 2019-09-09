export function getTrackTag(track) {
    let normalizedName = normalizeString(track.name);
    let normalizedArtist = getPrimaryArtist(track.artists);
    return `${normalizedName} - ${normalizedArtist}`;
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
