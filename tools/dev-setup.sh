#!/bin/bash

INSTALLDIR="/usr/local/lib"
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

function crun {
    "$@"
    local status=$?
    if [ $status -ne 0 ]; then
        echo "error with $1" >&2
    fi
    return $status
}

if [ "$1" == "uninstall" ]; then
  echo -e "${GREEN} Removing libsass and sassC...${NC}"
  rm "/usr/local/bin/sassc"
  #rm "/usr/local/bin/jsmin"
  rm -r "$INSTALLDIR/libsass"
  rm -r "$INSTALLDIR/sassc"
  #rm -r "$INSTALLDIR/jsmin"
  echo -e "${GREEN}Done. Bye!${NC}"
  exit 0
fi
if [[ $EUID -ne 0 ]]; then
  echo -e "${RED}You must be a root user, $(whoami)!${NC}" 2>&1
  exit 1
else
  mkdir "dist"
  echo -e "${GREEN}Cloning libsass and sassc...${NC}"
  rm -r "$INSTALLDIR/libsass"
  rm -r "$INSTALLDIR/sassc"
  crun git clone https://github.com/sass/libsass.git --branch 3.3.0 --depth 1 $INSTALLDIR/libsass
  if [ $? -ne 0 ]; then
    echo -e "${RED}Git clone of libsass failed${NC}"
    exit 1
  fi
  crun git clone https://github.com/sass/sassc.git --branch 3.3.0 --depth 1 $INSTALLDIR/sassc
  if [ $? -ne 0 ]; then
    echo -e "${RED}Git clone of sassc failed${NC}"
    exit 1
  fi
  echo -e "${GREEN}Done. Setting & flushing SASS_LIBSASS_PATH...${NC}"
  SASS_LIBSASS_PATH="$INSTALLDIR/libsass"
  export SASS_LIBSASS_PATH="$INSTALLDIR/libsass"
  #echo "SASS_LIBSASS_PATH=$SASS_LIBSASS_PATH" >> /etc/environment
  #source /etc/environment
  echo -e "${GREEN}SASS_LIBSASS_PATH is now set to $SASS_LIBSASS_PATH${NC}"

  echo -e "${GREEN}Building SassC...${NC}"
  cd "$INSTALLDIR/sassc/"
  crun make
  if [ $? -ne 0 ]; then
    echo -e "${RED}Build failed. See output log for more info, terminating${NC}"
    exit 1
  fi
  echo -e "${GREEN}Build successfully finished, creating symlink in /usr/local/bin ...${NC}"
  crun ln -s $INSTALLDIR/sassc/bin/sassc /usr/local/bin
  if [ $? -ne 0 ]; then
    echo -e "${RED}Symlinking failed!${NC}"
    exit 1
  fi
  echo -e "${GREEN}All done. SassC version:${NC}"
  sassc -v

#  echo -e "${GREEN}Cloning & installing jsmin.c ... ${NC}"
#  crun git clone https://github.com/douglascrockford/JSMin.git $INSTALLDIR/jsmin
#  if [ $? -ne 0 ]; then
#    echo -e "${RED}Cloning of jsmin failed!${NC}"
#    exit 1
#  fi
#  cd "$INSTALLDIR/jsmin"
#  crun gcc -o jsmin jsmin.c
#  if [ $? -ne 0 ]; then
#    echo -e "${RED}Compiling of jsmin failed!${NC}"
#    exit 1
#  fi
#  crun ln -s $INSTALLDIR/jsmin/jsmin /usr/local/bin/jsmin
#  if [ $? -ne 0 ]; then
#    echo -e "${RED}Symlinking of jsmin failed!${NC}"
#    exit 1
#  fi
  echo -e "${GREEN}Finished. Have a nice day/night! :)${NC}"
fi
