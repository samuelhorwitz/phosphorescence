export default function({$ga}) {
    let gdpr = document.getElementById('gdpr');
    let consent = localStorage.getItem('gdpr');
    if (consent !== 'true') {
        $ga.disable(); // this should already be disabled but just in case
        document.body.addEventListener('click', dismissGDPR, {capture: true, once: true});
        let agreeButton = document.getElementById('gdpr-consent');
        agreeButton.onclick = dismissGDPR;
        gdpr.removeAttribute('hidden');
    }
    function dismissGDPR() {
        $ga.enable();
        localStorage.setItem('gdpr', true);
        gdpr.setAttribute('hidden', true);
    }
}