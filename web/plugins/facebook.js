export default (ctx, inject) => {
    if (localStorage.getItem('gdpr') !== 'true') {
        fbq('consent', 'revoke');
    }
    if (process.env.NODE_ENV === 'production') {
        fbq('init', '544556836299318');
        fbq('track', 'PageView');
    }
    let injectedFbq = {
        enable() {
            fbq('consent', 'grant');
        },
        disable() {
            fbq('consent', 'revoke');
        }
    };
    inject('fbq', injectedFbq);
    ctx.$fbq = injectedFbq;
}