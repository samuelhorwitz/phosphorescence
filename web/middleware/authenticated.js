import {isNewUser} from '~/assets/session';

export default async function({redirect}) {
    if (await isNewUser()) {
        return redirect('/auth');
    }
};
