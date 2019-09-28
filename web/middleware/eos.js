import {initialize} from '~/assets/eos.js';

export default async function({error, $ga}) {
    try {
        await initialize();
    } catch (e) {
        $ga.exception(e, true);
        console.error(e);
        error({message: 'Could not initialize engine'});
    }
};
