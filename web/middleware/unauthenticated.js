import {isNewUser} from '~/assets/session';

export default function({redirect}) {
    if (!isNewUser()) {
        return redirect('/');
    }
};
