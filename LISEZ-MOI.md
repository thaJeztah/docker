Docker: le moteur de conteneurs [![Release](https://img.shields.io/github/release/docker/docker.svg)](https://github.com/docker/docker/releases/latest)
============================

Docker est un projet open source pour emballer, navire et exécuter une application
comme un conteneur léger.

Docker conteneurs sont à la fois indépendante du matériel * * et * plate-forme agnostique *.
Cela signifie qu'ils peuvent fonctionner n'importe où, à partir de votre ordinateur portable pour le plus grand
nuage calculer instance et tout le reste - et ils ne nécessitent pas
vous d'utiliser un système de langue, le cadre ou l'emballage particulier. Que
les grands blocs de construction pour le déploiement et la mise à l'échelle des applications web fait,
bases de données, et des services de back-end sans dépendre d'un empilement particulier
ou le fournisseur.

Docker a commencé comme une implémentation open-source du moteur de déploiement qui
pouvoirs [dotCloud](https://www.dotcloud.com), une plate-forme-as-a-Service populaire.
Il bénéficie directement à partir de l'expérience accumulée sur plusieurs années
d'opération à grande échelle et le soutien de centaines de milliers de
applications et bases de données.

![](docs/static_files/docker-logo-compressed.png "Docker")

## Divulgation de sécurité

La sécurité est très importante pour nous. Si vous avez une question concernant la sécurité,
s'il vous plaît communiquer les informations de manière responsable en envoyant un courriel à
security@docker.com et non en créant une question github.

## Mieux que VM

Une méthode commune pour la distribution d'applications et leur sandboxing
l'exécution est d'utiliser des machines virtuelles ou des machines virtuelles. VM formats typiques sont
La VMDK de VMware, la VDI VirtualBox d'Oracle, et ami d'Amazon EC2. En théorie
ces formats devraient permettre à chaque développeur de conditionner automatiquement
leur application dans une "machine" pour la distribution et le déploiement facile.
En pratique, cela se passe presque jamais, pour plusieurs raisons:

  * * * Taille: VM sont très grandes qui les rend peu pratique à stocker
     et transférer.
  * * * Performance: machines virtuelles en cours d'exécution consomme CPU et de mémoire importante,
    ce qui les rend peu pratique dans de nombreux scénarios, par exemple locale
    développement d'applications multi-niveaux, et de déploiement à grande échelle
    processeur et la mémoire-intensive des applications sur un grand nombre de
    des machines.
  * * * Portabilité: concurrence environnements VM ne jouent pas bien avec chaque
     autre. Bien que les outils de conversion existent, elles sont limitées et
     ajouter encore plus de frais généraux.
  * * * Hardware-centric: machines virtuelles ont été conçus avec des opérateurs de machines dans
    esprit, pas les développeurs de logiciels. En conséquence, ils offrent très
    outillage limité pour ce que les développeurs ont le plus besoin: la construction, les tests
    et la gestion de leur logiciel. Par exemple, les machines virtuelles offrent pas d'installations
    pour l'application des versions, la surveillance, la configuration, l'exploitation forestière ou
    découverte de service.

En revanche, Docker repose sur une méthode connue sous le nom de bac à sable différent
* * conteneurisation. Contrairement à la virtualisation traditionnelle, la conteneurisation
a lieu au niveau du noyau. La plupart des noyaux de systèmes d'exploitation modernes
désormais en charge les primitives nécessaires à la conteneurisation, y compris
Linux avec [openvz](https://openvz.org),
[vserver](http://linux-vserver.org) et, plus récemment,
[lxc](https://linuxcontainers.org/), avec Solaris
[zones](https://docs.oracle.com/cd/E26502_01/html/E29024/preface-1.html#scrolltoc),
et FreeBSD avec
[Prisons](https://www.freebsd.org/doc/handbook/jails.html).

Docker construit au-dessus de ces primitives de bas niveau pour offrir aux développeurs une
format portable et de l'environnement d'exécution qui permet de résoudre tous les quatre problèmes.
Docker conteneurs sont de petite taille (et leur transfert peuvent être optimisés avec
couches), ils ont fondamentalement mémoire et CPU-dessus de zéro, ils sont
complètement portable, et sont conçus à partir du sol avec un
centrée sur les applications de conception.

Peut-être le meilleur de tous, parce Docker fonctionne au niveau du système d'exploitation, il peut encore être
courir à l'intérieur d'une VM!

## Joue bien avec les autres

Docker ne vous oblige pas à acheter dans une programmation particulier
la langue, le cadre, le système d'emballage, ou la langue de configuration.

Votre application est un processus Unix? Est-il utiliser les fichiers, les connexions TCP,
variables d'environnement, les ruisseaux Unix standard et des arguments de ligne de commande
que les entrées et sorties? Puis Docker peut exécuter.

La construction de votre application peut être exprimée comme une séquence de tels
commandes? Puis Docker peut construire.

## Évasion dépendance enfer

Un problème commun pour les développeurs est la difficulté de gérer tous
les dépendances de leurs applications de manière simple et automatisée.

Ceci est habituellement difficile pour plusieurs raisons:

  * * Dépendances multiplates-formes *. Les applications modernes dépendent souvent
    une combinaison de bibliothèques et les binaires système, la langue-spécifiques
    paquets, des modules spécifiques cadres, les composants internes
    développé pour un autre projet, etc. Ces dépendances vivent dans
    différents «mondes» et nécessitent différents outils - ces outils
    généralement ne fonctionnent pas bien avec l'autre, nécessitant maladroite
    intégrations personnalisées.

  * * * Contradictoires dépendances. Différentes applications peuvent dépendre
    des versions différentes de la même dépendance. Manipuler les outils d'emballage
    ces situations avec divers degrés de facilité - mais ils ont tous
    les traiter de manières différentes et incompatibles, ce qui oblige à nouveau
    le développeur de faire un travail supplémentaire.

  * * * Dépendances personnalisés. Un développeur peut avoir besoin pour préparer une coutume
    version de la dépendance de leur application. Certains systèmes d'emballage
    peut gérer des versions personnalisées d'une dépendance, d'autres ne peuvent pas - et tout
    de les traiter différemment.


Docker résout le problème de la dépendance enfer en donnant le développeur d'un simple
façon d'exprimer * tous * les dépendances de leur application dans un seul endroit, alors que
la rationalisation du processus d'assemblage. Si cela vous fait penser
[XKCD 927](https://xkcd.com/927/), ne vous inquiétez pas. Docker ne
* * remplacer vos systèmes d'emballage préférées. Il orchestre tout simplement
leur utilisation d'une manière simple et reproductible. Comment fait-elle cela? Avec
couches.

Docker définit une construction que l'exécution d'une séquence de commandes Unix, un
après l'autre, dans le même récipient. Construire commandes modifient le
contenu du récipient (généralement par l'installation de nouveaux fichiers sur la
système de fichiers), la commande suivante modifie un peu plus, etc. Comme chaque
commande construire hérite le résultat des commandes précédentes, le
* commande * dans laquelle les commandes sont exécutées exprime * * dépendances.

Voici un processus typique Docker de construction:

`` `bash
À partir d'Ubuntu: 12.04
RUN apt-get jour && apt-get installer python -y python-pip boucle
RUN courbent -SSL https://github.com/shykes/helloflask/archive/master.tar.gz | goudron -xzv
RUN cd helloflask maître && PIP installer requirements.txt -r
`` `

Notez que Docker ne se soucie pas comment les dépendances * * sont construit - aussi longtemps
car ils peuvent être construits en exécutant une commande Unix dans un récipient.


Commencer
===============

Docker peut être installé sur votre ordinateur pour créer des applications ou
sur les serveurs pour les exécuter. Pour commencer, [vérifier l'installation
instructions dans le
Documentation](https://docs.docker.com/engine/installation/).

Nous offrons également un [tutoriel interactif](https://www.docker.com/tryit/)
pour apprendre rapidement les bases de l'utilisation Docker.

Exemples d'utilisation
==============

Docker peut être utilisé pour exécuter des commandes de courte durée, les démons de longue durée
(serveurs d'application, bases de données, etc.), des séances de shell interactif, etc.

Vous pouvez trouver une [liste de monde réel
Exemples]() dans le https://docs.docker.com/engine/examples/
Documentation.

Sous la capuche
--------------

Sous le capot, Docker est construite sur les composants suivants:

* Le
  [cgroups](https://www.kernel.org/doc/Documentation/cgroups/cgroups.txt)
  et
  [categories](http://man7.org/linux/man-pages/man7/namespaces.7.html)
  fonctionnalités du noyau Linux
* Le [Go](https://golang.org) langage de programmation
* Le [Docker Spécifications de l'image](https://github.com/docker/docker/blob/master/image/spec/v1.md)
* Le [Libcontainer Spécifications](https://github.com/opencontainers/runc/blob/master/libcontainer/SPEC.md)

Contribuer à Docker [![GoDoc](https://godoc.org/github.com/docker/docker?status.svg)](https://godoc.org/github.com/docker/docker)
======================

| ** ** Maître (Linux) | ** ** Expérimentale (linux) | ** ** De Windows | ** FreeBSD ** |
| ------------------ | ---------------------- | ------- - | --------- |
| [! [Jenkins Créer Status](https://jenkins.dockerproject.org/view/Docker/job/Docker%20Master/badge/icon)](https://jenkins.dockerproject.org/view/Docker/job/Docker%20Master/) | [! [Jenkins Créer Status](https://jenkins.dockerproject.org/view/Docker/job/Docker%20Master%20%28experimental%29/badge/icon)](https://jenkins.dockerproject.org/view/Docker/job/Docker%20Master%20%28experimental%29/) | [![Construire Status](http://jenkins.dockerproject.org/job/Docker%20Master%20(windows)/badge/icon)](http://jenkins.dockerproject.org/job/Docker%20Master%20(windows)/) | [![Construire Status](http://jenkins.dockerproject.org/job/Docker%20Master%20(freebsd)/badge/icon)](http://jenkins.dockerproject.org/job/Docker%20Master%20(freebsd)/) |

Vous voulez pirater sur Docker? Impressionnant! Nous avons [instructions pour vous aider
commencé à contribuer code ou la documentation](https://docs.docker.com/opensource/project/who-written-for/).

Ces instructions ne sont probablement pas parfait, s'il vous plaît laissez-nous savoir si quelque chose
sent erronée ou incomplète. Mieux encore, soumettre une PR et de les améliorer vous-même.

Obtenir le développement construit
==============================

Vous voulez exécuter Docker d'un maître construction? Vous pouvez télécharger
maître construit au [master.dockerproject.org](https://master.dockerproject.org).
Ils sont mis à jour à chaque commit fusionné dans la branche master.

Je ne sais pas comment utiliser cette nouvelle fonctionnalité super cool dans le maître construction? Chèque
les docs de maître à
[docs.master.dockerproject.org](http://docs.master.dockerproject.org).

Comment le projet est exécuté
======================

Docker est un projet très, très actif. Si vous voulez en savoir plus sur la façon dont il est géré,
ou si vous voulez vous impliquer davantage, le meilleur endroit pour commencer est [le répertoire du projet](https://github.com/docker/docker/tree/master/project).

Nous sommes toujours ouverts aux suggestions sur les améliorations de processus, et nous sommes toujours à la recherche pour plus de mainteneurs.

### Parler à d'autres utilisateurs et contributeurs Docker

<table class = "tg">
  <largeur de col = "45%">
  <largeur de col = "65%">
  <tr>
    <td> Internet & nbsp; Relais & nbsp; & chat nbsp; (IRC) </ td>
    <td>
      <p>
        IRC Une ligne directe à nos utilisateurs les plus avertis Docker; nous avons
        à la fois le <code> #docker </ code> et <code> # docker-dev </ code> sur le groupe
        <strong> irc.freenode.net </ strong>.
        IRC est une riche protocole de chat mais il peut submerger les nouveaux utilisateurs. Vous pouvez rechercher
        <a href="https://botbot.me/freenode/docker/#" target="_blank"> notre Archives Chat </a>.
      </ p>
      Lisez notre <a href="https://docs.docker.com/project/get-help/#irc-quickstart" target="_blank"> IRC Guide de démarrage </a> un moyen facile pour commencer.
    </ td>
  </ tr>
  <tr>
    <td> Google Groupes </ td>
    <td>
      Il existe deux groupes.
      <a href="https://groups.google.com/forum/#!forum/docker-user" target="_blank"> </a> Docker utilisateur
      est pour les personnes utilisant des conteneurs Docker.
      Le <a href="https://groups.google.com/forum/#!forum/docker-dev" target="_blank"> docker-dev </a>
      groupe est pour les contributeurs et autres personnes qui contribuent à la Docker
      projet.
    </ td>
  </ tr>
  <tr>
    <td> Twitter </ td>
    <td>
      Vous pouvez suivre <a href="https://twitter.com/docker/" target="_blank"> Docker Twitter RSS </a>
      pour obtenir des mises à jour sur nos produits. Vous pouvez également nous Tweet Questions ou tout simplement
      blogs d'actions ou des histoires.
    </ td>
  </ tr>
  <tr>
    <td> Stack Overflow </ td>
    <td>
      Stack Overflow a plus de 7000 des questions Docker cotées. Nous régulièrement
      surveiller <a href="https://stackoverflow.com/search?tab=newest&q=docker" target="_blank"> des questions Docker </a>
      et ainsi de faire beaucoup d'autres utilisateurs avertis Docker.
    </ td>
  </ tr>
</ table>

### Juridique

* Offert à vous courtoisie de notre conseiller juridique. Pour plus de contexte,
s'il vous plaît voir le document [AVIS](https://github.com/docker/docker/blob/master/NOTICE) dans cette pension. *

Utilisation et le transfert des Docker peut être soumis à certaines restrictions imposées par la
États-Unis et d'autres gouvernements.

Il est de votre responsabilité de vous assurer que votre utilisation et / ou le transfert ne
violer les lois applicables.

Pour plus d'informations, s'il vous plaît voir https://www.bis.doc.gov


Licences
=========
Docker est sous licence Apache License, Version 2.0. Voir
[LICENCE](https://github.com/docker/docker/blob/master/LICENSE) pour le plein
texte de la licence.

Autres projets connexes Docker
=============================
Il ya un certain nombre de projets en cours de développement qui sont basées sur de Docker
technologie de base. Ces projets élargir l'outillage construit autour de la
Docker plateforme pour élargir son application et l'utilité.

* [Registre Docker](https://github.com/docker/distribution): Registre
serveur pour Docker (hébergement / livraison des référentiels et des images)
* [Docker machine](https://github.com/docker/machine): la gestion de la machine
pour un monde de conteneurs-centric
* [Docker Swarm](https://github.com/docker/swarm): Un regroupement Docker-natale
système
* [Docker Compose](https://github.com/docker/compose) (anciennement la figure):
Définir et exécuter des applications multi-conteneurs
* [Kitematic](https://github.com/docker/kitematic): Le moyen le plus facile à utiliser
Docker sur Mac et Windows

Si vous connaissez un autre projet en cours de qui devrait être listé ici, s'il vous plaît aider
nous tenir cette liste à jour en soumettant une PR.

Impressionnant-Docker
==============
Vous pouvez trouver plus de projets, des outils et des articles liés à Docker dans la [Liste awesome-docker](https://github.com/veggiemonk/awesome-docker). Ajouter votre projet il.
