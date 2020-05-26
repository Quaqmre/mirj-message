A lot of work is here,

implement user method
create client struct and methods +
create room struct and methods +
create hub struct and methods, +

this repo developing for command line chat

TODO : 
    - more readable logs should be write, all log template shoul be same => Done
    - command line flags implementing => done
    
    
USAGE:
    After Server up `go run main.go`
    Client can connect server `go run client/*.go`
    CLIEND COMMAND:
    &ls room => list all room and total count => done
    &ls user => list all user in current room and total user count => done
    &ch <string> => change user name after connected server => done
    &joın <string> => join spesfic room with room name if exist => done
    &mk <string> => create room if not exist => done
    &ext => exit room => done 

    ClIENT COMMAND LINE ARGUMENT:
    -help => get help document
    -host => set host adresses
    -name => set username 
    -pass => set password
