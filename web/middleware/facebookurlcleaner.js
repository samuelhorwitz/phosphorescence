export default function({$ga}) {
    // remove dangling facebook hashmark garbage
    if (location.hash === '#_=_') {
        $ga.set('referrer', 'https://phosphor.me/facebook-oauth');
        history.replaceState(null, null, location.href.split("#")[0]);
    }
};
