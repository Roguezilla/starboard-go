- setup the bot on discord dev portal  
- you must enable the message content intent on your bot's application page:  
![](https://i.imgur.com/qL2etXv.png)
- when giving permission to the bot you must give at least the following permissions: **Read Messages/View Channels, Send Messages, Manage Messages, Embed Links, Add Reactions**  
- clone the repository
```bash
git clone https://github.com/Roguezilla/starboard-go.git
```
- **if you are on linux** run this command
```bash
git config core.fileMode false
```
- **if you are on linux** make sure everything is owned by 1 user if you like running things on different users
- **if you are on linux** i would recommend chmod 777ing everything
- after installing the requirements run:  
```bash
go run .
```
![](https://i.imgur.com/PXNRQog.png)
- all that's left is to setup the bot for your server:
```bash
sb!setup <archive_channel> <archive_emote> <archive_emote_amount>
```
![](https://i.imgur.com/ex6q23f.png)  