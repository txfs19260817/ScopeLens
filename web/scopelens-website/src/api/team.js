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

export const GetTeamsBySearchCriteria = (page, data) => {
    const headers = {
        'Content-Type': 'application/json',
    }
    return http.post("team/search", data, {
        headers: headers,
        params: {
            page: page
        }
    });
};