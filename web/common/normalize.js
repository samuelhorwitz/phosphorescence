export function getTrackTag(track) {
    let normalizedName = normalizeString(track.name);
    let normalizedArtist = normalizeString(track.artists[0].name);
    return `${normalizedName} - ${normalizedArtist}`;
}

function normalizeString(str) {
    return str.split(/[(\-\[]/, 1)[0].trim().normalize('NFD').replace(/[^a-zA-Z0-9]/g, '').toLowerCase();
}
