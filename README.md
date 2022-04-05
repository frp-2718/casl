# casl

Outil de comparaison des localisations de documents physiques dans Alma et dans le SUDOC à partir d'une liste de PPN.

## Utilisation

    ./casl fichier_ppn...

### Configuration

*fichier_ppn* contient un PPN par ligne.

Nécessite dans le répertoire de l'exécutable un fichier _config.json_ contenant
:
- le chemin vers le fichier de correspondance _alma-rcr.csv_
- la clé d'API Alma
- la liste des ILN concernés
- la liste des collections Alma ignorées, éventuellement vide
- la liste des RCR ignorés, éventuellement vide

_alma-rcr.csv_ établit la correspondance entre les bibliothèques Alma et les RCR du SUDOC. Format : `intitulé_alma,code_bib_alma,RCR,ILN`

### Résultat

Un fichier _resultats_XXXXXXX.csv_ contenant les colonnes suivantes :
1. PPN fautif
2. ILN concerné
3. Intitulé de la bibliothèque concernée dans Alma (seulement si le PPN est présent dans Alma)
4. Intitulé de la bibliothèque concernée dans le SUDOC (seulement si le PPN est présent dans Alma)
5. RCR concerné

Comme on recherche les anomalies, les colonnes 3 et 4 ne peuvent pas contenir
la même valeur en même temps :
- si la colonne 3 contient une valeur, alors le PPN existe dans Alma mais pas dans le SUDOC (et la colonne 4 est vide)
- si la colonne 4 contient une valeur, alors le PPN existe dans le SUDOC mais pas dans Alma (et la colonne 3 est vide)

