export default (ctx, inject) => {
    let oldTwq;
    let wasInit = false;
    function init() {
        if (wasInit) {
            return;
        }
        wasInit = true;
        twq('init', 'o2i1s');
        twq('track', 'PageView');
    }
    if (localStorage.getItem('gdpr') === 'true' && process.env.NODE_ENV === 'production') {
        init();
    }
    let injectedTwq = {
        enable() {
            if (oldTwq) {
                window.twq = oldTwq;
            }
            init();
        },
        disable() {
            oldTwq = twq;
            window.twq = function(){};
        }
    };
    inject('twq', injectedTwq);
    ctx.$twq = injectedTwq;
}