1. make a bot on the discord application portal  
you must enable the message content intent on your bot's application page:  
![](https://i.imgur.com/qL2etXv.png)
2. when giving permission to the bot you must give at least the following permissions:  
`Read Messages/View Channels, Send Messages, Manage Messages, Embed Links, Add Reactions`
3. clone the repository  
`git clone https://github.com/Roguezilla/starboard-go.git`  
also run  
`git config core.fileMode false`  
4. install reflex  
4.1 general  
`go install github.com/cespare/reflex@latest`  
note that if you use this method, you likely need to add `:$HOME/go/bin` to your `$PATH` export in either `$HOME/.profile` or `/etc/profile`(depending where you initially put it), the end result being:  
`export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin`    
4.2. ubuntu  
`sudo apt get install reflex`  
5. after installing the requirements run  
`go run .`  
![](https://i.imgur.com/hvOfUzT.png)
6. now `CTRL+C` and run:  
`chmod 777 start.sh`  
`./start.sh` (use this command from now one to run the bot)  
this will automatically restart the bot when it gets updated:  
![](https://i.imgur.com/FkGXoSU.png)
7. all that's left is to setup the bot for your server:  
`sb!setup <archive_channel> <archive_emote> <archive_emote_amount>`  
![](https://i.imgur.com/ex6q23f.png)  
