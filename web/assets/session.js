import {createCookie, readCookie, eraseCookie} from '~/assets/cookies';

export async function login(code) {
    let {access, refresh, expires} = await getTokens(code);
    setCookies(access, refresh, expires);
}

export async function refreshUser() {
    await getAccessToken();
}

export function authorizeUser() {
    return new Promise((resolve, reject) => {
        let iframe = document.createElement('iframe');
        iframe.src = `${process.env.API_ORIGIN}/spotify/authorize`;
        iframe.style.display = 'none';
        document.head.appendChild(iframe);
        let intervalId = setInterval(() => {
            if (readCookie('spotify_access')) {
                clearInterval(intervalId);
                clearTimeout(timeoutId);
                document.head.removeChild(iframe);
                resolve();
            }
        }, 100);
        let timeoutId = setTimeout(() => {
            clearInterval(intervalId);
            document.head.removeChild(iframe);
            reject();
            location.href = `${process.env.API_ORIGIN}/spotify/authorize`;
        }, 5000);
    });
}

function setCookies(access, refresh, expires) {
    createCookie('spotify_access', access, expires, true);
    if (refresh) {
        createCookie('spotify_refresh', refresh, 60 * 60 * 12);
    }
    else {
        eraseCookie('spotify_refresh');
    }
    createCookie('seen_user', true, 60 * 60 * 24 * 365);
}

export async function getAccessToken() {
    let accessToken = readCookie('spotify_access');
    if (!accessToken) {
        let refreshToken = getRefreshToken();
        if (!refreshToken) {
            await authorizeUser();
            return getAccessToken();
        }
        let {access, refresh, expires} = await getTokens(refreshToken, true);
        setCookies(access, refresh, expires);
        return access;
    }
    return accessToken;
}

export function quickReturnAccessTokenWithoutGuarantee() {
    getAccessToken();
    return readCookie('spotify_access');
}

export function accessTokenExists() {
    return !!readCookie('spotify_access');
}

export function isNewUser() {
    return !readCookie('seen_user');
}

function getRefreshToken() {
    return readCookie('spotify_refresh');
}

async function getTokens(code, isRefresh) {
    let type
    if (!isRefresh) {
        type = 'auth';
    } else {
        type = 'refresh';
    }
    let response = await fetch(`${process.env.API_ORIGIN}/spotify/tokens?code=${code}&type=${type}`, {cache: 'no-cache'});
    let {status, statusText} = response;
    if (status != 200) {
        throw new Error(`Could not get tokens: ${statusText}`);
    }
    let {access, refresh, expires} = await response.json();
    return {access, refresh, expires};
}

export function logout() {
    eraseCookie('spotify_access');
    eraseCookie('spotify_refresh');
    eraseCookie('seen_user');
    sessionStorage.clear();
    localStorage.clear();
}

export async function getUsersCountry() {
    let response = await fetch('https://api.spotify.com/v1/me', {
        method: 'GET',
        headers: {
            Authorization: `Bearer ${await getAccessToken()}`
        }
    });
    if (response && response.status == 200) {
        let {country} = await response.json();
        return country;
    }
    return null;
}

export async function getTrackWithFeatures(trackId) {
    let fullResponse = {};
    let trackResponse = await fetch(`https://api.spotify.com/v1/tracks/${trackId}`, {
        method: 'GET',
        headers: {
            Authorization: `Bearer ${await getAccessToken()}`
        }
    });
    if (trackResponse && trackResponse.status == 200) {
        fullResponse.track = await trackResponse.json();
    } else {
        return null;
    }
    let featuresResponse = await fetch(`https://api.spotify.com/v1/audio-features?ids=${trackId}`, {
        method: 'GET',
        headers: {
            Authorization: `Bearer ${await getAccessToken()}`
        }
    });
    if (featuresResponse && featuresResponse.status == 200) {
        let allFeatures = await featuresResponse.json();
        if (allFeatures.audio_features && allFeatures.audio_features.length > 0) {
            fullResponse.features = allFeatures.audio_features[0];
        } else {
            return null;
        }
    } else {
        return null;
    }
    return fullResponse;
}
