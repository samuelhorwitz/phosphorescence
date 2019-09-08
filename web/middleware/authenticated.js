export default async function({store, redirect}) {
    let userResponse = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
    if (!userResponse.ok) {
        return redirect('/auth');
    }
    let {user} = await userResponse.json();
    store.commit('user/user', user);
};
