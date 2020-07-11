import de from "./pokemonNames/de.json"
import en from "./pokemonNames/en.json"
import fr from "./pokemonNames/fr.json"
import ja from "./pokemonNames/ja.json"
import ko from "./pokemonNames/ko.json"
import zhHans from "./pokemonNames/zh-hans.json"
import zhHant from "./pokemonNames/zh-hant.json"


export const pmNames4Select = [
    en,
    zhHans,
    ja
].reduce((r, a) =>
    a.map((v, i) => (r[i] || []).concat(v)), []
).map(e => e.join('/'));