import {langToRegion} from '~/assets/l10n';
import bcp47 from 'bcp-47';

const isAdmin = 'dimension2';
const spotifyRegionGuess = 'dimension3';

export default async function({store, error, $ga}) {
    let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
    if (userResponse.ok) {
        let {user} = await userResponse.json();
        store.commit('user/user', user);
        console.debug('Logged in user found');
        if (user.spotifyId === '126149108') {
            $ga.set(isAdmin, 'true');
        }
        $ga.set('userId', user.gaId);
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
        $ga.set(spotifyRegionGuess, store.getters['user/country']);
        try {
            await Promise.race([
                new Promise(resolve => grecaptcha.ready(resolve)),
                new Promise((_, reject) => setTimeout(() => reject('RECAPTCHA readiness timeout'), 10000))
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
