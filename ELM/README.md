# Projet TcTurtle in Elm by G26
## 1. Présentation du Projet TcTurtle
TcTurtle Elm Edition est un interpréteur de commandes graphiques. Ce projet a été développé pour explorer les fondamentaux de la programmation fonctionnelle à travers le langage Elm.

L'objectif est simple : transformer des instructions textuelles en mouvements géométriques précis exécutés par une "tortue" virtuelle. L'application repose sur trois piliers techniques :

- Interprétation : Analyser les commandes saisies par l'utilisateur pour les traduire en données compréhensibles par le programme.
- Moteur de Mouvement : Un système de calcul basé sur la trigonométrie qui gère la position (x, y) et l'orientation de la tortue.
- Rendu SVG : Une visualisation fluide et vectorielle qui dessine le tracé en temps réel dans le navigateur.

Ce projet illustre la robustesse de l'Architecture Elm, garantissant une application sans erreurs à l'exécution et une gestion fluide de l'état du dessin. Que ce soit pour tracer un simple carré ou des fractales complexes via des boucles, TcTurtle fait le pont entre mathématiques et création numérique.
## 2. Cheatsheet des commandes

| Commandes | Description | Valeur par défaut |
| --- | --- | --- |
| fd [valeur] | Permet à la tortue d'avancer de [valeur] pixels tout en traçant un trait | 50 |
| rt [valeur] | Fait tourner la tortue de [valeur] degrés sur la droite | 90 |
| lt [valeur] | Fait tourner la tortue de [valeur] degrés sur la gauche | 90 |
| cercle [valeur] | Forme un cercle de rayon [valeur] | 50 |
| carre [valeur] | Forme un carré de coté [valeur] | 50 |
| etoile [valeur] | Forme une étoile de longueur [valeur] | 50 |
| clear | Efface les traits sans réinitialiser la position de la tortue | / |
| reset | Efface les traits et réinitialise la position de la tortue | / |

## 3. Pistes d'amélioration
Pour ce qui est des pistes qui restent à explorer ; Il serait plus agréable pour l'utilisateur de pouvoir utiliser sa touche entrée afin de rentrer une commande directement et que le texte disparaisse dès lors que la commande s'éxecute. De Plus de ombreuses formes et fonctionnalités peuvent se voir ajoutées au projet (ex : d'autres polygones, pouvoir changer la coueleur du pinceau, remplir le canva totalement d'une couleur voir même pouvoir changer de tortue).
