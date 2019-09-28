import {TextEncoder} from 'text-encoding-shim';
import {langToRegion} from '~/assets/l10n';
import bcp47 from 'bcp-47';

async function digestMessage(message) {
    const msgUint8 = new TextEncoder().encode(message);
    const hashBuffer = await crypto.subtle.digest('SHA-256', msgUint8);
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    return hashHex;
}

export default async function({store, error, $ga}) {
    let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
    if (userResponse.ok) {
        let {user} = await userResponse.json();
        store.commit('user/user', user);
        console.debug('Logged in user found');
        $ga.set('userId', await digestMessage(user.spotifyId));
    } else {
        store.commit('user/restore');
        if (!store.state.user.country) {
            let {language, region} = bcp47.parse(navigator.language);
            let regionFromLang = langToRegion[language];
            if (region && region.length === 2) {
                store.commit('user/country', region);
            } else if (!region && regionFromLang) {
                store.commit('user/country', regionFromLang);
            }
        }
        $ga.set('country-guess', store.getters['country']);
        try {
            await Promise.race([
                new Promise(resolve => grecaptcha.ready(resolve)),
                new Promise((_, reject) => setTimeout(() => reject('RECAPTCHA readiness timeout'), 1000))
            ]);
        } catch (e) {
            $ga.exception(e, true);
            console.error(e);
            error({message: 'Could not initialize RECAPTCHA'});
        }
        console.debug('User not logged in, RECAPTCHA ready');
    }
    console.debug(`User region is ${store.getters['user/country']}`)
};
