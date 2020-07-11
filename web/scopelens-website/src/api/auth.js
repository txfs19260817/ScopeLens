import http from "./index";

export const registerRequest = (data) => {
    return http.post("auth/register", data);
}

export const loginRequest = (data) => {
    return http.post("auth/login", data);
}