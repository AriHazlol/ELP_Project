const RESET = "\x1b[0m";
const ROUGE = "\x1b[31m";
const VERT = "\x1b[32m";
const JAUNE = "\x1b[33m";
const BLEU = "\x1b[34m";

export function messageBienvenue(record) {
    console.log(`\nBIENVENUE AU FLIP 7`);
    console.log(`Record actuel : ${JAUNE}${record} pts${RESET}\n`);
}

export function afficherEntete(manche, nom) {
    console.log(`\n--- MANCHE ${manche} | TOUR DE : ${JAUNE}${nom}${RESET} ---`);
}

export function afficherMain(main) {
    const rendu = main.map(c => {
        if (c.type === 'chiffre') return c.valeur;
        if (c.type === 'bonus') return `[+${c.valeur}]`;
        return `[${c.type.toUpperCase()}]`;
    }).join(" | ");
    console.log(`Main : ${rendu}${RESET}`);
}

export function afficherResultatFinal(scores) {
    console.log(`\n--- SCORE FINAL ---`);
    let maxScore = -1;
    let gagnant = "";

    scores.forEach((score, index) => {
        console.log(`Joueur ${index + 1} : ${score} pts`);
        if (score > maxScore) {
            maxScore = score;
            gagnant = `Joueur ${index + 1}`;
        } else if (score === maxScore) {
            gagnant = "Égalité";
        }
    });

    console.log(`${VERT}RÉSULTAT : ${gagnant.toUpperCase()} !${RESET}`);
}

export function alerte(message) {
    console.log(`${ROUGE}!! ${message} !!${RESET}`);
}

export function succes(message) {
    console.log(`${VERT}-> ${message}${RESET}`);
}

export function info(message) {
    console.log(`${BLEU}${message}${RESET}`);
}