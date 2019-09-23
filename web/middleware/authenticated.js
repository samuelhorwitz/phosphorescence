import {langToRegion} from '~/assets/l10n';
import bcp47 from 'bcp-47';

export default async function({store, error}) {
    let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
    if (userResponse.ok) {
        let {user} = await userResponse.json();
        store.commit('user/user', user);
        console.debug('Logged in user found');
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
        try {
            await Promise.race([
                new Promise(resolve => grecaptcha.ready(resolve)),
                new Promise((_, reject) => setTimeout(() => reject('RECAPTCHA readiness timeout'), 1000))
            ]);
        } catch (e) {
            console.error(e);
            error({message: 'Could not initialize RECAPTCHA'});
        }
        console.debug('User not logged in, RECAPTCHA ready');
    }
    console.debug(`User region is ${store.getters['user/country']}`)
};
