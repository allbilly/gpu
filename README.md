# gpu

# add new subtree
git remote add <alias> <external-repo-URL>
git subtree add -P <folder> <alias> <branch> --squash

git remote add turingas https://github.com/daadaada/turingas
git subtree add -P turingas turingas --squash
