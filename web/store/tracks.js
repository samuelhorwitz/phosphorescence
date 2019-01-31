import {initialize} from '~/store/abstract_playlist';
const {state, mutations, actions, getters} = initialize('tracks');
export {state, mutations, actions, getters};
