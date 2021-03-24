import http from "./index";

export const getSpeciesByImage = (file) => {
    const headers = {
        'Content-Type': 'multipart/form-data',
    }
    let formData = new FormData();
    formData.append('file', file);
    return http.post(process.env.VUE_APP_ADVANCED_URL+"predict", formData, {
        params:{
            type: 'team'
        }
    })
}