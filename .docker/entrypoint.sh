#!/bin/bash

#start app after wait the database be ready
sh -c "/wait && ./myapp rest"