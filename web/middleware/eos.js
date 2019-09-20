import {initialize} from '~/assets/eos.js';

export default async function({error}) {
    try {
        await initialize();
    } catch (e) {
        console.error(e);
        error({message: 'Could not initialize engine'});
    }
};
