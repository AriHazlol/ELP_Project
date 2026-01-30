import * as readline from 'node:readline/promises';
import { stdin as input, stdout as output } from 'node:process';
import * as Carte from './carte.js'; 
import * as Affichage from './affichage.js'; 
import { chargerScore, sauverRecord } from './sauvegarde.js';

const rl = readline.createInterface({ input, output });

function formatCarte(carte) {
    if (carte.type === 'chiffre') return carte.valeur;
    if (carte.type === 'bonus') return `BONUS +${carte.valeur}`;
    return carte.type.toUpperCase();
}

async function piocherAction(nom, main, paquet) {
    if (paquet.length === 0) {
        Affichage.alerte("Le paquet est vide !");
        return { perdu: false, fini: true };
    }

    const carte = paquet.pop();
    Affichage.info(`${nom} pioche : ${formatCarte(carte)}`);

    if (Carte.testerDoublon(main, carte)) {
        if (Carte.aUnBouclier(main)) {
            Affichage.alerte("DOUBLON ! Le BOUCLIER protège " + nom);
            Carte.consommerBouclier(main);
        } else {
            Affichage.alerte(`DOUBLON ! ${nom} perd ses points sur cette manche.`);
            return { perdu: true, fini: true };
        }
    }

    main.push(carte);

    if (carte.type === 'flip_three') {
        Affichage.alerte("EFFET FLIP THREE ! Piochez 3 cartes obligatoirement !");
        for (let i = 1; i <= 3; i++) {
            Affichage.info(`Pioche forcée ${i}/3...`);
            const res = await piocherAction(nom, main, paquet);
            if (res.perdu) return res;
        }
    }

    return { perdu: false, fini: false };
}

async function main() {
    const record = chargerScore();
    Affichage.messageBienvenue(record);

    // Choix du nombre de joueurs
    let nbInput = await rl.question("Combien de joueurs participent (1-5) ? ");
    let nbJoueurs = parseInt(nbInput);
    if (isNaN(nbJoueurs) || nbJoueurs < 1 || nbJoueurs > 5) {
        Affichage.alerte("Nombre invalide, on part sur 2 joueurs par défaut.");
        nbJoueurs = 2;
    }

    let scoresTotaux = new Array(nbJoueurs).fill(0);
    const paquet = Carte.melanger(Carte.creerPaquet());

    for (let m = 1; m <= 3; m++) {
        Affichage.info(`\n========== MANCHE ${m} / 3 ==========`);
        
        let mains = Array.from({ length: nbJoueurs }, () => []);
        let actifs = new Array(nbJoueurs).fill(true);

        // alternanace des joueurs
        while (actifs.some(a => a)) {
            for (let i = 0; i < nbJoueurs; i++) {
                if (!actifs[i]) continue;

                const nom = `Joueur ${i + 1}`;
                Affichage.afficherEntete(m, nom);
                Affichage.afficherMain(mains[i]);

                const choix = await rl.question(`${nom} : Piocher (p) ou Arrêter (a) ? `);
                
                if (choix.toLowerCase() === 'p') {
                    const resultat = await piocherAction(nom, mains[i], paquet);
                    if (resultat.perdu) { 
                        mains[i] = []; 
                        actifs[i] = false; 
                    }
                } else {
                    actifs[i] = false;
                }
            }
        }

        for (let i = 0; i < nbJoueurs; i++) {
            const gain = Carte.calculerScoreFinal(mains[i]);
            scoresTotaux[i] += gain;
            Affichage.succes(`Fin de manche J${i+1} : +${gain} pts`);
        }
    }

    Affichage.afficherResultatFinal(scoresTotaux);

    if (sauverRecord(Math.max(...scoresTotaux))) {
        Affichage.info("NOUVEAU RECORD ENREGISTRÉ !");
    }
    rl.close();
}

main();