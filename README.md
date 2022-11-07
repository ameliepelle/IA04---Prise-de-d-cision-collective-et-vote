# IA04---Prise-de-d-cision-collective-et-vote

## Récupération du projet
Le projet est disponible sur github
### Pour récupérer le code :

go mod init IA04---Prise-de-d-cision-collective-et-vote
go get github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote@latest
go get github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/comsoc@v0.0.0-20221107185036-5f177044450a

### Pour le lancer :
go run github.com/ameliepelle/IA04---Prise-de-d-cision-collective-et-vote/agt/restagentdemo/cmd/launch-all-rest-agents

## Ce qui a été implémenté :

### Méthodes de votes :
* Borda
* MajoritéSimple
* Approval
* Condorcet
* Kemeny

### Commandes: 
* /new_Ballot
* /vote
* /result

### Autres :
* Fonction de tiebreak et tiebreakFactory
