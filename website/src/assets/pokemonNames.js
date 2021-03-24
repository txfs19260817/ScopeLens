import de from "./pokemonNames/de.json"
import en from "./pokemonNames/en.json"
import fr from "./pokemonNames/fr.json"
import ja from "./pokemonNames/ja.json"
import ko from "./pokemonNames/ko.json"
import zhHans from "./pokemonNames/zh-hans.json"
import zhHant from "./pokemonNames/zh-hant.json"
import forme from "./pokemonNames/forme.json"
import {ProcessStr} from "./utils";

// Pokemon are grouped by gen
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

// "en/zh-hans/ja"
const pmNames = [
    en,
    zhHans,
    ja
].reduce((r, a) =>
    a.map((v, i) => (r[i] || []).concat(v)), []
).map(e => e.join('/'));


let pmNames4Select = [] // for v-autocomplete.items
let en2pmNames = {} // processed English names to pmNames
for (let i = 0; i < pmNames.length; i++) {
    // add dividers
    if (group[i] !== undefined) {
        if (i !== 0) {
            pmNames4Select.push({divider: true})
        }
        currentGroup = group[i]
        pmNames4Select.push({header: 'Generation ' + currentGroup})
    }
    // add pmNames
    let processedName = ProcessStr(en[i])
    pmNames4Select.push({
        name: pmNames[i],
        group: 'Generation ' + currentGroup,
        avatar: processedName + '.png'
    })
    en2pmNames[processedName] = pmNames[i]
}

export {pmNames4Select, en2pmNames, forme}