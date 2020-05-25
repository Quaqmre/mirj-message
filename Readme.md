A lot of work is here,

implement user method
create client struct and methods +
create room struct and methods +
create hub struct and methods,

this repo developing for command line chat

TODO : 
    - more readable logs should be write, all log template shoul be same
    - command line flags implementing
    - 
USAGE:

    CLIEND COMMAND:
    &ls room => list all room and total count
    &ls user => list all user in current room and total user count
    &ch <string> => change user name after connected server
    &join <string> => join spesfic room with room name if exist
    &crete <string> => create room if not exist
    &exit => exit room 

    SERVER COMMAND LINE ARGUMENT:
    -help => get help document
    -host => set host adresses
    -name => set username 
    -pass => set password