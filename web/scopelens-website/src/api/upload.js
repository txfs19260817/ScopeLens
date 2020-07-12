import http from "./index";

export const insertTeam = (data, token) => {
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': token
    }
    return http.post("team/post", data, {
        headers: headers
    });
}