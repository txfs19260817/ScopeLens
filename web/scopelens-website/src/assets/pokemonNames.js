import de from "./pokemonNames/de.json"
import en from "./pokemonNames/en.json"
import fr from "./pokemonNames/fr.json"
import ja from "./pokemonNames/ja.json"
import ko from "./pokemonNames/ko.json"
import zhHans from "./pokemonNames/zh-hans.json"
import zhHant from "./pokemonNames/zh-hant.json"
import {ProcessStr} from "./utils";

const group = {
    0: 1,
    151: 2,
    251: 3,
    386: 4,
    493: 5,
    649: 6,
    721: 7,
    809: 8
}
let currentGroup = 0

const pmNames = [
    en,
    zhHans,
    ja
].reduce((r, a) =>
    a.map((v, i) => (r[i] || []).concat(v)), []
).map(e => e.join('/'));

let pmNames4Select = []

for (let i = 0; i < pmNames.length; i++) {
    if (group[i] !== undefined) {
        if (i !== 0) {
            pmNames4Select.push({divider: true})
        }
        currentGroup = group[i]
        pmNames4Select.push({header: 'Generation ' + currentGroup})
    }
    pmNames4Select.push({
        name: pmNames[i],
        group: 'Generation ' + currentGroup,
        avatar: ProcessStr(pmNames[i]) + '.png'
    })
}

export {pmNames4Select}