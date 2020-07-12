/**
 * @return {string}
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
 * @return {string}
 */
export function ProcessStr(s) {
    return s.toString().replace(/[^a-z0-9]/gi, '').toLowerCase();
}

export function SortObjectArrayByValue(arr) {
    return arr.sort((a, b) => (b.value - a.value));
}

// https://stackoverflow.com/questions/36280818/how-to-convert-file-to-base64-in-javascript
export const toBase64 = file => new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = error => reject(error);
});