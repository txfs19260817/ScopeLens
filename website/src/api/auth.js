import http from "./index";

export const registerRequest = (data, recaptcha) => {
    return http.post("auth/register?recaptcha="+recaptcha, data);
}

export const loginRequest = (data) => {
    return http.post("auth/login", data);
}

export const checkToken = (token) => {
    const headers = {
        'Authorization': token
    }
    return http.get("auth/checktoken", {
        headers: headers
    })
}