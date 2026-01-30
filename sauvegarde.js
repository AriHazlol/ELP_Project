import fs from 'node:fs';

export function chargerScore() {
    if (!fs.existsSync('save.json')) return 0;
    const raw = fs.readFileSync('save.json');
    return JSON.parse(raw).record || 0;
}

export function sauverRecord(nouveau) {
    const actuel = chargerScore();
    if (nouveau > actuel) {
        fs.writeFileSync('save.json', JSON.stringify({ record: nouveau }));
        return true;
    }
    return false;
}