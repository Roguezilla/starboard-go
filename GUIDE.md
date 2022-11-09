- make a bot on the discord application portal  
you must enable the message content intent on your bot's application page:  
![](https://i.imgur.com/qL2etXv.png)
- when giving permission to the bot you must give at least the following permissions:  
`Read Messages/View Channels, Send Messages, Manage Messages, Embed Links, Add Reactions`
- clone the repository  
`git clone https://github.com/Roguezilla/starboard-go.git`  
also run  
`git config core.fileMode false`  
- install reflex  
`go install github.com/cespare/reflex@latest`  
**if you are on ubuntu you can use**  
`sudo apt get install reflex`  
- after installing the requirements run  
`go run .`  
![](https://i.imgur.com/hvOfUzT.png)
- now `CTRL+C` and run:  
`chmod 777 start.sh`  
`./start.sh` (use this command from now one to run the bot)  
this will automatically restart the bot when it gets updated:  
![](https://i.imgur.com/FkGXoSU.png)
- all that's left is to setup the bot for your server:  
`sb!setup <archive_channel> <archive_emote> <archive_emote_amount>`
![](https://i.imgur.com/ex6q23f.png)  