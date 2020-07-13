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

export const getTeams = (page) => {
    return http.get("/team/teams", {
        params: {
            page: page
        }
    })
};