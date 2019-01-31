export function createCookie(name, value, seconds) {
    let expires;
    if (seconds) {
        let date = new Date();
        date.setTime(date.getTime() + (seconds * 1000));
        expires = `; expires=${date.toGMTString()}`;
    }
    else {
        expires = '';
    }
    document.cookie = `${name}=${value}${expires}; path=/`;
}

export function readCookie(name) {
    let nameEQ = `${name}=`;
    let ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1, c.length);
        }
        if (c.indexOf(nameEQ) === 0) {
            return c.substring(nameEQ.length, c.length);
        }
    }
    return null;
}

export function eraseCookie(name) {
    createCookie(name, '', -1);
}
