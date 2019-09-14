export default function() {
    // remove dangling facebook hashmark garbage
    if (location.hash === '#_=_') {
        history.replaceState(null, null, location.href.split("#")[0]);
    }
};
