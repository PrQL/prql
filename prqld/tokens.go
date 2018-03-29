package main

import (
  "strings"
  "strconv"

  "github.com/prql/prql/util"
  log "github.com/sirupsen/logrus"
)


type TokenEntry struct {
  living bool

  user     string
  password string
  hostName string
  dbname   string

  origins []string
}

var (
  TokenPool = make(map[string]TokenEntry)
)


/**
* Token File Entry Schema
*
* token:username:password:hostName:dbname:origins:living
*
* token - 32-character string generated by the cli. Used to identify credentials 
*         and passed to the program from the client in the X-PrQL-Token header.
*
* username - The database username.
* 
* password - The database user's password.
* 
* hostName - The tag that was defined by the user while using the cli to store the credentials
*              of the database server.
*
* dbname - The name of the database to create the connection to.
*
* origins - A comma separated list of authorized origins that this entry's token can
*           be used with. If left empty, then all origins are authorized.
*
* living - A boolean which indicates whether this entry will spawn a lifelong connection
*          to the specified database whenever the system starts up.
*/

func PopulateTokenPool(refresh bool) {
  if refresh {
    TokenPool = make(map[string]TokenEntry) 
  }

  entries := util.ParseEntryFile("/var/lib/prql/tokens")

  for i, parts := range entries {
    if len(parts) != 7 {
      log.Error("Invalid token entry at line " + strconv.Itoa(i + 1)) 
      continue
    }

    origins := strings.Split(parts[5], ",")

    living, err := strconv.ParseBool(parts[6])
    if err != nil {
      living = false 
    }

    TokenPool[parts[0]] = TokenEntry {
      user: parts[1], 
      password: parts[2], 
      hostName: parts[3], 
      dbname: parts[4], 
      origins: origins, 
      living: living,
    }
  }
}
