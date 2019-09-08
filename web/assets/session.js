export async function getAccessToken() {
    let tokenResponse = await fetch(`${process.env.API_ORIGIN}/spotify/token`, {credentials: 'include'});
    let {token} = await tokenResponse.json();
    if (!token) {
        location.reload();
        return;
    }
    return token;
}

export function authorizeUserRedirect(rememberMe) {
    let qp = rememberMe ? '?permanent=true' : '';
    location.href = `${process.env.API_ORIGIN}/spotify/authorize${qp}`;
}

export async function isNewUser() {
    let response = await fetch(`${process.env.API_ORIGIN}/user/me`, {credentials: 'include'});
    return response.status === 401;
}

export function logout() {
    sessionStorage.clear();
    localStorage.clear();
    location.href = `${process.env.API_ORIGIN}/authenticate/logout`;
}
