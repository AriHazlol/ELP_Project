export function creerPaquet() {
    let paquet = [];
    // Cartes chiffres classiques (1 à 12) 
    for (let i = 1; i <= 12; i++) {
        for (let j = 0; j < i; j++) {
            paquet.push({ valeur: i, type: 'chiffre' });
        }
    }
    // Multiplicateurs (x2) et Boucliers 
    for (let i = 0; i < 2; i++) {
        paquet.push({ valeur: 2, type: 'x2' });
        paquet.push({ valeur: 0, type: 'bouclier' });
    }
    // Cartes Bonus : +2, +4, +6, +8, +10 
    [2, 4, 6, 8, 10].forEach(val => {
        paquet.push({ valeur: val, type: 'bonus' });
    });
    // Carte Flip Three 
    for (let i = 0; i < 3; i++) {
        paquet.push({ valeur: 0, type: 'flip_three' });
    }

    return paquet;
}

export function melanger(paquet) {
    // Mélange de Fisher-Yates 
    for (let i = paquet.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [paquet[i], paquet[j]] = [paquet[j], paquet[i]];
    }
    return paquet;
}

export function calculerScoreFinal(main) {
    let base = 0;
    let multiplicateur = 1;
    main.forEach(c => {
        if (c.type === 'chiffre') base += c.valeur;
        if (c.type === 'bonus') base += c.valeur; 
        if (c.type === 'x2') multiplicateur *= 2;
    });
    return base * multiplicateur;
}

export function testerDoublon(main, carte) {
    if (carte.type !== 'chiffre') return false;
    return main.some(c => c.type === 'chiffre' && c.valeur === carte.valeur); 
}

export function aUnBouclier(main) {
    return main.some(c => c.type === 'bouclier'); 
}

export function consommerBouclier(main) {
    const index = main.findIndex(c => c.type === 'bouclier');
    if (index !== -1) main.splice(index, 1);
}