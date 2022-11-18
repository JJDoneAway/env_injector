#/bin/zsh

export ODJ_DEP_BACKEND4ADDRESSBOOK_URI=https://johannes-address-book-01.dev.sit.sys.odj.cloud
export ODJ_DEP_BACKEND4ADDRESSBOOK_SERVICE="Some strange foo like $%&$§"

cp -rf dist_original/* dist
clear
echo "####################"
echo "run it for the first time to see the replacements"
echo "####################"
go run github.com/JJDoneAway/env_injector ./dist ./.env.odj

echo "####################"
echo "run it for the second time to check if every thing hase been replaced"
echo "####################"
go run github.com/JJDoneAway/env_injector ./dist ./.env.odj

unset ODJ_DEP_BACKEND4ADDRESSBOOK_URI
unset ODJ_DEP_BACKEND4ADDRESSBOOK_SERVICE