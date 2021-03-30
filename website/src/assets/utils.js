/**
 * Convert timestamp to string
 * @param timestamp
 * @returns {string}
 * @constructor
 */
export function DateConversion(timestamp) {
    let date = new Date(Date.parse(timestamp));
    let Y = date.getFullYear() + '-';
    let M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + '-';
    let D = (date.getDate() < 10 ? '0' + date.getDate() : date.getDate()) + ' ';
    let h = (date.getHours() < 10 ? '0' + date.getHours() : date.getHours()) + ':';
    let m = (date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()) + ':';
    let s = (date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds());
    return Y + M + D + h + m + s
}

/**
 * Remove all non-alphanumeric characters and lower every characters
 * @param {string} s
 * @returns {string}
 * @constructor
 */
export function ProcessStr(s) {
    return s.toString().replace(/[^a-z0-9]/gi, '').toLowerCase();
}

export function SortObjectArrayByValue(arr) {
    return arr.sort((a, b) => (b.value - a.value));
}

/**
 * Convert a file to Base64 encoded string
 * Ref: https://stackoverflow.com/questions/36280818/how-to-convert-file-to-base64-in-javascript
 * @param file
 * @returns {Promise<unknown>}
 */
export const toBase64 = file => new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = error => reject(error);
});

/**
 * Redirect To CloudFront accelerated resources
 * @param {string} src
 * @returns {string}
 * @constructor
 */
export function RedirectToCloudFront(src) {
    if (process.env.VUE_APP_CDN_URL.length > 0) {
        return src.toString().replace(`https://scopelens.s3.ap-northeast-1.amazonaws.com`, process.env.VUE_APP_CDN_URL)
    }
    return src
}