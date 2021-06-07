import http from "./index";

const advapiUrl = process.env.VUE_APP_ADVANCED_URL + "predict";

export const getSpeciesByImage = (file) => {
    const headers = {
        'Content-Type': 'multipart/form-data',
    }
    let formData = new FormData();
    formData.append('file', file);
    // a baseURL overriding
    return http({ method: 'post', url: advapiUrl, baseURL: '/', data: formData, params: { type: 'team' } })
}