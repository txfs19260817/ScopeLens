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

export const getTeamsByLikes = (page) => {
    return http.get("/team/likes", {
        params: {
            page: page
        }
    })
};

export const getTeamByID = (id) => {
    return http.get("/team/teams/" + id)
};

export const getTeamsBySearchCriteria = (page, data) => {
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

export const getPokemonUsageByFormat = (format) => {
    return http.get("/team/usage/" + format)
}

export const insertLikeByUsername = (username, id, token) => {
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': token
    }

    const data = {
        username: username,
        id: id
    }

    return http.post("user/like", data, {
        headers: headers
    });
}

export const getLikedTeamsByUsername = (page, username) => {
    return http.get("/team/likes/" + username, {
        params: {
            page: page
        }
    })
}

export const getUploadedTeamsByUsername = (page, username) => {
    return http.get("/team/uploaded/" + username, {
        params: {
            page: page
        }
    })
}